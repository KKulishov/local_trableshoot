package proc

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ProcessThread содержит информацию о процессе и его потоках
type ProcessThread struct {
	PID        int
	Name       string
	Command    string
	Threads    int
	CPUUsage   float64
	MemUsageMB float64
}

// GetProcessThreadsInfo собирает информацию о процессах и их потоках
func GetProcessThreadsInfo() ([]ProcessThread, error) {
	var processes []ProcessThread
	var wg sync.WaitGroup
	dir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	totalMemMB := getTotalMemoryMB()

	for _, file := range files {
		pid, err := strconv.Atoi(file)
		if err != nil || pid <= 0 {
			continue
		}

		wg.Add(1)
		go func(pid int) {
			defer wg.Done()
			if proc, err := getProcessInfo(pid, totalMemMB); err == nil {
				processes = append(processes, proc)
			}
		}(pid)
	}

	wg.Wait()
	return processes, nil
}

// getProcessInfo собирает информацию о конкретном процессе
func getProcessInfo(pid int, totalMemMB float64) (ProcessThread, error) {
	var proc ProcessThread
	proc.PID = pid

	// Читаем имя процесса и CPU использование из stat
	statFile := fmt.Sprintf("/proc/%d/stat", pid)
	file, err := os.Open(statFile)
	if err != nil {
		return proc, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 {
			proc.Name = strings.Trim(fields[1], "()")
		}
		if len(fields) >= 14 { // Поля для CPU
			utime, _ := strconv.ParseFloat(fields[13], 64)
			stime, _ := strconv.ParseFloat(fields[14], 64)
			cutime, _ := strconv.ParseFloat(fields[15], 64)
			cstime, _ := strconv.ParseFloat(fields[16], 64)
			proc.CPUUsage = utime + stime + cutime + cstime
		}
	}

	// Читаем количество потоков и использование памяти из status
	statusFile := fmt.Sprintf("/proc/%d/status", pid)
	file, err = os.Open(statusFile)
	if err != nil {
		return proc, err
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Threads:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				proc.Threads, _ = strconv.Atoi(parts[1])
			}
		} else if strings.HasPrefix(line, "VmRSS:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				v, _ := strconv.ParseFloat(parts[1], 64)
				proc.MemUsageMB = v / 1024 // kB to MB
			}
		}
	}

	// Читаем команду запуска из cmdline
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(cmdlineFile)
	if err != nil {
		return proc, err
	}
	proc.Command = strings.ReplaceAll(string(data), "\x00", " ")

	return proc, nil
}

// getTotalMemoryMB получает общее количество RAM в системе (в MB)
func getTotalMemoryMB() float64 {
	meminfo, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(meminfo), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				total, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				return total / 1024 // KB -> MB
			}
		}
	}
	return 0
}

// ShowTopThreads выводит топ-20 процессов по количеству потоков
func ShowTopThreads(file *os.File) error {
	processes, err := GetProcessThreadsInfo()
	if err != nil {
		return err
	}

	// Сортируем процессы по количеству потоков
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Threads > processes[j].Threads
	})

	// Вывод HTML
	fmt.Fprintln(file, "<h3>Top 20 Processes by Threads</h3>")
	fmt.Fprintln(file, "<table border='1'>")
	fmt.Fprintln(file, "<tr><th>PID</th><th>Name</th><th>Command</th><th>Threads Count</th></tr>")
	for i, proc := range processes {
		if i >= 20 {
			break
		}
		row := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%d</td></tr>", proc.PID, proc.Name, proc.Command, proc.Threads)
		fmt.Fprintln(file, row)
	}
	fmt.Fprintln(file, "</table>")
	fmt.Fprintln(file, "<table border='1'>")
	fmt.Fprintln(file, "<h3> Man threads diagnostics: </h3>")
	fmt.Fprintf(file, "<p>Show stat threads for pid: pidstat -t -p [PID]</p>\n")
	fmt.Fprintf(file, "<p>Show all threads in system: ps -eLf | wc -l</p>\n")
	return nil
}

// ShowTopThreadsByCPU выводит топ-20 процессов по CPU
/* ToDo
func ShowTopThreadsByCPU(file *os.File) error {
	processes, err := GetProcessThreadsInfo()
	if err != nil {
		return err
	}

	// Сортируем по использованию CPU
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPUUsage > processes[j].CPUUsage
	})

	// Вывод HTML
	fmt.Fprintln(file, "<h3>Top 20 Processes by CPU Usage</h3>")
	fmt.Fprintln(file, "<table border='1'>")
	fmt.Fprintln(file, "<tr><th>PID</th><th>Name</th><th>Command</th><th>Threads</th><th>CPU Usage</th><th>Mem Usage (MB)</th></tr>")
	for i, proc := range processes {
		if i >= 20 {
			break
		}
		row := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td>%.2f</td><td>%.2f</td></tr>", proc.PID, proc.Name, proc.Command, proc.Threads, proc.CPUUsage, proc.MemUsageMB)
		fmt.Fprintln(file, row)
	}
	fmt.Fprintln(file, "</table>")
	return nil
}
*/

// ShowTopThreadsByMem выводит топ-20 процессов по памяти
func ShowTopThreadsByMem(file *os.File) error {
	processes, err := GetProcessThreadsInfo()
	if err != nil {
		return err
	}

	// Сортируем по использованию памяти
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].MemUsageMB > processes[j].MemUsageMB
	})

	// Вывод HTML
	fmt.Fprintln(file, "<h3>Top 20 Threads by Memory Usage</h3>")
	fmt.Fprintln(file, "<table border='1'>")
	fmt.Fprintln(file, "<tr><th>PID</th><th>Name</th><th>Command</th><th>Threads count</th><th>Mem Usage (MB)</th></tr>")
	for i, proc := range processes {
		if i >= 20 {
			break
		}
		row := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%d</td><td>%.2f</td></tr>", proc.PID, proc.Name, proc.Command, proc.Threads, proc.MemUsageMB)
		fmt.Fprintln(file, row)
	}
	fmt.Fprintln(file, "</table>")
	return nil
}
