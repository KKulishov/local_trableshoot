package main

import (
	"fmt"
	"local_trableshoot/configs"
	"local_trableshoot/internal/flags"
	"local_trableshoot/internal/hostname"
	"local_trableshoot/internal/platform"
	"local_trableshoot/internal/platform/linux"
	"local_trableshoot/internal/rotate"
	"runtime"
	"time"
)

func main() {
	var diag platform.Diagnostic

	name_host := hostname.HostName()
	currentTime := time.Now().Format("02.01.2006_15:04:05")

	fileNamequick := fmt.Sprintf("%s/report_%s_%s.html", *flags.ReportDir, name_host, currentTime)
	fileName := fmt.Sprintf("%s/full_report_%s_%s.html", *flags.ReportDir, name_host, currentTime)
	// Создаем файл отчета с помощью функции из configs
	file, err := configs.CreateReportFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	filequick, err := configs.CreateReportFile(fileNamequick)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// определяем ос и логику диагнотики
	switch os := runtime.GOOS; os {
	case "linux":
		diag = &linux.LinuxDiagnostic{}
	default:
		fmt.Println("Unsupported platform")
		return
	}

	diag.BaseDiagnostics(filequick)
	diag.FullDiagnostics(file)

	// Очистка старых отчетов
	rotate.CleanUpOldReports(*flags.ReportDir, "report_", *flags.CountRotate)
	rotate.CleanUpOldReports(*flags.ReportDir, "full_report_", *flags.CountRotate)
	fmt.Println("Отчет о процессах создан:", fileNamequick)
	fmt.Println("Отчет о процессах создан:", fileName)
}
