package s3

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

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
func UploadToS3(cfg *S3Config, filePath string) error {
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

	// Загружаем файл в bucket
	objectName := filePath // Имя файла в bucket
	contentType := "text/plain"

	// Загружаем файл на S3
	info, err := minioClient.FPutObject(ctx, cfg.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("не удалось загрузить файл в S3: %v", err)
	}

	fmt.Printf("Файл успешно загружен в S3. ETag: %s, VersionID: %s\n", info.ETag, info.VersionID)
	return nil
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
		err = UploadToS3(cfg, filePath)
		if err != nil {
			fmt.Printf("Ошибка при загрузке в S3: %v\n", err)
		}
	}
}
