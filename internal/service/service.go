package service

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/HT4w5/forklift/internal/config"
	"github.com/HT4w5/forklift/internal/run"
	"github.com/go-co-op/gocron/v2"
)

const (
	startProfileRetries = 5
)

type Service struct {
	cfg     *config.Config
	patches map[string]any
	cron    gocron.Scheduler
	inst    *run.Inst
	logger  *slog.Logger
}

func MakeService(cfg *config.Config) (*Service, error) {
	// Create logger
	var level slog.Level
	switch cfg.Log.Level {
	case "":
		level = slog.LevelError // Default
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	case "info":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "none":
	default:
		return nil, fmt.Errorf("invalid log level \"%s\"", cfg.Log.Level)
	}

	// Create cron scheduler
	cron, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	var logger *slog.Logger
	if cfg.Log.Level == "none" {
		logger = slog.New(slog.DiscardHandler)
	} else {
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: level,
			},
		))
	}

	svc := Service{
		cfg:     cfg,
		patches: make(map[string]any),
		cron:    cron,
		logger:  logger,
	}

	// Create patch map
	for _, v := range cfg.Patches {
		svc.patches[v.Tag] = v.Content
	}

	// Verify specified patches exist
	for _, v := range cfg.Profile.Patches {
		_, ok := svc.patches[v]
		if !ok {
			return nil, fmt.Errorf("patch \"%s\" not defined", v)
		}
	}

	return &svc, nil
}

func (svc *Service) Start() error {
	svc.logger.Info("starting service")
	// Get profile
	success := false
	var profile any
	var err error
	for i := 0; i < startProfileRetries; i++ {
		profile, err = svc.compileProfile()
		if err == nil {
			success = true
			break
		}
		svc.logger.Warn("failed to compile profile", "error", err)
	}

	if !success {
		return fmt.Errorf("failed to compile profile: %w", err)
	}

	// Create instance
	svc.inst, err = run.Create(svc.cfg.Exec, profile)
	if err != nil {
		svc.logger.Error("failed to create instance", "error", err)
		return err
	}

	// Start cron
	_, err = svc.cron.NewJob(
		gocron.CronJob(
			svc.cfg.Profile.Update,
			false,
		),
		gocron.NewTask(
			svc.reload,
		),
	)
	if err != nil {
		svc.logger.Error("failed to setup cron job", "error", err)
		return err
	}
	svc.cron.Start()
	return nil
}

func (svc *Service) Stop() error {
	svc.logger.Info("stopping service")
	// Shutdown cron
	err := svc.cron.Shutdown()
	if err != nil {
		svc.logger.Warn("failed to shutdown cron", "error", err)
	}

	// Destroy instance
	err = svc.inst.Destroy()
	if err != nil {
		svc.logger.Warn("failed to destroy instance", "error", err)
		return err
	}

	return nil
}

// Recompile profile and restart sing-box
func (svc *Service) reload() {
	svc.logger.Info("reloading instance")
	// Reload profile
	profile, err := svc.compileProfile()
	if err != nil {
		svc.logger.Error("failed to compile profile", "error", err)
	}

	// Destroy previous instance
	err = svc.inst.Destroy()
	if err != nil {
		svc.logger.Error("failed to destroy previous instance", "error", err)
	}

	// Start new instance
	svc.inst, err = run.Create(svc.cfg.Exec, profile)
	if err != nil {
		svc.logger.Error("failed to create new instance", "error", err)
	}
}
