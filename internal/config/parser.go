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
	var pendingTags []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// собираем pendingTags до Host
		if strings.HasPrefix(line, "# tags:") {
			tags := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "# tags:")), ",")
			pendingTags = pendingTags[:0] // сбрасываем
			for i := range tags {
				pendingTags = append(pendingTags, strings.TrimSpace(tags[i]))
			}
			continue
		}

		if strings.HasPrefix(line, "Host ") {
			// сохраняем предыдущий хост
			if current.Host != "" {
				configs = append(configs, current)
			}

			// создаём новый и вшиваем pendingTags
			current = SSHConfig{
				Host: strings.TrimSpace(strings.TrimPrefix(line, "Host")),
				Tags: append([]string{}, pendingTags...), // копируем
			}
			pendingTags = nil
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

	// добавим последний хост
	if current.Host != "" {
		configs = append(configs, current)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}
