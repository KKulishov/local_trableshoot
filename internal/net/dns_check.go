package net

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Функция для получения списка DNS-серверов из /etc/resolv.conf
func getDNSServers() ([]string, error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл /etc/resolv.conf: %v", err)
	}
	defer file.Close()

	var servers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				servers = append(servers, fields[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла /etc/resolv.conf: %v", err)
	}

	return servers, nil
}

// Функция для проверки доступности DNS-сервера на порту 53
func checkDNSServerAvailability(file *os.File, server string) bool {
	conn, err := net.DialTimeout("udp", net.JoinHostPort(server, "53"), 2*time.Second)
	if err != nil {
		fmt.Fprintf(file, "<p>DNS сервер <b>%s</b> не доступен: %v</p>\n", server, err)
		return false
	}
	defer conn.Close()
	fmt.Fprintf(file, "<p>DNS сервер <b>%s</b> доступен</p>\n", server)
	return true
}

// Функция для выполнения tracepath на UDP порту 53
func tracepathToDNSServer(file *os.File, server string) {
	fmt.Fprintf(file, "<h3>Трассировка до DNS сервера %s:</h3>\n<pre>", server)

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 14*time.Second)
	defer cancel()

	// Выполняем команду с учетом контекста
	cmd := exec.CommandContext(ctx, "tracepath", "-p53", server)
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Fprintf(file, "Таймаут для трассировки до сервера %s\n", server)
	} else if err != nil {
		fmt.Fprintf(file, "Ошибка выполнения tracepath для сервера %s: %v\n", server, err)
	} else {
		fmt.Fprintln(file, string(output))
	}
	fmt.Fprintln(file, "</pre>")
}

func CheckDnS(file *os.File) {
	// Начинаем HTML-отчет
	fmt.Fprintln(file, "<html><head><title>DNS Check Report</title></head><body>")
	fmt.Fprintln(file, "<h1>Отчет о проверке DNS-серверов</h1>")

	// Получаем список DNS-серверов
	servers, err := getDNSServers()
	if err != nil {
		fmt.Fprintf(file, "<p>Ошибка получения списка DNS-серверов: %v</p>\n", err)
	} else {
		// Проверяем каждый сервер на доступность и выполняем tracepath
		for _, server := range servers {
			if !checkDNSServerAvailability(file, server) {
				tracepathToDNSServer(file, server)
			}
		}
	}
	// Заканчиваем HTML-отчет
	fmt.Fprintln(file, "</body></html>")
}
