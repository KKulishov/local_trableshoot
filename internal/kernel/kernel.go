package kernel

import (
	"local_trableshoot/internal/format"
	"os"
)

func GetKernelAndModules(file *os.File) {
	// Current
	format.WriteHeader(file, "Current")
	currentOutput := format.ExecuteCommand("uname", "-a")
	format.WritePreformatted(file, currentOutput)

	// Boot options
	format.WriteHeader(file, "Boot options")
	bootOptionsOutput := format.ExecuteCommand("cat", "/proc/cmdline")
	format.WritePreformatted(file, bootOptionsOutput)

	// Modules
	format.WriteHeader(file, "Modules")
	modulesOutput := format.ExecuteCommand("lsmod")
	format.WritePreformatted(file, modulesOutput)

	// Last messages
	format.WriteHeader(file, "Last messages")
	lastMessagesOutput := format.ExecuteCommand("dmesg")
	lastMessages := format.ExecuteCommand("tail", "-n", "50")
	format.WritePreformatted(file, lastMessagesOutput+lastMessages)
}
