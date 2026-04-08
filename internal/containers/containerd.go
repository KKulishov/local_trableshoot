package containers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type criStatsResponse struct {
	Stats []criContainerStat `json:"stats"`
}

type criContainerStat struct {
	Attributes struct {
		ID       string            `json:"id"`
		Labels   map[string]string `json:"labels"`
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
	} `json:"attributes"`
	CPU struct {
		UsageNanoCores struct {
			Value string `json:"value"`
		} `json:"usageNanoCores"`
		UsageCoreNanoSeconds struct {
			Value string `json:"value"`
		} `json:"usageCoreNanoSeconds"`
	} `json:"cpu"`
	Memory struct {
		WorkingSetBytes struct {
			Value string `json:"value"`
		} `json:"workingSetBytes"`
	} `json:"memory"`
}

func GetContainerdStatCpu(file *os.File) {
	file.WriteString("<html><head><title>Containerd Stats Report</title></head><body>\n")
	file.WriteString("<h2>Top 10 Containers by CPU Usage</h2>\n")
	file.WriteString("<table border='1'>\n")
	file.WriteString("<tr><th>Namespace</th><th>Pod name</th><th>Container Name</th><th>CPU %</th></tr>\n")

	stats, err := getCRIStats()
	if err != nil {
		fmt.Fprintln(file, "Ошибка при выполнении команды:", err)
		return
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].cpuPercent > stats[j].cpuPercent
	})

	if len(stats) > 10 {
		stats = stats[:10]
	}

	for _, stat := range stats {
		row := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%.2f%%</td></tr>\n", stat.namespace, stat.podName, stat.name, stat.cpuPercent)
		file.WriteString(row)
	}

	file.WriteString("</table>\n")
	file.WriteString("</body></html>\n")
}

func GetContainerdStatMem(file *os.File) {
	file.WriteString("<html><head><title>Containerd Stats Report</title></head><body>\n")
	file.WriteString("<h2>Top 10 Containers by MEM Usage</h2>\n")
	file.WriteString("<table border='1'>\n")
	file.WriteString("<tr><th>Namespace</th><th>Pod name</th><th>Container Name</th><th>MEM Usage</th></tr>\n")

	stats, err := getCRIStats()
	if err != nil {
		fmt.Fprintln(file, "Ошибка при выполнении команды:", err)
		return
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].memUsageBytes > stats[j].memUsageBytes
	})

	if len(stats) > 10 {
		stats = stats[:10]
	}

	for _, stat := range stats {
		row := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n", stat.namespace, stat.podName, stat.name, stat.memUsageHuman)
		file.WriteString(row)
	}

	file.WriteString("</table>\n")
	file.WriteString("</body></html>\n")
}

type containerdStatRow struct {
	namespace     string
	podName       string
	name          string
	cpuPercent    float64
	memUsageBytes int64
	memUsageHuman string
}

func getCRIStats() ([]containerdStatRow, error) {
	cmd := exec.Command("crictl", "stats", "--output", "json")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if strings.TrimSpace(stderr.String()) != "" {
			return nil, fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return nil, err
	}

	var payload criStatsResponse
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		return nil, err
	}

	result := make([]containerdStatRow, 0, len(payload.Stats))
	for _, item := range payload.Stats {
		cpuNanoCores, _ := strconv.ParseInt(item.CPU.UsageNanoCores.Value, 10, 64)
		memVal, _ := strconv.ParseInt(item.Memory.WorkingSetBytes.Value, 10, 64)

		name := strings.TrimSpace(item.Attributes.Metadata.Name)
		if name == "" {
			name = item.Attributes.ID
		}

		result = append(result, containerdStatRow{
			namespace: firstNonEmpty(
				item.Attributes.Labels["io.kubernetes.pod.namespace"],
				item.Attributes.Labels["pod.namespace"],
				item.Attributes.Labels["namespace"],
				"-",
			),
			podName: firstNonEmpty(
				item.Attributes.Labels["io.kubernetes.pod.name"],
				item.Attributes.Labels["pod.name"],
				"-",
			),
			name:          name,
			cpuPercent:    nanoCoresToPercent(cpuNanoCores),
			memUsageBytes: memVal,
			memUsageHuman: humanizeBytes(memVal),
		})
	}

	return result, nil
}

func humanizeBytes(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	units := []string{"KiB", "MiB", "GiB", "TiB"}
	val := float64(size)
	unitIdx := -1
	for val >= 1024 && unitIdx < len(units)-1 {
		val /= 1024
		unitIdx++
	}

	if unitIdx < 0 {
		return fmt.Sprintf("%d B", size)
	}
	return fmt.Sprintf("%.2f %s", val, units[unitIdx])
}

func nanoCoresToPercent(nanoCores int64) float64 {
	// 1 core = 1e9 nanoCores, so percent of one core is nanoCores/1e9*100.
	return (float64(nanoCores) / 1_000_000_000.0) * 100.0
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}
