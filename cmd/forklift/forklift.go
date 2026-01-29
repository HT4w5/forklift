package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/HT4w5/forklift/internal/config"
	"github.com/HT4w5/forklift/internal/meta"
	"github.com/HT4w5/forklift/internal/service"
)

const (
	exitConfigErr = iota
	exitServiceErr
	exitShutdownErr
)

func main() {
	var configPath string
	var showVersion bool
	var showHelp bool

	flag.StringVar(&configPath, "c", "", "path to config file")
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.BoolVar(&showVersion, "v", false, "show version info")
	flag.BoolVar(&showVersion, "version", false, "show version info")
	flag.BoolVar(&showHelp, "h", false, "show help")
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Print(meta.VersionMultiline())
		os.Exit(0)
	}

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "Error: missing config file path\n")
		flag.Usage()
		os.Exit(exitConfigErr)
	}

	// Load configuration
	cfg := config.Default()
	if err := cfg.Load(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(exitConfigErr)
	}

	// Create service
	svc, err := service.MakeService(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create service: %v\n", err)
		os.Exit(exitServiceErr)
	}

	// Start service
	if err := svc.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start service: %v\n", err)
		os.Exit(exitServiceErr)
	}

	fmt.Println("Service started successfully")

	// Wait for shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("Received signal: %v, shutting down...\n", sig)

	if err := svc.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during shutdown: %v", err)
		os.Exit(exitShutdownErr)
	}
}
