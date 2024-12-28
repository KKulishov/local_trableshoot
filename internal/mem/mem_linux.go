package mem

import (
	"bufio"
	"fmt"
	"local_trableshoot/internal/cgroups"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type ProcessSwap struct {
	PID     int
	Swap    float64
	Name    string
	Command string
}

func AddProcessesByMem(file *os.File) ([]int, error) {
	// Заголовок секции
	fmt.Fprintln(file, "<h3>Processes by MEM</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Выполнение команды для получения списка процессов по памяти
	cmd := exec.Command("sh", "-c", `ps -ewwwo pid,size,command --sort -size | head -n 20 | awk '{ pid=$1 ; printf("%7s ", pid) }{ hr=$2/1024 ; printf("%8.2f Mb ", hr) } { for ( x=3 ; x<=NF ; x++ ) { printf("%s ",$x) } print "" }'`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при выполнении команды: %s\n", err)
	} else {
		// Запись результата в файл
		fmt.Fprintln(file, string(output))
	}

	fmt.Fprintln(file, "</pre></div>")
	var pids []int
	// Парсинг вывода
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		// Пропускаем заголовок и пустые строки
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		// Разделяем строку на поля
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue // Пропускаем строки с недостаточным количеством полей
		}

		// Парсим PID (первое поле)
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue // Пропускаем строки с некорректным PID
		}

		// Добавляем PID в массив
		pids = append(pids, pid)
	}

	return pids, nil
}

func SaveContainersToHTML(file *os.File, containerInfos []cgroups.ContainerInfo) error {
	// Начало HTML-страницы и таблицы
	_, err := file.WriteString("<html><body><h3 id=\"Containers\">Container Info by top Pid, mem used for Pods </h3><table border='1'>")
	if err != nil {
		return err
	}

	// Заголовки таблицы
	_, err = file.WriteString("<tr><th>PID</th><th>Pod</th><th>Namespace</th><th>Container Name</th></tr>")
	if err != nil {
		return err
	}

	// Строки таблицы
	for _, info := range containerInfos {
		row := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			info.PID, info.PodName, info.Namespace, info.ContainerName)
		_, err := file.WriteString(row)
		if err != nil {
			return err
		}
	}

	// Закрытие таблицы и HTML-страницы
	_, err = file.WriteString("</table></body></html>")
	if err != nil {
		return err
	}

	return nil
}

func ShowMem(file *os.File) {
	// Show mem linux
	format.WriteHeader(file, "Show all mem")
	currentOutput := format.ExecuteCommand("free", "-m")
	format.WritePreformatted(file, currentOutput)
}

func getProcessInfo(pid int) (ProcessSwap, error) {
	var proc ProcessSwap
	proc.PID = pid

	// Получаем имя процесса из stat
	statFile := fmt.Sprintf("/proc/%d/stat", pid)
	data, err := os.ReadFile(statFile)
	if err != nil {
		return proc, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 2 {
		return proc, fmt.Errorf("invalid stat file format")
	}
	proc.Name = strings.Trim(fields[1], "()")

	// Считываем swap usage из smaps
	smapsFile := fmt.Sprintf("/proc/%d/smaps", pid)
	file, err := os.Open(smapsFile)
	if err != nil {
		return proc, err
	}
	defer file.Close()

	var swapSum float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "Swap:") {
			parts := strings.Fields(scanner.Text())
			if len(parts) > 1 {
				swap, _ := strconv.ParseFloat(parts[1], 64)
				swapSum += swap
			}
		}
	}
	proc.Swap = swapSum / 1024 // Конвертируем в MB

	// Получаем команду запуска
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdlineData, err := os.ReadFile(cmdlineFile)
	if err != nil {
		return proc, err
	}
	proc.Command = strings.ReplaceAll(string(cmdlineData), "\x00", " ")

	return proc, nil
}

func ShowProcessesSwapUsage(file *os.File) error {
	var processes []ProcessSwap

	// Перебираем все процессы в /proc
	dir, err := os.Open("/proc")
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		pid, err := strconv.Atoi(file)
		if err == nil && pid > 0 { // Проверяем, что это числовой PID
			if proc, err := getProcessInfo(pid); err == nil && proc.Swap > 0 {
				processes = append(processes, proc)
			}
		}
	}

	// Сортируем по использованию swap
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Swap > processes[j].Swap
	})

	// Заголовок HTML-таблицы
	fmt.Fprintln(file, "<h3>Processes by SWAP top 20 utilization</h3>")
	fmt.Fprintln(file, "<table border='1'>")
	fmt.Fprintln(file, "<tr><th>PID</th><th>SWAP (MB)</th><th>Program</th><th>COMMAND</th></tr>")

	// Выводим топ 20 процессов по использованию swap
	for i, proc := range processes {
		if i >= 20 {
			break
		}
		row := fmt.Sprintf(
			"<tr><td>%d</td><td>%.2f</td><td>%s</td><td>%s</td></tr>",
			proc.PID, proc.Swap, proc.Name, proc.Command,
		)
		fmt.Fprintln(file, row)
	}
	// Закрываем таблицу
	fmt.Fprintln(file, "</table>")

	return nil
}
