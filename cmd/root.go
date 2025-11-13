package cmd

import (
	"dependency-visualizer/internal/config"
	"dependency-visualizer/internal/fetch"
	"dependency-visualizer/internal/validate"
	"flag"
	"fmt"
	"os"
)

func Execute() {
	cfg := &config.Config{}

	flag.StringVar(&cfg.Name, "name", "", "Имя анализируемого пакета")
	flag.StringVar(&cfg.Repo, "url", "", "URL-адрес репозитория")
	flag.StringVar(&cfg.TestMode, "mode", "online", "Режим работы (online/offline)")
	flag.StringVar(&cfg.Version, "version", "", "Версия пакета")
	flag.Parse()

	if err := validate.All(cfg); err != nil {
		fmt.Println("Ошибка в параметрах:", err)
		os.Exit(1)
	}

	fmt.Println("\n--- Параметры запуска ---")
	fmt.Printf("Пакет: %s\n", cfg.Name)
	fmt.Printf("URL: %s\n", cfg.Repo)
	fmt.Printf("Режим: %s\n", cfg.TestMode)
	fmt.Printf("Версия: %s\n", cfg.Version)
	fmt.Println("-------------------------\n")

	// получаем зависимости из NuGet
	deps, err := fetch.GetNuGetDependencies(cfg.Name, cfg.Version)
	if err != nil {
		fmt.Println("Ошибка при получении зависимостей:", err)
		os.Exit(1)
	}

	if len(deps) == 0 {
		fmt.Println("У пакета нет прямых зависимостей.")
		return
	}

	fmt.Printf("Прямые зависимости пакета %s (%s):\n", cfg.Name, cfg.Version)
	for _, d := range deps {
		fmt.Printf("  - %s %s\n", d.ID, d.Range)
	}
}
