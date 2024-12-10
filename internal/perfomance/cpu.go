package perfomance

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// HeavyMathOperations выполняет сложные математические расчеты
func heavyMathOperations() (float64, string) {
	start := time.Now()

	// Количество итераций для теста
	iterations := 1000000
	result := 0.0

	for i := 1; i <= iterations; i++ {
		// Выполняем математические операции
		// Выполняет расчет, который включает в себя квадратный корень, логарифм и синус для каждого числа от 1 до миллиона iterations.
		// Вычисляется суммарное значение result для всех операций.
		result += math.Sqrt(float64(i)) * math.Log(float64(i)) * math.Sin(float64(i))
	}

	duration := time.Since(start)

	durationStr := duration.String()

	//fmt.Printf("Результат вычислений: %.5f\n", result)
	//fmt.Printf("Время выполнения: %s\n", duration)
	return result, durationStr
}

// SetHighPriority изменяет приоритет процесса
func setHighPriority() error {
	pid := os.Getpid()
	cmd := exec.Command("renice", "-n", "-20", "-p", fmt.Sprintf("%d", pid))
	return cmd.Run()
}

func RunCpuResults(file *os.File) {
	// Заголовок секции
	fmt.Fprintln(file, "<h3>Perfomance CPU</h3>")
	fmt.Fprintln(file, "<div><pre>")

	// Устанавливаем максимальный приоритет (если есть права)
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		err := setHighPriority()
		if err != nil {
			fmt.Println("Не удалось установить высокий приоритет:", err)
		} else {
			fmt.Println("Высокий приоритет установлен.")
		}
	}

	// Выполняем вычисления
	opsCpu, timeOps := heavyMathOperations()

	fmt.Fprintln(file, `Выполняет расчет, который включает в себя квадратный корень, логарифм и синус для каждого числа от 1 до миллиона`)
	fmt.Fprintf(file, "Результат вычислений: %.5f\n", opsCpu)
	fmt.Fprintf(file, "Время выполнения: %s\n", timeOps)

	fmt.Fprintln(file, "</pre></div>")

}

func RunCpuResultsSdout() {
	// Устанавливаем максимальный приоритет (если есть права)
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		err := setHighPriority()
		if err != nil {
			fmt.Println("Не удалось установить высокий приоритет:", err)
		} else {
			fmt.Println("Высокий приоритет установлен.")
		}
	}

	// Выполняем вычисления
	opsCpu, timeOps := heavyMathOperations()

	fmt.Println(`Выполняет расчет, который включает в себя квадратный корень, логарифм и синус для каждого числа от 1 до миллиона`)
	fmt.Printf("Результат вычислений: %.5f\n", opsCpu)
	fmt.Printf("Время выполнения: %s\n", timeOps)

}
