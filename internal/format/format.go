package format

import (
	"fmt"
	"os"
	"os/exec"
)

func WriteHeader(file *os.File, header string) {
	file.WriteString(fmt.Sprintf("<h3>%s</h3>\n", header))
}

func WriteHeaderWithID(file *os.File, header string, id string) {
	if id != "" {
		file.WriteString(fmt.Sprintf("<h3 id=\"%s\">%s</h3>\n", id, header))
	} else {
		file.WriteString(fmt.Sprintf("<h3>%s</h3>\n", header))
	}
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

// WriteHTMLHeader записывает начальную часть HTML-документа.
func WriteHTMLHeader(file *os.File) {
	file.WriteString(`<!DOCTYPE html>
		<html lang="ru">
		<head>
		    <meta http-equiv="Content-Type" content="text/html; charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>Local troubleshoot</title>
		</head>
		<body>
		`)
}

// WriteHTMLFooter закрывает HTML-документ.
func WriteHTMLFooter(file *os.File) {
	file.WriteString(`</body>
</html>
`)
}

func ListAnchorReport(file *os.File) {
	file.WriteString(`<h3>Menu</h1>
	<nav>
		<ul>
			<li><a href="#Atop">Atop info</a></li>
			<li><a href="#Process">Process run system</a></li>
			<li><a href="#Network">Network info system</a></li>
			<li><a href="#Disck">Disck info system</a></li>
			<li><a href="#Error">Error log system</a></li>
		</ul>
	</nav>
		`)
}
