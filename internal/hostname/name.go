package hostname

import (
	"fmt"
	"local_trableshoot/configs"
	"local_trableshoot/internal/format"
	"os"
	"time"
)

func GetVersionApp(file *os.File) {
	format.WriteHeader(file, "Version app trableshoot")
	format.WritePreformatted(file, configs.Version)
}

func GetHostName(file *os.File) {
	// Hostname
	format.WriteHeader(file, "Hostname")
	hostnameOutput := format.ExecuteCommand("hostname")
	format.WritePreformatted(file, hostnameOutput)

	// Date and Time
	format.WriteHeader(file, "Дата и время отчёта")
	currentTime := time.Now().Format("02.01.2006 time: 15:04:05")
	format.WritePreformatted(file, currentTime)
}

func HostName() string {
	// Получение имени хоста
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Ошибка получения имени хоста:", err)
	}
	return hostname
}
