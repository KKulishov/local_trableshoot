package kernel

import (
	"bufio"
	"fmt"
	"local_trableshoot/internal/format"
	"os"
	"regexp"
	"strings"
	"time"
)

func GetKernelAndModules(file *os.File) {
	// Current
	format.WriteHeader(file, "Current")
	currentOutput := format.ExecuteCommand("uname", "-a")
	format.WritePreformatted(file, currentOutput)

	// Boot options
	format.WriteHeader(file, "Boot options")
	bootOptionsOutput := format.ExecuteCommand("cat", "/proc/cmdline")
	format.WritePreformatted(file, bootOptionsOutput)

	// Modules
	format.WriteHeader(file, "Modules")
	modulesOutput := format.ExecuteCommand("lsmod")
	format.WritePreformatted(file, modulesOutput)

	// Last messages
	//format.WriteHeader(file, "Last messages")
	//lastMessagesOutput := format.ExecuteCommand("dmesg")
	//lastMessages := format.ExecuteCommand("tail", "-n", "50")
	//format.WritePreformatted(file, lastMessagesOutput+lastMessages)
}

func parseLogDate(logLine string) (time.Time, error) {
	// Define a layout matching the log format
	const layout = "Jan 2 15:04:05"
	now := time.Now()
	year := now.Year()

	// Extract the date substring from the log line
	dateStr := logLine[:15] // "Nov 15 08:16:45"
	logTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// Assuming the log is from the current year
	// Adjust if the current time was before the log time to handle year roll-over
	logTime = logTime.AddDate(year-logTime.Year(), 0, 0)
	if logTime.After(now) {
		logTime = logTime.AddDate(-1, 0, 0)
	}

	return logTime, nil
}

func GetErrorKernel(file *os.File) {
	// Заголовок раздела в HTML
	format.WriteHeaderWithID(file, "Kernel Log Errors", "Error")

	// Объявляем регулярное выражение для поиска ключевых слов ошибок
	errorRegex := regexp.MustCompile(`(?i)(error|critical|fail|panic|warn)`)
	yesterday := time.Now().Add(-24 * time.Hour)

	// Функция для обработки лога и записи ошибок в файл за последние сутки
	processLogFile := func(filePath string) {
		logFile, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(file, "<p>Лог %s не найден, пропускаем...</p>\n", filePath)
			return
		}
		defer logFile.Close()

		fmt.Fprintf(file, "<h3>Ошибки за последнии 24ч в %s:</h3>\n<pre>", filePath)
		scanner := bufio.NewScanner(logFile)
		for scanner.Scan() {
			line := scanner.Text()

			// Extract the date from the log
			logTime, err := parseLogDate(line)
			if err != nil {
				fmt.Fprintf(file, "<p>Ошибка парсинга строки: %s: %v</p>\n", line, err)
				continue
			}

			// Filter messages from the last 24 hours
			if logTime.After(yesterday) && errorRegex.MatchString(line) {
				fmt.Fprintln(file, line)
			}
		}
		fmt.Fprintln(file, "</pre>")

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(file, "<p>Ошибка при чтении %s: %v</p>\n", filePath, err)
		}
	}

	// Проверка dmesg и запись ошибок
	fmt.Fprintf(file, "<h3>Ошибки в dmesg:</h3>\n<pre>")
	dmesgOutput := format.ExecuteCommand("dmesg")
	scanner := bufio.NewScanner(strings.NewReader(dmesgOutput))
	for scanner.Scan() {
		line := scanner.Text()
		if errorRegex.MatchString(line) {
			fmt.Fprintln(file, line)
		}
	}
	fmt.Fprintln(file, "</pre>")

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(file, "<p>Ошибка при чтении dmesg: %v</p>\n", err)
	}

	// Обработка файлов /var/log/kern.log, /var/log/kernel.log, /var/log/messages
	processLogFile("/var/log/kern.log")
	processLogFile("/var/log/kernel.log")
	processLogFile("/var/log/messages")
}
