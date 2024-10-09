package format

import (
	"fmt"
	"os"
	"os/exec"
)

func WriteHeader(file *os.File, header string) {
	file.WriteString(fmt.Sprintf("<h3>%s</h3>\n", header))
}

func WritePreformatted(file *os.File, content string) {
	file.WriteString(fmt.Sprintf("<div><pre>%s</pre></div>\n", content))
}

func ExecuteCommand(command string, args ...string) string {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error executing command: %s", err)
	}
	return string(output)
}
