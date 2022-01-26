package hostname_manager

import (
	"bufio"
	"os"
	"strings"
)

type HostNameManager struct {
	fileName   string
	containers []string
}

func NewHostNameManager(fileName string, containers []string) *HostNameManager {
	return &HostNameManager{
		fileName:   fileName,
		containers: containers,
	}
}

func (h *HostNameManager) getHostsFileConfig() ([]string, error) {
	file, err := os.Open(h.fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	return txtlines, nil
}

func (h *HostNameManager) getFileBorders() ([]string, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dirPathSlice := strings.Split(path, "/")
	currDir := dirPathSlice[len(dirPathSlice)-1]

	border := "# ==== Happy docker-compose DNS for " + currDir + " ==="

	return []string{border, border}, nil
}

func (h *HostNameManager) cleanUpConfig(borders []string, config []string) []string {
	var newConfig []string
	writeLines := true

	for _, line := range config {
		if writeLines && line == borders[0] {
			writeLines = false
			continue
		}

		if !writeLines && line == borders[1] {
			writeLines = true
			continue
		}

		if writeLines {
			newConfig = append(newConfig, line)
		}
	}

	return newConfig
}

func (h *HostNameManager) generateConfig(borders []string, containers []string) []string {
	var lines []string
	lines = append(lines, borders[0])

	for _, container := range containers {
		lines = append(lines, "127.0.0.1\t"+container)
	}

	lines = append(lines, borders[1])
	return lines
}

func (h *HostNameManager) updateConfig(config []string, newConfig []string) error {
	file, err := os.OpenFile(h.fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	file.Truncate(0)
	writer := bufio.NewWriter(file)

	for _, line := range config {
		_, _ = writer.WriteString(line + "\n")
	}

	for _, line := range newConfig {
		_, _ = writer.WriteString(line + "\n")
	}

	writer.Flush()

	return nil
}

func (h *HostNameManager) Install() error {
	config, err := h.getHostsFileConfig()
	if err != nil {
		return err
	}

	containers := h.containers
	borders, err := h.getFileBorders()
	if err != nil {
		return err
	}

	config = h.cleanUpConfig(borders, config)
	err = h.updateConfig(config, h.generateConfig(borders, containers))
	return err
}

func (h *HostNameManager) UnInstall() error {
	config, err := h.getHostsFileConfig()
	if err != nil {
		return err
	}

	borders, err := h.getFileBorders()
	if err != nil {
		return err
	}

	config = h.cleanUpConfig(borders, config)
	err = h.updateConfig(config, nil)
	return err
}
