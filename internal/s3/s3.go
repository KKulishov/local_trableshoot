package s3

import (
	"bytes"
	"context"
	"fmt"
	"local_trableshoot/internal/hostname"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"local_trableshoot/internal/flags"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Структура для хранения конфигурации
type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

var (
	proxy_S3_Host = *flags.ProxyS3Host
)

// Функция для чтения конфигурации из файла
func LoadConfig(filepath string) (*S3Config, error) {
	data, err := os.ReadFile(filepath) // Используем os.ReadFile для загрузки файла
	if err != nil {
		return nil, err
	}

	config := &S3Config{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

			switch key {
			case "endpoint_url":
				config.Endpoint = value
			case "access_key_id":
				config.AccessKeyID = value
			case "secret_access_key":
				config.SecretAccessKey = value
			case "use_ssl":
				config.UseSSL = value == "true"
			case "bucket_name":
				config.BucketName = value
			}
		}
	}

	return config, nil
}

// Функция для загрузки файла в S3 с использованием MinIO
func UploadToS3(cfg *S3Config, hostName, filePath string) error {
	// Создаем клиента MinIO
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("не удалось создать клиента MinIO: %v", err)
	}

	// Создаем контекст с таймаутом для операций S3
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Проверяем, существует ли указанный bucket, и создаем его, если его нет
	exists, err := minioClient.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return fmt.Errorf("не удалось проверить существование bucket: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			return fmt.Errorf("не удалось создать bucket: %v", err)
		}
	}

	// Читаем содержимое файла в байтовый массив
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	// Определяем путь и имя файла на S3
	objectName := filepath.Join(hostName, filepath.Base(filePath))
	contentType := "text/html; charset=utf-8"

	// Загружаем файл на S3
	info, err := minioClient.PutObject(ctx, cfg.BucketName, objectName,
		bytes.NewReader(fileData), int64(len(fileData)),
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: map[string]string{"charset": "UTF-8"},
		})
	if err != nil {
		return fmt.Errorf("не удалось загрузить файл в S3: %v", err)
	}

	if proxy_S3_Host == "" {
		fmt.Printf("Файл успешно загружен в S3. Path: %s, ETag: %s, VersionID: %s\n", objectName, info.ETag, info.VersionID)
	} else {
		fmt.Printf("Файл успешно загружен в S3. Path: https://%s/%s, ETag: %s, VersionID: %s\n", proxy_S3_Host, objectName, info.ETag, info.VersionID)
	}
	return nil
}

/*
// Функция для загрузки файла в S3 с использованием MinIO
func UploadToS3(cfg *S3Config, hostName, filePath string) error {
	// Создаем клиента MinIO
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("не удалось создать клиента MinIO: %v", err)
	}

	// Создаем контекст с таймаутом для операций S3
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверяем, существует ли указанный bucket, и создаем его, если его нет
	exists, err := minioClient.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return fmt.Errorf("не удалось проверить существование bucket: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			return fmt.Errorf("не удалось создать bucket: %v", err)
		}
	}

	// Читаем содержимое файла в байтовый массив
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	// Создаем путь на S3 с учетом имени хоста
	objectName := filepath.Join(hostName, filepath.Base(filePath))
	contentType := "text/html; charset=utf-8"
	//contentType := "text/plain"

	// Загружаем файл на S3
	info, err := minioClient.PutObject(ctx, cfg.BucketName, objectName,
		bytes.NewReader(fileData), int64(len(fileData)),
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: map[string]string{"charset": "UTF-8"},
		})
	if err != nil {
		return fmt.Errorf("не удалось загрузить файл в S3: %v", err)
	}

		info, err := minioClient.FPutObject(ctx, cfg.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType,
			ContentEncoding: "UTF-8",
			UserMetadata:    map[string]string{"charset": "UTF-8"},
		})
		if err != nil {
			return fmt.Errorf("не удалось загрузить файл в S3: %v", err)
		}


	fmt.Printf("Файл успешно загружен в S3. Path: %s, ETag: %s, VersionID: %s\n", objectName, info.ETag, info.VersionID)
	return nil
}
*/

// DeleteOldFiles оставляет последние `retainCount` файлов в папке `hostPath`
func deleteOldFiles(cfg *S3Config, hostPath string, retainCount int) error {
	// Список для хранения информации о файлах
	var objectInfos []minio.ObjectInfo

	// Создаем клиента MinIO
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("не удалось создать клиента MinIO: %v", err)
	}

	// Инициализируем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Получаем список объектов в указанной директории
	for object := range minioClient.ListObjects(ctx, cfg.BucketName, minio.ListObjectsOptions{
		Prefix:    hostPath,
		Recursive: true,
	}) {
		if object.Err != nil {
			return object.Err
		}
		objectInfos = append(objectInfos, object)
	}

	// Проверяем, нужно ли удалять файлы
	if len(objectInfos) <= retainCount {
		fmt.Println("Нет необходимости удалять файлы, количество файлов меньше или равно указанному лимиту.")
		return nil
	}

	// Сортируем файлы по дате последнего изменения (от старых к новым)
	sort.Slice(objectInfos, func(i, j int) bool {
		return objectInfos[i].LastModified.Before(objectInfos[j].LastModified)
	})

	// Удаляем старые файлы, оставляем только `retainCount` последних файлов
	for i := 0; i < len(objectInfos)-retainCount; i++ {
		object := objectInfos[i]
		err := minioClient.RemoveObject(ctx, cfg.BucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("error deleting file %s: %w", object.Key, err)
		}
		fmt.Printf("Удален файл %s\n", object.Key)
	}

	return nil
}

func Rotation_s3_bucket(filePath string, retain int) {
	configPath := os.Getenv("HOME") + "/.config/report_send_s3"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Файл конфигурации s3 не найден: %v\n", err)
	} else {
		cfg, err := LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Ошибка при загрузке конфигурации: %v\n", err)
			return
		}
		// Получаем имя хоста
		hostName := hostname.HostName()
		deleteOldFiles(cfg, hostName, retain)
	}
}

// отправка отчета в s3
func Send_report_file(filePath string) {
	// Загружаем конфигурацию из файла для s3
	configPath := os.Getenv("HOME") + "/.config/report_send_s3"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Файл конфигурации s3 не найден: %v\n", err)
	} else {
		cfg, err := LoadConfig(configPath)
		if err != nil {
			fmt.Printf("Ошибка при загрузке конфигурации: %v\n", err)
			return
		}
		// Получаем имя хоста
		hostName := hostname.HostName()
		// Загружаем файл в S3 с учетом имени хоста
		err = UploadToS3(cfg, hostName, filePath)
		if err != nil {
			fmt.Printf("Ошибка при загрузке в S3: %v\n", err)
		}
	}
}
