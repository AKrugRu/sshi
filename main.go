package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/akrugru/sshi/internal/config"
	"github.com/akrugru/sshi/internal/ssh"
	"github.com/akrugru/sshi/internal/tui"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nInterrupted")
		os.Exit(0)
	}()

	var configFile string

	if len(os.Args) >= 2 {
		configFile = filepath.Clean(os.Args[1])
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Unable to determine home directory: %v", err)
		}
		defaultPath := filepath.Join(homeDir, ".ssh", "config")
		if _, err := os.Stat(defaultPath); err == nil {
			configFile = defaultPath
		} else {
			fmt.Println("Usage: sshi <path_to_ssh_config_file>")
			fmt.Println("No argument provided and no ~/.ssh/config found.")
			os.Exit(1)
		}
	}

	configs, err := config.ParseSSHConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to parse SSH config: %v", err)
	}
	if len(configs) == 0 {
		log.Fatalf("No valid SSH hosts found in %s", configFile)
	}

	selected, err := tui.SelectHost(configs)
	if err != nil {
		log.Fatalf("TUI failed: %v", err)
	}
	if selected != nil {
		if err := ssh.ConnectToSSH(*selected); err != nil {
			log.Fatalf("SSH connection failed: %v", err)
		}
	}
}
