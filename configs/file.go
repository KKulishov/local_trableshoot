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

func CheckAndCreateDir(path string) error {
	// Проверка существует ли каталог
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		// Если каталог не существует, создаем его
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("не удалось создать каталог: %v", err)
		}
		fmt.Println("Каталог", path, "успешно создан.")
	} else if err != nil {
		// Если произошла ошибка, отличная от "каталог не существует"
		return fmt.Errorf("ошибка при проверке каталога: %v", err)
	} else {
		// Если каталог существует
		//fmt.Println("Каталог", path, "уже существует.")
	}

	return nil
}
