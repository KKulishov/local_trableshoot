package net

import (
	"bufio"
	"context"
	"fmt"
	"local_trableshoot/internal/flags"
	"local_trableshoot/internal/format"
	"local_trableshoot/internal/hostname"
	"local_trableshoot/internal/rotate"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ToDo analyze network utilization for pid/command
// add cilium/ebpf or go-tcpdump
func isCommandAvailable(command string) bool {
	cmd := exec.Command("which", command)
	output, err := cmd.CombinedOutput()
	return err == nil && len(strings.TrimSpace(string(output))) > 0
}

func TrablNetBase(file *os.File) {
	format.WriteHeader(file, "Ifconfig")
	ifconfig := format.ExecuteCommand("ifconfig", "-a")
	format.WritePreformatted(file, ifconfig)

	format.WriteHeader(file, "Net dev")
	net_dev := format.ExecuteCommand("cat", "/proc/net/dev")
	format.WritePreformatted(file, net_dev)

	format.WriteHeader(file, "Net protocols")
	net_protocols := format.ExecuteCommand("cat", "/proc/net/protocols")
	format.WritePreformatted(file, net_protocols)

	format.WriteHeader(file, "Slabtop")
	slabtop := format.ExecuteCommand("slabtop", "-o")
	format.WritePreformatted(file, slabtop)

	format.WriteHeader(file, "Netstat status")
	netstat := format.ExecuteCommand("netstat", "-tulpn")
	format.WritePreformatted(file, netstat)
}

func AnalyzeSoftirqdWithPS(file *os.File) {
	format.WriteHeader(file, "Ksoftirqd cpu used")
	cmd := exec.Command("sh", "-c", "ps aux --sort -pcpu | grep ksoftirqd | head -n 1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(file, "Ошибка выполнения ps: %v\n", err)
		return
	}

	fmt.Fprintln(file, "Результаты анализа через ps:")
	fmt.Fprintln(file, string(output))
}

// анализ ksoftirqd, то /proc/interrupts может дать представление о частоте прерываний.
func AnalyzeInterrupts(file *os.File) {
	format.WriteHeader(file, "Interrupts")
	interruptsPath := "/proc/interrupts"

	// Открываем файл
	interruptsFile, err := os.Open(interruptsPath)
	if err != nil {
		fmt.Fprintf(file, "Ошибка открытия файла %s: %v\n", interruptsPath, err)
		return
	}
	defer interruptsFile.Close()

	// Читаем файл построчно
	scanner := bufio.NewScanner(interruptsFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "IR-") || strings.Contains(line, "SOFT") {
			fmt.Fprintln(file, line)
		}
	}
}

// PerfAnalyzSoftirqd анализирует самый нагруженный процесс ksoftirqd с использованием perf
func PerfAnalyzSoftirqd(file *os.File) string {
	if !isCommandAvailable("perf") {
		fmt.Fprintln(file, "<h3>Утилита perf не найдена</h3>")
		fmt.Fprintln(file, "<div><pre>Убедитесь, что perf установлена в системе.</pre></div>")
		return "Утилита perf не найдена"
	}
	// Заголовок секции в отчете
	fmt.Fprintln(file, "<h3 id=\"PerfAnalyzSoftirqd\">Softirqd Performance Analysis</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Находим PID самого нагруженного процесса ksoftirqd
	cmd := exec.Command("sh", "-c", "ps aux --sort -pcpu | grep ksoftirqd | head -n 1 | awk '{print $2}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при выполнении команды для поиска ksoftirqd: %s\n", err)
		fmt.Fprintln(file, "</pre></div>")
		//return
	}

	// Получаем PID из команды
	pid := strings.TrimSpace(string(output))
	if pid == "" {
		fmt.Fprintln(file, "Не найден нагруженный процесс ksoftirqd.")
		fmt.Fprintln(file, "</pre></div>")
		//return
	}

	fmt.Fprintf(file, "Самый нагруженный процесс ksoftirqd: PID=%s\n", pid)

	name_host := hostname.HostName()
	currentTime := time.Now().Format("02.01.2006_15:04:05")
	perf_report := fmt.Sprintf("%s/perf_%s_%s.data", *flags.ReportDir, name_host, currentTime)
	// Формируем команду для анализа perf
	perfCmd := fmt.Sprintf("perf record -o %s -F 99 -a -g -c 3 -p %s -- sleep 5", perf_report, pid)
	cmd = exec.Command("sh", "-c", perfCmd)
	// Запускаем команду
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при выполнении команды perf: %s\n", err)
		fmt.Fprintln(file, "</pre></div>")
		//return
	}
	// Записываем результаты в отчет
	fmt.Fprintln(file, "Результаты perf анализа:")
	fmt.Fprintln(file, "Результаты perf анализа:", perf_report)
	fmt.Fprintln(file, "HowTo analyze perf report")
	fmt.Fprintln(file, "perf report -i ", perf_report)
	fmt.Fprintln(file, "perf report -i ", perf_report, " --stdio")
	fmt.Fprintln(file, "If yuo need interactive perf pid trablshoot, you can run: perf top -g --no-children -F 99 -c 3 -p {PID}")
	fmt.Fprintln(file, "</pre></div>")
	fmt.Fprintln(file, "<p>Анализ завершен.</p>")
	rotate.CleanUpOldReports(*flags.ReportDir, "perf_", *flags.CountRotate)
	return string(perf_report)
}

