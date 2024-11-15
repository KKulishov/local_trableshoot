package disk

import (
	"bufio"
	"fmt"
	"local_trableshoot/internal/format"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

type ProcessInfo struct {
	TID     int
	PRIO    string
	User    string
	Read    string
	Write   string
	SwapIn  string
	IO      string
	Command string
}

func AppDiskUtilization(file *os.File) error {
	cmd := exec.Command("sudo", "iotop", "-o", "-b", "-n", "5")
	var output strings.Builder
	cmd.Stdout = &output

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ошибка при запуске команды: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ошибка при выполнении команды: %w", err)
	}

	processes := parseIOTopOutput(output.String())

	if err := saveToHTML(file, processes); err != nil {
		return fmt.Errorf("ошибка при сохранении данных: %w", err)
	}
	return nil
}

func parseIOTopOutput(output string) []ProcessInfo {
	var processes []ProcessInfo
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "    TID ") ||
			strings.HasPrefix(line, "Total DISK READ:") ||
			strings.HasPrefix(line, "Actual DISK READ:") {
			// Пропускаем строки заголовков или итогов
			continue
		}

		// Удаление `b'` и конечного `'` символа, если они присутствуют
		line = strings.TrimPrefix(line, "b'")
		line = strings.TrimSuffix(line, "'")

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		tid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}

		process := ProcessInfo{
			TID:     tid,
			PRIO:    fields[1],
			User:    fields[2],
			Read:    fields[3],
			Write:   fields[4],
			SwapIn:  fields[5],
			IO:      fields[6],
			Command: strings.Join(fields[7:], " "),
		}
		processes = append(processes, process)
	}

	return processes
}

func saveToHTML(file *os.File, processes []ProcessInfo) error {
	_, err := file.WriteString("<html><body><h3 id=\"Disck\">Disk Utilization Report</h3><table border='1'>")
	if err != nil {
		return err
	}

	_, err = file.WriteString("<tr><th>TID</th><th>PRIO</th><th>User</th><th>Read</th><th>Write</th><th>SwapIn</th><th>IO</th><th>Command</th></tr>")
	if err != nil {
		return err
	}

	for _, process := range processes {
		row := fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			process.TID, process.PRIO, process.User, process.Read, process.Write, process.SwapIn, process.IO, process.Command)
		_, err := file.WriteString(row)
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString("</table></body></html>")
	if err != nil {
		return err
	}

	return nil
}
