package net

import (
	"bufio"
	"context"
	"fmt"
	"local_trableshoot/internal/flags"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
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

// Проверка доступности DNS-сервера на порту 53
func checkDNSServerAvailability(server string, checkName string) (bool, error) {
	if checkName == "" {
		checkName = "ya.ru." // значение по умолчанию для проверки DNS resolve (полностью квалифицированное имя)
	} else if !strings.HasSuffix(checkName, ".") {
		checkName += "." // добавляем завершающую точку, если ее нет
	}

	// Создаем DNS-клиент
	client := dns.Client{
		Timeout: 2 * time.Second,
	}

	// Создаем DNS-запрос
	msg := new(dns.Msg)
	msg.SetQuestion(checkName, dns.TypeA) // Запрос на получение A-записи (IPv4) для checkName

	// Отправляем запрос
	resp, _, err := client.Exchange(msg, net.JoinHostPort(server, "53"))
	if err != nil {
		return false, fmt.Errorf("ошибка при обращении к DNS-серверу %s: %v", server, err)
	}

	// Проверяем, есть ли ответ в секции ответа
	if len(resp.Answer) == 0 {
		return false, fmt.Errorf("DNS-сервер %s не вернул ответа", server)
	}

	return true, nil
}

// Функция для выполнения tracepath на UDP порту 53
func tracepathToDNSServer(server string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 14*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "tracepath", "-p53", server)
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Sprintf("Таймаут для трассировки до сервера %s\n", server)
	} else if err != nil {
		return fmt.Sprintf("Ошибка выполнения tracepath для сервера %s: %v\n", server, err)
	}
	return string(output)
}

func CheckDnS(file *os.File) {
	fmt.Fprintln(file, "<html><head><title>DNS Check Report</title></head><body>")
	fmt.Fprintln(file, "<h1>Отчет о проверке DNS-серверов</h1>")

	servers, err := getDNSServers()
	if err != nil {
		fmt.Fprintf(file, "<p>Ошибка получения списка DNS-серверов: %v</p>\n", err)
		fmt.Fprintln(file, "</body></html>")
		return
	}

	var wg sync.WaitGroup
	results := make(chan string, len(servers)) // Канал для сбора результатов

	for _, server := range servers {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()

			// Проверяем доступность DNS-сервера
			available, err := checkDNSServerAvailability(server, *flags.CheckNameDns)
			if available {
				results <- fmt.Sprintf("<p>DNS сервер <b>%s</b> доступен</p>\n", server)
			} else {
				if err != nil {
					results <- fmt.Sprintf("<p>DNS сервер <b>%s</b> не доступен: %v</p>\n", server, err)
				}

				// Выполняем трассировку
				trace := tracepathToDNSServer(server)
				results <- fmt.Sprintf("<h3>Трассировка до DNS сервера %s:</h3>\n<pre>%s</pre>", server, trace)
			}
		}(server)
	}

	// Закрытие канала результатов после завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()

	// Сбор и запись результатов
	for result := range results {
		fmt.Fprintln(file, result)
	}

	fmt.Fprintln(file, "</body></html>")
}
