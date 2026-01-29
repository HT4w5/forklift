package run

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/HT4w5/forklift/internal/config"
)

const (
	tempDirPrefix  = "forklift-"
	configFileName = "/config.json"
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

	err := in.cmd.Process.Kill()
	if err != nil {
		return err
	}

	return in.cmd.Wait()
}
