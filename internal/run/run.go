package run

import (
	"encoding/json"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/HT4w5/forklift/internal/config"
)

const (
	tempDirPrefix   = "forklift-"
	configFileName  = "/config.json"
	shutdownTimeout = time.Second * 10
)

// Represents a running sing-box instance
type Inst struct {
	cfg     config.ExecConfig
	tempDir string
	cmd     *exec.Cmd
}

func Create(cfg config.ExecConfig, profile any) (*Inst, error) {
	// Look for sing-box binary
	path, err := exec.LookPath(cfg.Path)
	if err != nil {
		return nil, err
	}

	// Create temporary working directory
	tempDir, err := os.MkdirTemp("", tempDirPrefix)
	if err != nil {
		return nil, err
	}

	// Write config
	configBytes, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}

	cfgPath := tempDir + configFileName

	err = os.WriteFile(
		cfgPath,
		configBytes,
		0600,
	)
	if err != nil {
		return nil, err
	}

	// Create command
	cmd := exec.Command(
		path,
		"run",
		"-c",
		cfgPath,
		"-D",
		tempDir,
	)

	// Forward log
	if cfg.LogFwd {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Start process
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return &Inst{
		cfg:     cfg,
		tempDir: tempDir,
		cmd:     cmd,
	}, nil
}

func (in *Inst) Destroy() error {
	// Clean up tempDir in the end
	defer os.RemoveAll(in.tempDir)

	if in.cmd.Process == nil {
		return nil
	}

	err := in.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil && err != os.ErrProcessDone {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- in.cmd.Wait()
	}()

	timer := time.NewTimer(shutdownTimeout)
	defer timer.Stop()

	select {
	case <-done:
		return nil
	case <-timer.C:
		if killErr := in.cmd.Process.Kill(); killErr != nil && killErr != os.ErrProcessDone {
			return killErr
		}
		return nil
	}
}
