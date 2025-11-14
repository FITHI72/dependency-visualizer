package fetch

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadTestRepo читает зависимости из локального файла тестового репозитория.
func LoadTestRepo(path string) (map[string][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть тестовый файл: %v", err)
	}
	defer file.Close()

	repo := make(map[string][]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("неверный формат строки: %s", line)
		}

		pkg := strings.TrimSpace(parts[0])
		deps := strings.Fields(strings.TrimSpace(parts[1]))

		repo[pkg] = deps
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return repo, nil
}
