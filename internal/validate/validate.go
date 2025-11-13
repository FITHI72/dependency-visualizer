package validate

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"dependency-visualizer/internal/config"
)

func All(cfg *config.Config) error {
	if err := Name(cfg.Name); err != nil {
		return err
	}
	if err := Repo(cfg.Repo); err != nil {
		return err
	}
	if err := TestMode(cfg.TestMode); err != nil {
		return err
	}
	if err := Version(cfg.Version); err != nil {
		return err
	}
	return nil
}

func Name(n string) error {
	if n == "" {
		return errors.New("missing required parameter --name")
	}
	if filepath.Base(n) != n {
		return errors.New("invalid package name (must not contain path separators)")
	}
	return nil
}

func Repo(r string) error {
	if r == "" {
		return errors.New("missing required parameter --repo")
	}
	u, err := url.Parse(r)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https" || u.Scheme == "git" || u.Scheme == "ssh") {
		return nil
	}
	if _, err := os.Stat(r); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("repo path does not exist: %s", r)
		}
		return fmt.Errorf("cannot access repo path: %v", err)
	}
	return nil
}

func TestMode(m string) error {
	switch m {
	case "off", "local", "remote":
		return nil
	default:
		return fmt.Errorf("invalid --test-mode: %s (allowed: off, local, remote)", m)
	}
}

func Version(v string) error {
	if v == "" {
		return errors.New("missing required parameter --version")
	}
	if v == "latest" {
		return nil
	}

	// Проверяем формат версии по SemVer
	matched, _ := regexp.MatchString(`^(\d+\.\d+\.\d+)([-+][A-Za-z0-9.-]+)?$`, v)
	if !matched {
		return fmt.Errorf("invalid --version: %s (expected semver like 1.2.3 or 'latest')", v)
	}
	return nil
}
