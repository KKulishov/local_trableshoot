package mem

import (
	"fmt"
	"os"
	"os/exec"
)

func AddProcessesByMem(file *os.File) {
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
}
