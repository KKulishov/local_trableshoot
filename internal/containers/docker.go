package containers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type NetworkStats struct {
	Container string
	Name      string
	NetIO     string
}

type DiskUsageStats struct {
	Container string
	Name      string
	Size      int // в байтах для удобства сортировки
	SizeHuman string
}

// Функция для получения статистики использования диска Docker и записи в файл
func GetDockerStatDisk(file *os.File) {
	// Записываем заголовок HTML
	file.WriteString("<html><head><title>Docker Disk Usage Report</title></head><body>\n")
	file.WriteString("<h2>Top 10 Containers by Disk Usage</h2>\n")
	file.WriteString("<table border='1'>\n")
	file.WriteString("<tr><th>Container ID</th><th>Container Name</th><th>Disk Size</th></tr>\n")

	// Выполнение команды docker system df -v
	cmd := exec.Command("docker", "system", "df", "-v")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(file, "Ошибка при выполнении команды:", err)
		return
	}

	// Парсинг вывода команды
	lines := strings.Split(out.String(), "\n")
	var stats []DiskUsageStats
	for _, line := range lines {
		if !strings.Contains(line, "CONTAINER") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		// Конвертируем размер в байты для сортировки
		size := parseSize(parts[3])

		stats = append(stats, DiskUsageStats{
			Container: parts[0],
			Name:      parts[1],
			Size:      size,
			SizeHuman: parts[3],
		})
	}

	// Сортировка по размеру
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Size > stats[j].Size // Сортировка по убыванию
	})

	// Ограничиваем до 10 лучших
	if len(stats) > 10 {
		stats = stats[:10]
	}

	// Записываем статистику в файл
	for _, stat := range stats {
		row := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n", stat.Container, stat.Name, stat.SizeHuman)
		file.WriteString(row)
	}

	// Закрываем таблицу и HTML
	file.WriteString("</table>\n")
	file.WriteString("</body></html>\n")
}

// Функция для получения статистики сети Docker и записи в файл
func GetDockerStatNetwork(file *os.File) {
	// Записываем заголовок HTML
	file.WriteString("<html><head><title>Docker Network Stats Report</title></head><body>\n")
	file.WriteString("<h2>Top 10 Containers by Network Usage</h2>\n")
	file.WriteString("<table border='1'>\n")
	file.WriteString("<tr><th>Container ID</th><th>Container Name</th><th>Network I/O</th></tr>\n")

	// Выполнение команды docker stats
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.Container}},{{.Name}},{{.NetIO}}")

	// Получение вывода команды
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(file, "Ошибка при выполнении команды:", err)
		return
	}

	// Разбиваем вывод на строки
	lines := strings.Split(out.String(), "\n")

	// Список для хранения статистики сети контейнеров
	var stats []NetworkStats

	// Обработка строк вывода
	for _, line := range lines {
		if line == "" {
			continue // Игнорируем пустые строки
		}

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue // Игнорируем некорректные строки
		}

		stats = append(stats, NetworkStats{
			Container: parts[0],
			Name:      parts[1],
			NetIO:     parts[2],
		})
	}

	// Сортировка по сетевому трафику
	sort.Slice(stats, func(i, j int) bool {
		// Получаем байты входящего и исходящего трафика для сравнения
		inI, outI := parseNetIO(stats[i].NetIO)
		inJ, outJ := parseNetIO(stats[j].NetIO)
		return (inI + outI) > (inJ + outJ)
	})

	// Ограничиваем до 10 лучших
	if len(stats) > 10 {
		stats = stats[:10]
	}

	// Записываем статистику в файл
	for _, stat := range stats {
		row := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n", stat.Container, stat.Name, stat.NetIO)
		file.WriteString(row)
	}

	// Закрываем таблицу и HTML
	file.WriteString("</table>\n")
	file.WriteString("</body></html>\n")
}

// parseNetIO разбирает NetIO в байты для входящего и исходящего трафика
func parseNetIO(netIO string) (inBytes, outBytes int) {
	parts := strings.Split(netIO, " / ")
	if len(parts) != 2 {
		return 0, 0
	}
	inBytes = parseSize(parts[0])
	outBytes = parseSize(parts[1])
	return
}

// parseSize парсит строку размера в байты (например, "10MB" в 10485760)
func parseSize(sizeStr string) int {
	sizeStr = strings.TrimSpace(sizeStr)
	var multiplier int
	if strings.HasSuffix(sizeStr, "kB") {
		multiplier = 1024
		sizeStr = strings.TrimSuffix(sizeStr, "kB")
	} else if strings.HasSuffix(sizeStr, "MB") {
		multiplier = 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "MB")
	} else if strings.HasSuffix(sizeStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = strings.TrimSuffix(sizeStr, "GB")
	} else {
		return 0
	}

	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return 0
	}

	return int(size * float64(multiplier))
}
