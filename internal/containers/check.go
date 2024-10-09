package containers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type ContainerStats struct {
	Container string
	Name      string
	CPUPerc   string
	MemUsage  string
}

// ToDo check CRI

// Функция для получения статистики CPU Docker и записи в файл
func GetDockerStatCpu(file *os.File) {
	// Записываем заголовок HTML
	file.WriteString("<html><head><title>Docker Stats Report</title></head><body>\n")
	file.WriteString("<h2>Top 10 Containers by CPU Usage</h2>\n")
	file.WriteString("<table border='1'>\n")
	file.WriteString("<tr><th>Container ID</th><th>Container Name</th><th>CPU %</th></tr>\n")

	// Выполнение команды docker stats
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.Container}},{{.Name}},{{.CPUPerc}}")

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

	// Список для хранения статистики контейнеров
	var stats []ContainerStats

	// Обработка строк вывода
	for _, line := range lines {
		if line == "" {
			continue // Игнорируем пустые строки
		}

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue // Игнорируем некорректные строки
		}

		stats = append(stats, ContainerStats{
			Container: parts[0],
			Name:      parts[1],
			CPUPerc:   parts[2],
		})
	}

	// Сортировка по CPUPerc
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].CPUPerc > stats[j].CPUPerc // Сортировка по убыванию
	})

	// Ограничиваем до 10 лучших
	if len(stats) > 10 {
		stats = stats[:10]
	}

	// Записываем статистику в файл
	for _, stat := range stats {
		row := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n", stat.Container, stat.Name, stat.CPUPerc)
		file.WriteString(row)
	}

	// Закрываем таблицу и HTML
	file.WriteString("</table>\n")
	file.WriteString("</body></html>\n")
}

// GetDockerStatMem получает топ-10 контейнеров по использованию памяти и записывает в файл
func GetDockerStatMem(file *os.File) {
	// Записываем заголовок HTML
	file.WriteString("<html><head><title>Docker Stats Report</title></head><body>\n")
	file.WriteString("<h2>Top 10 Containers by MEM Usage</h2>\n")
	file.WriteString("<table border='1'>\n")
	file.WriteString("<tr><th>Container ID</th><th>Container Name</th><th>MEM %</th></tr>\n")

	// Выполнение команды docker stats
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{.Container}},{{.Name}},{{.MemUsage}}")

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

	// Список для хранения статистики контейнеров
	var stats []ContainerStats

	// Обработка строк вывода
	for _, line := range lines {
		if line == "" {
			continue // Игнорируем пустые строки
		}

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue // Игнорируем некорректные строки
		}

		stats = append(stats, ContainerStats{
			Container: parts[0],
			Name:      parts[1],
			MemUsage:  parts[2],
		})
	}

	// Сортировка по CPUPerc
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].MemUsage > stats[j].MemUsage // Сортировка по убыванию
	})

	// Ограничиваем до 10 лучших
	if len(stats) > 10 {
		stats = stats[:10]
	}

	// Записываем статистику в файл
	for _, stat := range stats {
		row := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>\n", stat.Container, stat.Name, stat.MemUsage)
		file.WriteString(row)
	}

	// Закрываем таблицу и HTML
	file.WriteString("</table>\n")
	file.WriteString("</body></html>\n")

}