// TcpDumpAnalyze запускает tcpdump с ограничением по времени и записывает результаты в файл
func TcpDumpAnalyze(file *os.File) string {
	if !isCommandAvailable("tcpdump") {
		fmt.Fprintln(file, "<h3>Утилита tcpdump не найдена</h3>")
		fmt.Fprintln(file, "<div><pre>Убедитесь, что tcpdump установлена в системе.</pre></div>")
		return "Утилита tcpdump не найдена"
	}
	// Заголовок секции в отчете
	fmt.Fprintln(file, "<h3 id=\"TcpDumpAnalyze\">TCP Dump Analysis</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Получаем имя хоста
	name_host, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при получении имени хоста: %s\n", err)
		fmt.Fprintln(file, "</pre></div>")
		return ""
	}

	// Форматируем текущее время
	currentTime := time.Now().Format("02.01.2006_15:04:05")
	// Устанавливаем путь к файлу отчета
	dump_report := fmt.Sprintf("%s/tcpdump_%s_%s.pcap", *flags.ReportDir, name_host, currentTime)

	// Формируем команду для tcpdump
	tcpdumpCmd := exec.Command("tcpdump", "-i", "any", "-nnneeee", "-w", dump_report)

	// Создаем контекст с тайм-аутом 6 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Связываем контекст с командой
	tcpdumpCmd = exec.CommandContext(ctx, "tcpdump", "-i", "any", "-nnneeee", "-w", dump_report)

	// Запускаем команду
	err = tcpdumpCmd.Start()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при запуске tcpdump: %s\n", err)
		fmt.Fprintln(file, "</pre></div>")
		return ""
	}

	// Ожидаем завершения команды
	err = tcpdumpCmd.Wait()

	// Проверяем, завершилась ли команда из-за тайм-аута
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Fprintln(file, "tcpdump завершен из-за превышения времени (6 секунд).")
	} else if err != nil {
		fmt.Fprintf(file, "Ошибка при выполнении tcpdump: %s\n", err)
	}

	// Записываем результаты в отчет
	fmt.Fprintln(file, "Результаты tcpdump анализа:")
	fmt.Fprintf(file, "Файл дампа: %s\n", dump_report)
	fmt.Fprintln(file, "Help man used:")
	fmt.Fprintln(file, "tcpdump -r ", dump_report)
	fmt.Fprintln(file, "Фильтрация по IP-адресу:  tcpdump -r ", dump_report, " host 192.168.1.1")
	fmt.Fprintln(file, "Фильтрация по порту:  tcpdump -r ", dump_report, " port 80")
	fmt.Fprintln(file, "Фильтрация по протоколу::  tcpdump -r ", dump_report, " tcp")
	fmt.Fprintln(file, "<p>Анализ завершен.</p>")

	rotate.CleanUpOldReports(*flags.ReportDir, "tcpdump_", 2)

	return dump_report
}

/*
func TcpDumpAnalyze(file *os.File) string {
	// Заголовок секции в отчете
	fmt.Fprintln(file, "<h3 id=\"TcpDumpAnalyze\">TCP Dump Analysis</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Получаем имя хоста
	name_host, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при получении имени хоста: %s\n", err)
		fmt.Fprintln(file, "</pre></div>")
		return ""
	}

	// Форматируем текущее время
	currentTime := time.Now().Format("02.01.2006_15:04:05")
	// Устанавливаем путь к файлу отчета
	dump_report := fmt.Sprintf("%s/tcpdump_%s_%s.pcap", *flags.ReportDir, name_host, currentTime)

	// Формируем команду для tcpdump
	tcpdumpCmd := fmt.Sprintf("tcpdump -i any -nnneeee -w %s", dump_report)

	// Запускаем команду
	cmd := exec.Command("sh", "-c", tcpdumpCmd)
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(file, "Ошибка при запуске tcpdump: %s\n", err)
		fmt.Fprintln(file, "</pre></div>")
		return ""
	}
	// Записываем результаты в отчет
	fmt.Fprintln(file, "Результаты tcpdump анализа:")
	fmt.Fprintln(file, "Файл дампа:", dump_report)
	fmt.Fprintln(file, "</pre></div>")
	fmt.Fprintln(file, "<p>Анализ завершен.</p>")

	return dump_report
}
*/
