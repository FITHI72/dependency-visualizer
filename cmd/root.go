package cmd

import (
	"dependency-visualizer/internal/config"
	"dependency-visualizer/internal/fetch"
	"dependency-visualizer/internal/graph"
	"dependency-visualizer/internal/validate"
	"flag"
	"fmt"
	"os"
	"strings"
)

func Execute() {
	cfg := &config.Config{}
	var operation = flag.String("op", "graph", "operation mode: graph or order")

	flag.StringVar(&cfg.Name, "name", "", "Имя анализируемого пакета")
	flag.StringVar(&cfg.Repo, "url", "", "URL-адрес репозитория")
	flag.StringVar(&cfg.TestMode, "mode", "online", "Режим работы (online/offline/test)")
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
	fmt.Printf("Операция: %s\n", *operation)
	fmt.Println("-------------------------")

	// Работа в тестовом режиме
	if cfg.TestMode == "test" {
		repo, err := fetch.LoadTestRepo(cfg.Repo)
		if err != nil {
			fmt.Println("Ошибка загрузки тестового репозитория:", err)
			os.Exit(1)
		}

		g := graph.NewGraph()
		g.BuildDFS(cfg.Name, func(pkg string) []string {
			return repo[pkg]
		})

		switch *operation {
		case "graph":
			fmt.Println("\nГраф зависимостей:")
			g.PrintGraph()
		case "order":
			order, cycle := g.LoadOrder(cfg.Name)

			fmt.Println("\nПорядок загрузки зависимостей:")

			if cycle != nil {
				fmt.Println("⚠ Обнаружен цикл в зависимостях:", cycle)
				return
			}
			fmt.Println(strings.Join(order, " → "))
		default:
			fmt.Printf("Неизвестная операция: %s (используй 'graph' или 'order')\n", *operation)
		}
	} else {
		fmt.Println("Режим remote: пока используется только для прямых зависимостей (этап 2)")
	}
}
