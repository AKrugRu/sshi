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
		raw := scanner.Text()
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		// Обрабатываем комментарий с тегами — допускаем ведущие пробелы перед #
		if strings.HasPrefix(line, "#") {
			afterHash := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			// ожидаем "tags:"
			if strings.HasPrefix(afterHash, "tags:") {
				tagsStr := strings.TrimSpace(strings.TrimPrefix(afterHash, "tags:"))
				parts := strings.Split(tagsStr, ",")
				pendingTags = pendingTags[:0] // очистить slice (без аллокации нового)
				for _, t := range parts {
					t = strings.TrimSpace(t)
					if t != "" {
						pendingTags = append(pendingTags, t)
					}
				}
			}
			continue
		}

		// Обрабатываем Host — берём всё, что идёт после слова "Host"
		if strings.HasPrefix(line, "Host ") {
			// если был предыдущий текущий хост — сохраним его
			if current.Host != "" {
				configs = append(configs, current)
			}

			// создаём новый current и вшиваем pendingTags
			// безопасно извлекаем часть после 'Host '
			hostPart := strings.TrimSpace(strings.TrimPrefix(line, "Host"))
			// hostPart может начинаться с пробела — уберём его
			hostPart = strings.TrimSpace(hostPart)

			current = SSHConfig{
				Host: strings.TrimSpace(hostPart),
				Tags: append([]string{}, pendingTags...), // копируем slice
			}
			// сбросим pendingTags: если для следующего блока не будет тега, останется nil/пустой
			pendingTags = nil
			continue
		}

		// остальные директивы — только если у нас есть текущий хост
		if current.Host == "" {
			// директивы до первого Host игнорируем
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key, val := fields[0], strings.Join(fields[1:], " ")

		// убрать возможные кавычки вокруг значения
		val = strings.Trim(val, `"'`)

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

	// добавим последний хост, если есть
	if current.Host != "" {
		configs = append(configs, current)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}
