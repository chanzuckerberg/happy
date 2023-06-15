package config_manager

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

func normalizeKey(key string) string {
	reg, err := regexp.Compile(`[^a-zA-Z0-9_]+`)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(reg.ReplaceAllString(key, "-")))
}

func findAllDockerfiles(path string) ([]string, error) {
	logrus.Debugf("Searching for Dockerfiles in %s", path)
	paths := []string{}

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".terraform" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".dockerignore") {
			return nil
		}
		if !strings.HasSuffix(path, "Dockerfile") && !strings.Contains(path, "Dockerfile.") {
			return nil
		}
		paths = append(paths, path)
		return nil
	})

	return paths, err
}
