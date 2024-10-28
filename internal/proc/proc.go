package proc

import (
	"fmt"
	"html"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
)

func AddProcessesByCPU(file *os.File) {
	// Заголовок секции
	fmt.Fprintln(file, "<h3>Processes by CPU</h3>")
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
