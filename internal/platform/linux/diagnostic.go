//go:build linux
// +build linux

// internal/platform/linux/linux.go
package linux

import (
	"fmt"
	"local_trableshoot/internal/cgroups"
	"local_trableshoot/internal/containers"
	"local_trableshoot/internal/disk"
	"local_trableshoot/internal/flags"
	"local_trableshoot/internal/format"
	"local_trableshoot/internal/hostname"
	"local_trableshoot/internal/kernel"
	"local_trableshoot/internal/load"
	"local_trableshoot/internal/mem"
	"local_trableshoot/internal/net"
	"local_trableshoot/internal/perfomance"
	"local_trableshoot/internal/proc"
	"local_trableshoot/internal/top"
	"os"
)

type LinuxDiagnostic struct{}

// TODO parallel run
// Создаем wait group для синхронизации горутин
//var wg sync.WaitGroup
// Добавляем список процессов&памяти в HTM

var dir_kubelet string = "/var/log/containers"

func (d *LinuxDiagnostic) FullDiagnostics(file *os.File) {
	format.WriteHTMLHeader(file)
	format.ListAnchorReport(file)

	hostname.GetVersionApp(file)
	hostname.GetHostName(file)
	proc.ShowAllCpu(file)
	load.GetLA(file)
	mem.ShowMem(file)
	top.GetSummary(file)
	if *flags.AtopReport {
		top.Get_atop_processes_lists(file)
	}
	if *flags.ContainerFlag == "docker" {
		containers.GetDockerStatCpu(file)
		containers.GetDockerStatMem(file)
		//containers.GetDockerStatDisk(file)
		//containers.GetDockerStatNetwork(file)
	}

	getProcess_to_ns(file, dir_kubelet)
	getMem_to_ns(file, dir_kubelet)
	mem.ShowProcessesSwapUsage(file)

	proc.GetProcessesTree(file)
	net.PrintNetStat(file, "tcp")
	net.PrintNetStat(file, "udp")
	net.GetConnections(file)
	net.GetNetworkStats(file)
	net.GetTrableConnections(file)
	if *flags.CheckDns {
		net.CheckDnS(file)
	}
	disk.AppDiskUtilization(file)
	disk.GetDisksInfo(file)
	kernel.GetErrorKernel(file)
	kernel.GetKernelAndModules(file)

	format.WriteHTMLFooter(file)
}

func (d *LinuxDiagnostic) BaseDiagnostics(file *os.File) {
	format.WriteHTMLHeader(file)
	format.ListAnchorReport(file)

	hostname.GetVersionApp(file)
	hostname.GetHostName(file)
	proc.ShowAllCpu(file)
	load.GetLA(file)
	mem.ShowMem(file)
	perfomance.RunCpuResults(file)
	if *flags.ContainerFlag == "docker" {
		containers.GetDockerStatCpu(file)
		containers.GetDockerStatMem(file)
		//containers.GetDockerStatDisk(file)
		//containers.GetDockerStatNetwork(file)
	}
	getProcess_to_ns(file, dir_kubelet)
	getMem_to_ns(file, dir_kubelet)
	mem.ShowProcessesSwapUsage(file)

	if *flags.AtopReport {
		top.Get_atop_processes_lists(file)
	}
	net.PrintNetStat(file, "tcp")
	net.PrintNetStat(file, "udp")
	kernel.GetErrorKernel(file)

	format.WriteHTMLFooter(file)
}

func (d *LinuxDiagnostic) NetowrDiagnosics(file *os.File) {
	format.WriteHTMLHeader(file)

	net.TrablNetBase(file)
	net.AnalyzeSoftirqdWithPS(file)
	net.TcpDumpAnalyze(file)
	//net.AnalyzeInterrupts(file)
	//net.PerfAnalyzSoftirqd(file)

	format.WriteHTMLFooter(file)
}

func getProcess_to_ns(file *os.File, kuber_dir string) {
	pids, err := proc.AddProcessesByCPU(file)
	if err != nil {
		fmt.Println("Error getting processes by CPU:", err)
		return
	}

	_, err_dir_kubelet := os.Stat(kuber_dir)

	if err_dir_kubelet == nil {
		// Обрабатываем полученные PID
		containerInfos, err := cgroups.ProcessPIDs(pids)
		if err != nil {
			//fmt.Println("Error processing PIDs:", err)
			return
		}
		// Выводим информацию о контейнерах
		proc.SaveContainersToHTML(file, containerInfos)
	}
}

func getMem_to_ns(file *os.File, kuber_dir string) {
	pids, err := mem.AddProcessesByMem(file)
	if err != nil {
		fmt.Println("Error getting processes by MEM:", err)
		return
	}
	_, err_dir_kubelet := os.Stat(kuber_dir)
	if err_dir_kubelet == nil {
		// Обрабатываем полученные PID
		containerInfos, err := cgroups.ProcessPIDs(pids)
		if err != nil {
			//fmt.Println("Error processing PIDs:", err)
			return
		}
		// Выводим информацию о контейнерах
		mem.SaveContainersToHTML(file, containerInfos)
	}

}

// top network traffic used process
// tracert до не стабильнго соединения
// ToDo add lsof
// ToDo add ping& traceroute
// ToDo strace&perf
//wg.Wait()
