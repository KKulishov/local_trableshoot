## Base structure 


ToDo

```bash
project-root/
│
├── cmd/
│   └── app/
│       └── main.go      # Главная точка входа программы
│
├── configs/             # Конфигурационные файлы и настройки
│   └── config.go
│
├── internal/
│   ├── model/           # Модели и структуры данных (если они нужны)
│   │   └── model.go
│   ├── mem/             # Логика работы с памятью
│   │   ├── mem.go       # Основная реализация
│   │   ├── mem_linux.go # Логика для Linux
│   │   └── mem_windows.go # Логика для Windows
│   ├── proc/            # Логика работы с процессами
│   │   ├── proc.go      # Общая логика
│   │   ├── proc_linux.go  # Реализация для Linux
│   │   └── proc_windows.go # Реализация для Windows
│   ├── disk/            # Логика работы с дисковыми устройствами
│   │   ├── disk.go
│   ├── net/             # Логика работы с сетью
│   │   ├── net.go
│   └── platform/        # Логика для определения ОС и её функционала
│       ├── platform.go
│       ├── linux/
|           └── diagnostic.go # логика реализация отчетных данных
│
└── go.mod               # Модуль Go
```

