package cmd

import (
	"fmt"
	"log"
	"os"

	"dependency-visualizer/internal/config"
	"dependency-visualizer/internal/validate"
)

func Execute() {
	cfg, err := config.ParseFlags()
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(2)
	}

	// Валидация параметров
	if err := validate.All(cfg); err != nil {
		log.Printf("error: %v", err)
		os.Exit(2)
	}

	// Вывод конфигурации
	fmt.Printf("name=%s\n", cfg.Name)
	fmt.Printf("repo=%s\n", cfg.Repo)
	fmt.Printf("test-mode=%s\n", cfg.TestMode)
	fmt.Printf("version=%s\n", cfg.Version)
}
