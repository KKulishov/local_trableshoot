package net

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type NetStat struct {
	Interface string
	RxBytes   int64
	TxBytes   int64
	Command   string
}

func getCommandByPID(pid int) (string, error) {
	path := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(path) // Используем os.ReadFile вместо ioutil.ReadFile
	if err != nil {
		return "", fmt.Errorf("could not read cmdline: %v", err)
	}
	cmdline := strings.ReplaceAll(string(data), "\x00", " ")
	return strings.TrimSpace(cmdline), nil
}

func getNetStatsByPID(pid int) ([]NetStat, error) {
	path := fmt.Sprintf("/proc/%d/net/dev", pid)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	var netStats []NetStat
	scanner := bufio.NewScanner(file)

	// Пропускаем заголовки
	for i := 0; i < 2; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected format in %s", path)
		}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)

		if len(fields) < 17 {
			return nil, fmt.Errorf("unexpected format in %s", path)
		}

		interfaceName := strings.TrimSuffix(fields[0], ":")
		rxBytes, err1 := strconv.ParseInt(fields[1], 10, 64)
		txBytes, err2 := strconv.ParseInt(fields[9], 10, 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("error parsing bytes for interface %s", interfaceName)
		}

		netStats = append(netStats, NetStat{
			Interface: interfaceName,
			RxBytes:   rxBytes,
			TxBytes:   txBytes,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return netStats, nil
}

func RunAnalNet() {
	processes, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Printf("Failed to read /proc: %v\n", err)
		return
	}
	for _, proc := range processes {
		if !proc.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(proc.Name()) // Проверяем, является ли имя "PID"
		if err != nil || pid <= 0 {
			continue
		}
		command, err := getCommandByPID(pid)
		if err != nil {
			// Пропускаем процессы, для которых не можем получить команду
			continue
		}
		netStats, err := getNetStatsByPID(pid)
		if err != nil {
			// Пропускаем, если не можем прочитать сетевую статистику
			continue
		}

		for _, stat := range netStats {
			stat.Command = command
			fmt.Printf("PID: %d, Command: %s, Interface: %s, RxBytes: %d, TxBytes: %d\n",
				pid, stat.Command, stat.Interface, stat.RxBytes, stat.TxBytes)
		}
	}
}
