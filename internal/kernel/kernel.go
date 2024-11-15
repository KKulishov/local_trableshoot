package kernel

import (
	"bufio"
	"fmt"
	"local_trableshoot/internal/format"
	"os"
	"regexp"
	"strings"
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

func GetErrorKernel(file *os.File) {
	// Заголовок раздела в HTML
	format.WriteHeader(file, "Kernel Log Errors")

	// Объявляем регулярное выражение для поиска ключевых слов ошибок
	errorRegex := regexp.MustCompile(`(?i)(error|critical|fail|panic|warn)`)

	// Функция для обработки лога и записи ошибок в файл
	processLogFile := func(filePath string) {
		logFile, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(file, "<p>Лог %s не найден, пропускаем...</p>\n", filePath)
			return
		}
		defer logFile.Close()

		fmt.Fprintf(file, "<h3>Ошибки в %s:</h3>\n<pre>", filePath)
		scanner := bufio.NewScanner(logFile)
		for scanner.Scan() {
			line := scanner.Text()
			if errorRegex.MatchString(line) {
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

	// Обработка файлов /var/log/kernel.log и /var/log/messages
	processLogFile("/var/log/kern.log")
	processLogFile("/var/log/kernel.log")
	processLogFile("/var/log/messages")
}
