package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Dependency описывает одну зависимость пакета.
type Dependency struct {
	ID    string
	Range string
}

// GetNuGetDependencies получает прямые зависимости из NuGet API.
func GetNuGetDependencies(pkgName, version string) ([]Dependency, error) {
	// Формируем URL запроса.
	url := fmt.Sprintf("https://api.nuget.org/v3/registration5-gz-semver2/%s/index.json", strings.ToLower(pkgName))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к NuGet API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("NuGet API вернул статус %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ: %w", err)
	}

	// Структуры для парсинга JSON.
	var data struct {
		Items []struct {
			Items []struct {
				CatalogEntry struct {
					Version          string `json:"version"`
					DependencyGroups []struct {
						Dependencies []struct {
							ID    string `json:"id"`
							Range string `json:"range"`
						} `json:"dependencies"`
					} `json:"dependencyGroups"`
				} `json:"catalogEntry"`
			} `json:"items"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %w", err)
	}

	// Ищем нужную версию
	for _, outer := range data.Items {
		for _, inner := range outer.Items {
			entry := inner.CatalogEntry
			if strings.EqualFold(entry.Version, version) {
				// Нашли нужную версию — собираем зависимости.
				var deps []Dependency
				for _, group := range entry.DependencyGroups {
					for _, dep := range group.Dependencies {
						deps = append(deps, Dependency{ID: dep.ID, Range: dep.Range})
					}
				}
				return deps, nil
			}
		}
	}

	return nil, errors.New("указанная версия не найдена в NuGet")
}
