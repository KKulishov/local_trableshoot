package main

import (
	"fmt"
	"local_trableshoot/configs"
	"local_trableshoot/internal/flags"
	"local_trableshoot/internal/hostname"
	"local_trableshoot/internal/net"
	"local_trableshoot/internal/platform"
	"local_trableshoot/internal/platform/linux"
	"local_trableshoot/internal/rotate"
	"local_trableshoot/internal/s3"
	"runtime"
	"sync"
	"time"
)

func main() {
	var diag platform.Diagnostic

	// проверка существ. каталога для отчетов , если его нет то создаем
	configs.CheckAndCreateDir(*flags.ReportDir)

	name_host := hostname.HostName()
	currentTime := time.Now().Format("02.01.2006_15:04:05")

	fileNamequick := fmt.Sprintf("%s/report_%s_%s.html", *flags.ReportDir, name_host, currentTime)
	fileName := fmt.Sprintf("%s/full_report_%s_%s.html", *flags.ReportDir, name_host, currentTime)
	fileNameNetwork := fmt.Sprintf("%s/network_report_%s_%s.html", *flags.ReportDir, name_host, currentTime)
	// Создаем файл отчета с помощью функции из configs
	file, err := configs.CreateReportFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	filenetwork, err := configs.CreateReportFile(fileNameNetwork)
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

	// если задан RunRotateS3 аргумент при запуске то идет чистка S3 бакета
	if *flags.RunRotateS3 {
		s3.Rotation_s3_bucket(fileName, *flags.CountRotate_S3)
		s3.Rotation_s3_bucket(fileNamequick, *flags.CountRotate_S3)
		s3.Rotation_s3_bucket(fileNameNetwork, *flags.CountRotate_S3)
		return
	}

	// ToDo add в отдельной горутине
	if *flags.NetworkAnalyze {
		diag.NetowrDiagnosics(filenetwork)
		perfReportPath := net.PerfAnalyzSoftirqd(filenetwork)
		dumpReportPath := net.TcpDumpAnalyze(filenetwork)
		fmt.Println("Отчет о процессах создан:", fileNameNetwork)
		s3.Send_report_file(fileNameNetwork)
		s3.Send_report_file(perfReportPath)
		s3.Send_report_file(dumpReportPath)
	}
	//

	// Используем sync.WaitGroup для синхронизации
	var wg sync.WaitGroup

	// Каналы для синхронизации зависимостей
	baseDiagnosticsDone := make(chan struct{})
	fullDiagnosticsDone := make(chan struct{})

	// Запуск BaseDiagnostics в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		diag.BaseDiagnostics(filequick)
		close(baseDiagnosticsDone) // Сигнализируем, что BaseDiagnostics завершен
	}()
	// Запуск FullDiagnostics в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		//<-baseDiagnosticsDone // Ждем завершения BaseDiagnostics
		diag.FullDiagnostics(file)
		close(fullDiagnosticsDone) // Сигнализируем, что FullDiagnostics завершен
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-fullDiagnosticsDone // Ждем завершения BaseDiagnostics
		rotate.CleanUpOldReports(*flags.ReportDir, "report_", *flags.CountRotate)
		rotate.CleanUpOldReports(*flags.ReportDir, "full_report_", *flags.CountRotate)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-baseDiagnosticsDone // Ждем завершения BaseDiagnostics
		s3.Send_report_file(fileNamequick)
		fmt.Println("Отчет о процессах создан:", fileNamequick)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-fullDiagnosticsDone // Ждем завершения BaseDiagnostics
		s3.Send_report_file(fileName)
		fmt.Println("Отчет о процессах создан:", fileName)
	}()
	// Ожидаем завершения всех горутин
	wg.Wait()
	//diag.BaseDiagnostics(filequick)
	//diag.FullDiagnostics(file)
	// Очистка старых отчетов
	//rotate.CleanUpOldReports(*flags.ReportDir, "report_", *flags.CountRotate)
	//rotate.CleanUpOldReports(*flags.ReportDir, "full_report_", *flags.CountRotate)
	//fmt.Println("Отчет о процессах создан:", fileNamequick)
	//fmt.Println("Отчет о процессах создан:", fileName)
	// Загружаем конфигурацию из файла для s3
	//s3.Send_report_file(fileName)
	//s3.Send_report_file(fileNamequick)

}
