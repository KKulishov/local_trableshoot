package cgroups

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// сопостовление pid процесса c pods/namespace/container при условии что /var/log/containers существует
type ContainerInfo struct {
	PID           int    // PID процесса
	PodName       string // Имя пода
	Namespace     string // Namespace
	ContainerName string // Имя контейнера
}

// getContainerInfo извлекает информацию о контейнере по PID
func getContainerInfo(pid int) (*ContainerInfo, error) {
	cgroupPath := fmt.Sprintf("/proc/%d/cgroup", pid)
	file, err := os.Open(cgroupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cgroup file for PID %d: %w", pid, err)
	}
	defer file.Close()

	var containerID string
	scanner := bufio.NewScanner(file)

	// Читаем только первую строку
	if scanner.Scan() {
		line := scanner.Text()
		//fmt.Printf("PID %d first cgroup content: %s\n", pid, line) // Отладочный вывод
		parts := strings.Split(line, "/")
		if len(parts) > 0 {
			containerID = parts[len(parts)-1] // Берём последний элемент пути
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read cgroup file for PID %d: %w", pid, err)
	}

	if containerID == "" {
		//fmt.Printf("No container ID found for PID %d\n", pid)
		return nil, fmt.Errorf("container ID not found for PID %d", pid)
	}

	//fmt.Printf("Looking for container ID: %s\n", containerID) // Отладочный вывод

	logDir := "/var/log/containers"
	var logFile string
	err = filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, containerID) {
			logFile = path
			//fmt.Printf("Found log file for PID %d: %s\n", pid, logFile) // Отладочный вывод
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find log file for container ID %s: %w", containerID, err)
	}

	if logFile == "" {
		//fmt.Printf("No log file found for container ID: %s\n", containerID)
		return nil, fmt.Errorf("log file not found for container ID: %s", containerID)
	}

	filename := filepath.Base(logFile)
	//fmt.Printf("Parsing log file name: %s\n", filename) // Отладочный вывод

	parts := strings.Split(filename, "_")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid log file format: %s", filename)
	}

	info := &ContainerInfo{
		PID:           pid,
		PodName:       parts[0],
		Namespace:     parts[1],
		ContainerName: strings.TrimSuffix(parts[2], fmt.Sprintf("-%s.log", containerID)),
	}

	return info, nil
}

// processPIDs обрабатывает список PID и собирает информацию о контейнерах
func ProcessPIDs(pids []int) ([]ContainerInfo, error) {
	var containerInfos []ContainerInfo

	for _, pid := range pids {
		info, err := getContainerInfo(pid)
		if err != nil {
			//fmt.Printf("Error processing PID %d: %v\n", pid, err)
			continue // Пропускаем ошибки для отдельных PID
		}
		containerInfos = append(containerInfos, *info)
	}

	return containerInfos, nil
}
