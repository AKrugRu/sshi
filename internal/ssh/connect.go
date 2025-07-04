package ssh

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/akrugru/sshi/internal/config"
)

func sanitize(input string) string {
	if strings.ContainsAny(input, ";|&$`()") {
		log.Fatalf("Invalid input in SSH config: %q", input)
	}
	return input
}

func ConnectToSSH(cfg config.SSHConfig) error {
	args := []string{"-t"}

	if user := sanitize(cfg.User); user != "" {
		args = append(args, "-l", user)
	}
	if port := sanitize(cfg.Port); port != "" {
		args = append(args, "-p", port)
	}
	if key := sanitize(cfg.IdentityFile); key != "" {
		args = append(args, "-i", key)
		if _, err := os.Stat(key); os.IsNotExist(err) {
			log.Fatalf("IdentityFile does not exist: %s", key)
		}
	}
	host := cfg.HostName
	if host == "" {
		host = cfg.Host
	}
	if host == "" {
		log.Fatalf("Missing HostName and Host in SSH config")
	}
	args = append(args, sanitize(host))

	log.Printf("Executing: ssh %s", strings.Join(args, " "))
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
