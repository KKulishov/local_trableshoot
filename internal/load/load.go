package load

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/load"
)

func GetLA(file *os.File) {
	loadAvg, err := load.Avg()
	if err != nil {
		fmt.Fprintf(file, "<h3>Трассировка до DNS сервера %s:</h3>\n<pre>", err)
	}

	// Выводим данные
	fmt.Fprintln(file, "<h3>Load Average system</h3>")
	fmt.Fprintln(file, "1 min: ", loadAvg.Load1, " ;")
	fmt.Fprintln(file, "5 mim: ", loadAvg.Load15, " ;")
	fmt.Fprintln(file, "15 mim: ", loadAvg.Load15, " ;")

}
