package top

import (
	"fmt"
	"local_trableshoot/configs"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
	"strings"
	"time"
	// Импортируем пакет configs
)

func Get_atop_processes_lists(file *os.File) {
	// Проверка на наличие утилиты atop
	if !isCommandAvailable("atop") {
		fmt.Fprintln(file, "<h3>Утилита atop не найдена</h3>")
		fmt.Fprintln(file, "<div><pre>Убедитесь, что atop установлена в системе.</pre></div>")
	} else {
		// Определяем временные параметры
		start := time.Now().Add(-15 * time.Minute).Format("15:04")
		end := time.Now().Format("15:04")
		today := time.Now().Format("20060102")
		logFile := fmt.Sprintf("%s/atop_%s", configs.AtopLogDir, today)
		// Функция для записи секции
		writeSection := func(header, command string) {
			fmt.Fprintln(file, fmt.Sprintf("<h3>%s</h3>", header))
			fmt.Fprintln(file, "<div><pre>")
			// Выполнение команды
			cmd := exec.Command("sh", "-c", command)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Fprintf(file, "Ошибка при выполнении команды %s: %s\n", header, err)
			} else {
				fmt.Fprintln(file, string(output))
			}
			fmt.Fprintln(file, "</pre></div>")
		}

		// Запись информации о процессах по CPU
		writeSection("Atop processes by CPU, in the last 15 minutes", fmt.Sprintf("atopsar -r \"%s\" -b %s -e %s -O", logFile, start, end))
		// Запись информации о процессах по MEM
		writeSection("Atop processes by MEM, in the last 15 minutes", fmt.Sprintf("atopsar -r \"%s\" -b %s -e %s -G", logFile, start, end))
		// Запись информации о процессах по IOPS
		writeSection("Atop processes by IOPS, in the last 15 minutes", fmt.Sprintf("atopsar -r \"%s\" -b %s -e %s -D", logFile, start, end))
		// Проверка на наличие утилиты netatop
		if !isCommandAvailable("netatop") {
			fmt.Fprintln(file, "<h3>Утилита netatop не обнаружена</h3>")
			fmt.Fprintln(file, "<div><pre>Убедитесь, что netatop установлена в системе.</pre></div>")
		} else {
			// Запись информации о процессах по NET
			writeSection("Atop processes by NET, in the last 15 minutes", fmt.Sprintf("atopsar -r \"%s\" -b %s -e %s -N", logFile, start, end))
		}
	}
}

func GetSummary(file *os.File) {
	// Проверка на наличие утилиты atop
	if !isCommandAvailable("atop") {
		fmt.Fprintln(file, "<h3>Утилита atop не найдена</h3>")
		fmt.Fprintln(file, "<div><pre>Убедитесь, что atop установлена в системе.</pre></div>")
	} else {
		// Выполнение команды для вывода информации о CPU с помощью atop
		format.WriteHeaderWithID(file, "Show atop process", "Atop")
		cpuCmd := exec.Command("sh", "-c", "atop -L 180 -a 1 1 | sed -rn '1,/^\\s+/ p' | tail -n +3 | head -n -1")
		cpuOutput, err := cpuCmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(file, "Ошибка при выполнении команды для CPU:", err)
		} else {
			// Запись результатов в файл
			fmt.Fprintln(file, "<div><pre>")
			fmt.Fprintln(file, string(cpuOutput))
			fmt.Fprintln(file, "</pre></div>")
		}
	}

	// Запись заголовка для раздела "Sessions"
	fmt.Fprintln(file, "<h3>Sessions</h3>")

	// Выполнение команды для получения информации о сессиях
	sessionsCmd := exec.Command("who", "-a", "-H")
	sessionsOutput, err := sessionsCmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(file, "Ошибка при выполнении команды для Sessions:", err)
	} else {
		// Запись результатов в файл
		fmt.Fprintln(file, "<div><pre>")
		fmt.Fprintln(file, string(sessionsOutput))
		fmt.Fprintln(file, "</pre></div>")
	}
}

func isCommandAvailable(command string) bool {
	cmd := exec.Command("which", command)
	output, err := cmd.CombinedOutput()
	return err == nil && len(strings.TrimSpace(string(output))) > 0
}
