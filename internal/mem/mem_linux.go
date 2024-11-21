package mem

import (
	"fmt"
	"local_trableshoot/internal/cgroups"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
