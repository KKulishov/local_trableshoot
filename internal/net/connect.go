package net

import (
	"fmt"
	"html"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/net"
)

func GetConnections(file *os.File) {
	// TCP sockets
	format.WriteHeader(file, "Listened TCP sockets")
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

// getNetStatSummary возвращает сетевую статистику (аналог `netstat`)
func getNetStatSummary(protocol string) (map[string]int, error) {
	// Получаем список всех сетевых соединений
	conns, err := net.Connections(protocol)
	if err != nil {
		return nil, err
	}

	// Словарь для хранения статистики по состояниям
	stateCount := make(map[string]int)

	// Подсчитываем количество каждого состояния
	for _, conn := range conns {
		state := strings.ToUpper(conn.Status) // Состояние соединения (например, TIME_WAIT, ESTABLISHED)
		stateCount[state]++
	}

	return stateCount, nil
}

func PrintNetStat(file *os.File, protocol string) {
	stats, err := getNetStatSummary(protocol)
	if err != nil {
		fmt.Printf("Error retrieving %s stats: %v\n", protocol, err)
		return
	}

	// Сортируем состояния по количеству соединений
	type kv struct {
		Key   string
		Value int
	}
	var sortedStats []kv
	for k, v := range stats {
		sortedStats = append(sortedStats, kv{k, v})
	}
	sort.Slice(sortedStats, func(i, j int) bool {
		return sortedStats[i].Value < sortedStats[j].Value
	})

	// Добавляем заголовок для секции SYN-SENT
	fmt.Fprintf(file, "<h3 id=\"Network\">Общая сетевая статистика по %s</h3>", protocol)
	fmt.Fprintln(file, "<div><pre>")
	// Выводим статистику
	for _, kv := range sortedStats {
		//fmt.Printf("%d %s\n", kv.Value, kv.Key)
		fmt.Fprintln(file, kv.Value, kv.Key)
	}
	fmt.Fprintln(file, "</pre></div>")
}
