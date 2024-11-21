package proc

import (
	"fmt"
	"html"
	"local_trableshoot/internal/cgroups"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func AddProcessesByCPU(file *os.File) ([]int, error) {
	// Заголовок секции
	fmt.Fprintln(file, "<h3 id=\"Process\">Processes by CPU</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Выполнение команды для получения списка процессов по CPU
	cmd := exec.Command("sh", "-c", "ps -ewwwo pcpu,pid,user,command --sort -pcpu | head -n 20")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при выполнении команды: %s\n", err)
	} else {
		// Запись результата в файл
		fmt.Fprintln(file, html.EscapeString(string(output)))
	}

	fmt.Fprintln(file, "</pre></div>")

	// Парсинг вывода для извлечения PID
	lines := strings.Split(string(output), "\n")
	var pids []int
	for i, line := range lines {
		// Пропускаем заголовок и пустые строки
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		// Разбиваем строку на поля
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue // Пропускаем строки с недостаточным количеством полей
		}

		// Парсим PID (второе поле)
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue // Пропускаем строки с некорректным PID
		}

		pids = append(pids, pid)
	}

	return pids, nil
}

func SaveContainersToHTML(file *os.File, containerInfos []cgroups.ContainerInfo) error {
	// Начало HTML-страницы и таблицы
	_, err := file.WriteString("<html><body><h3 id=\"Containers\">Container Info by top Pid, cpu used for Pods </h3><table border='1'>")
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

func GetProcessesTree(file *os.File) {
	// Заголовок секции
	fmt.Fprintln(file, "<h3>Processes by Tree</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Выполняем команду ps auxf
	cmd := exec.Command("ps", "auxf")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при выполнении команды: %s\n", err)
	} else {
		// Экранируем специальные символы и записываем результат в файл
		fmt.Fprintln(file, html.EscapeString(string(output)))
	}

	fmt.Fprintln(file, "</pre></div>")
}

func ShowAllCpu(file *os.File) {
	// Show mem linux
	format.WriteHeader(file, "Show all cpu")
	currentOutput := format.ExecuteCommand("nproc")
	format.WritePreformatted(file, currentOutput)
}
