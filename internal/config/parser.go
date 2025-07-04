package config

import (
	"bufio"
	"os"
	"strings"
)

type SSHConfig struct {
	Host         string
	HostName     string
	User         string
	Port         string
	IdentityFile string
	Tags         []string
}

func ParseSSHConfig(filePath string) ([]SSHConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var configs []SSHConfig
	var current SSHConfig
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# tags:") {
			if current.Host == "" {
				continue
			}
			tags := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "# tags:")), ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			current.Tags = tags
			continue
		}

		if strings.HasPrefix(line, "Host ") {
			if current.Host != "" {
				configs = append(configs, current)
			}
			current = SSHConfig{Host: strings.TrimSpace(strings.TrimPrefix(line, "Host"))}
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key, val := fields[0], strings.Join(fields[1:], " ")

		switch key {
		case "HostName":
			current.HostName = val
		case "User":
			current.User = val
		case "Port":
			current.Port = val
		case "IdentityFile":
			current.IdentityFile = val
		}
	}

	if current.Host != "" {
		configs = append(configs, current)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}
