//go:build linux
// +build linux

// internal/platform/linux/linux.go
package linux

import (
	"local_trableshoot/internal/containers"
	"local_trableshoot/internal/disk"
	"local_trableshoot/internal/flags"
	"local_trableshoot/internal/format"
	"local_trableshoot/internal/hostname"
	"local_trableshoot/internal/kernel"
	"local_trableshoot/internal/load"
	"local_trableshoot/internal/mem"
	"local_trableshoot/internal/net"
	"local_trableshoot/internal/proc"
	"local_trableshoot/internal/top"
	"os"
)

type LinuxDiagnostic struct{}

// TODO parallel run
// Создаем wait group для синхронизации горутин
//var wg sync.WaitGroup
// Добавляем список процессов&памяти в HTM

func (d *LinuxDiagnostic) FullDiagnostics(file *os.File) {
	format.WriteHTMLHeader(file)

	hostname.GetVersionApp(file)
	hostname.GetHostName(file)
	proc.ShowAllCpu(file)
	load.GetLA(file)
	mem.ShowMem(file)
	top.GetSummary(file)
	top.Get_atop_processes_lists(file)
	if *flags.ContainerFlag == "docker" {
		containers.GetDockerStatCpu(file)
		containers.GetDockerStatMem(file)
		containers.GetDockerStatDisk(file)
		containers.GetDockerStatNetwork(file)
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
	disk.AppDiskUtilization(file)
	disk.GetDisksInfo(file)
	kernel.GetErrorKernel(file)
	kernel.GetKernelAndModules(file)

	format.WriteHTMLFooter(file)
}

func (d *LinuxDiagnostic) BaseDiagnostics(file *os.File) {
	format.WriteHTMLHeader(file)

	hostname.GetVersionApp(file)
	hostname.GetHostName(file)
	proc.ShowAllCpu(file)
	load.GetLA(file)
	mem.ShowMem(file)
	if *flags.ContainerFlag == "docker" {
		containers.GetDockerStatCpu(file)
		containers.GetDockerStatMem(file)
		containers.GetDockerStatDisk(file)
		containers.GetDockerStatNetwork(file)
	}
	proc.AddProcessesByCPU(file)
	mem.AddProcessesByMem(file)
	top.Get_atop_processes_lists(file)
	kernel.GetErrorKernel(file)

	format.WriteHTMLFooter(file)
}

// top network traffic used process
// tcpdump
// arp
// tracert до не стабильнго соединения
// ToDo add lsof
// ToDo add ping& traceroute
// ToDo strace&perf
//wg.Wait()
