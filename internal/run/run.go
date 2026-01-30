package run

import (
	"encoding/json"
	"fmt"
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

	// Check if process exited immediately
	exited := make(chan error, 1)
	go func() {
		exited <- cmd.Wait()
	}()

	select {
	case err := <-exited:
		os.RemoveAll(tempDir) // Clean up temp directory
		return nil, fmt.Errorf("process exited immediately: %w", err)
	case <-time.After(10 * time.Second):
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

	select {
	case <-done:
		return nil
	case <-time.After(shutdownTimeout):
		if killErr := in.cmd.Process.Kill(); killErr != nil && killErr != os.ErrProcessDone {
			return killErr
		}
		return nil
	}
}
