package disk

import (
	"local_trableshoot/internal/format"
	"os"
)

func GetDisksInfo(file *os.File) {
	// Devices
	format.WriteHeader(file, "Devices")
	devicesOutput := format.ExecuteCommand("lsblk")
	format.WritePreformatted(file, devicesOutput)

	// Space Usage
	format.WriteHeader(file, "Space usage")
	spaceUsageOutput := format.ExecuteCommand("df", "-h")
	format.WritePreformatted(file, spaceUsageOutput)

	// Inodes Usage
	format.WriteHeader(file, "Inodes usage")
	inodesUsageOutput := format.ExecuteCommand("df", "-i")
	format.WritePreformatted(file, inodesUsageOutput)

	// MDADM
	format.WriteHeader(file, "MDADM")
	mdadmOutput := format.ExecuteCommand("cat", "/proc/mdstat")
	format.WritePreformatted(file, mdadmOutput)

	// Mounted
	format.WriteHeader(file, "Mounted")
	mountedOutput := format.ExecuteCommand("cat", "/proc/mounts")
	format.WritePreformatted(file, mountedOutput)
}
