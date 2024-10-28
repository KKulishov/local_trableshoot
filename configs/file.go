package configs

import (
	"fmt"
	"os"
)

var Version = "unknown"

// путь к логам atop
var AtopLogDir = "/var/log/atop"

// CreateReportFile создает HTML-файл и возвращает его указатель
func CreateReportFile(filePath string) (*os.File, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании файла: %w", err)
	}
	return file, nil
}
