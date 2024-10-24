package main

import (
	"fmt"
	"local_trableshoot/configs"
	"local_trableshoot/internal/containers"
	"local_trableshoot/internal/disk"
	"local_trableshoot/internal/flags"
	"local_trableshoot/internal/hostname"
	"local_trableshoot/internal/kernel"
	"local_trableshoot/internal/mem"
	"local_trableshoot/internal/net"
	"local_trableshoot/internal/proc"
	"local_trableshoot/internal/rotate"
	"local_trableshoot/internal/top"
	"time"
)

func main() {
	name_host := hostname.HostName()
	currentTime := time.Now().Format("02.01.2006_15:04:05")
	// /var/log or /tmp
	fileName := fmt.Sprintf("/var/log/report_%s_%s.html", name_host, currentTime)

	// Создаем файл отчета с помощью функции из configs
	file, err := configs.CreateReportFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// TODO parallel run
	// Создаем wait group для синхронизации горутин
	//var wg sync.WaitGroup
	// Добавляем список процессов&памяти в HTML

	hostname.GetHostName(file)
	top.GetSummary(file)
	top.Get_atop_processes_lists(file)
	if *flags.ContainerFlag == "docker" {
		containers.GetDockerStatCpu(file)
		containers.GetDockerStatMem(file)
	}
	proc.AddProcessesByCPU(file)
	proc.GetProcessesTree(file)
	mem.AddProcessesByMem(file)
	net.GetConnections(file)
	net.GetNetworkStats(file)
	net.GetTrableConnections(file)
	if *flags.CheckDns {
		net.CheckDnS(file)
	}
	disk.GetDisksInfo(file)
	kernel.GetKernelAndModules(file)

	// top network traffic used process
	// tcpdump
	// arp
	// tracert до не стабильнго соединения
	// ToDo add lsof
	// ToDo add ping& traceroute
	// try upload to s3
	//wg.Wait()

	// Очистка старых отчетов
	rotate.CleanUpOldReports("/var/log", "report_", *flags.CountRotate)
	fmt.Println("Отчет о процессах создан:", fileName)
}
