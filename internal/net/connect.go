package net

import (
	"fmt"
	"html"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
	"strings"
)

func GetConnections(file *os.File) {
	// TCP sockets
	format.WriteHeaderWithID(file, "Listened TCP sockets", "Network")
	tcpOutput := format.ExecuteCommand("ss", "-nlpt")
	format.WritePreformatted(file, tcpOutput)

	// UDP sockets
	format.WriteHeader(file, "Listened UDP sockets")
	udpOutput := format.ExecuteCommand("ss", "-nlpu")
	format.WritePreformatted(file, udpOutput)

	// Unix sockets
	format.WriteHeader(file, "Listened Unix sockets")
	unixOutput := format.ExecuteCommand("ss", "-nlpu")
	format.WritePreformatted(file, unixOutput)

	// Connections by IP
	format.WriteHeader(file, "Connections by IP")
	connectionsOutput := format.ExecuteCommand("ss", "-anlpe", "-A", "inet")
	ips := strings.Split(connectionsOutput, "\n")
	ipFrequency := make(map[string]int)

	for _, line := range ips[2:] {
		fields := strings.Fields(line)
		if len(fields) > 4 {
			ip := strings.Split(fields[4], ":")[0]
			ipFrequency[ip]++
		}
	}

	for ip, count := range ipFrequency {
		format.WritePreformatted(file, fmt.Sprintf("%s: %d", ip, count))
	}
}

// AddSocketStates получает список сокетов в состояниях TIME_WAIT и SYN_SENT и записывает их в HTML.
func GetTrableConnections(file *os.File) {
	// Добавляем заголовок для секции TIME_WAIT
	fmt.Fprintln(file, "<h3>Sockets in TIME-WAIT</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Выполняем команду ss для TIME_WAIT
	cmdTimeWait := exec.Command("sh", "-c", "ss -tuanp | grep TIME-WAIT")
	outputTimeWait, err := cmdTimeWait.CombinedOutput()
	if err != nil {
		// fmt.Fprintf(file, "TIME-WAIT не обнаружены: %s\n", err)
		fmt.Fprintf(file, "TIME-WAIT не обнаружены")
	} else {
		// Экранируем специальные символы и записываем результат в файл
		fmt.Fprintln(file, html.EscapeString(string(outputTimeWait)))
	}
	fmt.Fprintln(file, "</pre></div>")

	// Добавляем заголовок для секции SYN-SENT
	fmt.Fprintln(file, "<h3>Sockets in SYN-SENT</h3>")
	fmt.Fprintln(file, "<div><pre>")
	fmt.Fprint(file, `SYN-SENT Если это состояние долго не меняется, скорее всего, соединение не может быть установлено.
	Обычно это означает, что сервер не отвечает, или что есть сетевые проблемы между клиентом и сервером.`)
	fmt.Fprintln(file, "<div><pre>")
	// Выполняем команду ss для SYN_SENT
	cmdSynSent := exec.Command("sh", "-c", "ss -tuanp | grep SYN-SENT")
	outputSynSent, err := cmdSynSent.CombinedOutput()
	if err != nil {
		// fmt.Fprintf(file, "SYN-SENT не обнаружены: %s\n", err)
		fmt.Fprintf(file, "SYN-SENT не обнаружены\n")
	} else {
		// Экранируем специальные символы и записываем результат в файл
		fmt.Fprintln(file, html.EscapeString(string(outputSynSent)))
	}
	// Добавляем заголовок для секции SYN-RECEIVED:
	fmt.Fprintln(file, "<h3>Sockets in SYN-RECEIVED</h3>")
	fmt.Fprintln(file, "<div><pre>")
	fmt.Fprint(file, `SYN-RECEIVED Если это состояние длится долго, это может означать, что клиент недоступен или отклонил соединение, 
	либо есть сетевые проблемы, препятствующие ответу клиента.`)
	fmt.Fprintln(file, "<div><pre>")
	// Выполняем команду ss для SYN-RECEIVED
	cmdSynRECEIVED := exec.Command("sh", "-c", "ss -tuanp | grep SYN-RECEIVED")
	outputSynRECEIVED, err := cmdSynRECEIVED.CombinedOutput()
	if err != nil {
		fmt.Fprintf(file, "SYN-RECEIVED не обнаружены: %s\n", err)
	} else {
		// Экранируем специальные символы и записываем результат в файл
		fmt.Fprintln(file, html.EscapeString(string(outputSynRECEIVED)))
	}

	fmt.Fprintln(file, "</pre></div>")
}
