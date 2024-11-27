package rotate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CleanUpOldReports удаляет все, кроме последних 'retain' файлов, которые соответствуют шаблону.
func CleanUpOldReports(dir string, prefix string, retain int) {
	// Получаем список файлов в директории
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Ошибка при чтении директории:", err)
		return
	}

	// Фильтруем файлы, оставляя только те, что начинаются с prefix и имеют указанные расширения
	var reportFiles []fs.DirEntry
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) &&
			(strings.HasSuffix(file.Name(), ".html") ||
				strings.HasSuffix(file.Name(), ".data") ||
				strings.HasSuffix(file.Name(), ".pcap")) {
			reportFiles = append(reportFiles, file)
		}
	}

	// Если файлов меньше или равно количеству retain, ничего не делаем
	if len(reportFiles) <= retain {
		return
	}

	// Сортируем файлы по времени модификации (сначала старые)
	sort.Slice(reportFiles, func(i, j int) bool {
		infoI, errI := reportFiles[i].Info()
		infoJ, errJ := reportFiles[j].Info()
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	// Удаляем все файлы, кроме последних retain
	for _, file := range reportFiles[:len(reportFiles)-retain] {
		filePath := filepath.Join(dir, file.Name())
		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("Не удалось удалить файл %s: %v\n", file.Name(), err)
		} else {
			fmt.Printf("Удален старый отчет: %s\n", file.Name())
		}
	}
}
