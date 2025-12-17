package main

import (
	"fmt"
)

// СТАДИЯ 1: Генератор чисел
// chan<- int означает: "В этот канал можно ТОЛЬКО ПИСАТЬ"
func generateNumbers(out chan<- int) {
	for i := 1; i <= 5; i++ {
		fmt.Println("📤 Генератор: Отправляю", i)
		out <- i // Кидаем число в трубу
	}
	// ВАЖНО: Закрываем канал, когда закончили.
	// Это сгинал для следующего этапа: "Данных больше не будет, расходимся".
	close(out)
}

// СТАДИЯ 2: Обработчик (Квадрат)
// <-chan int (ТОЛЬКО ЧИТАТЬ)
// chan<- int (ТОЛЬКО ПИСАТЬ)
func squareNumbers(in <-chan int, out chan<- int) {
	// range по каналу работает, пока канал не закроют (close)
	for num := range in {
		result := num * num
		fmt.Printf(" ⚙️ Обработчик: %d -> %d\n", num, result)
		out <- result // Кидаем дальше
	}
	close(out) // Тоже закрываем выходной канал
}

// СТАДИЯ 3: Умножитель
func multiplyByTwo(in <-chan int, out chan<- int) {
	for num := range in {
		result := num * 2
		fmt.Printf(" Умножитель: %d -> %d\n", num, result)
		out <- result
	}
	close(out)
}

func main() {
	fmt.Println("=== PIPELINE SYSTEM v1.0 ===")

	// 1. Создаем трубы (Каналы)
	// make(chan Тип)
	numbersChan := make(chan int)
	squaredChan := make(chan int)
	finalChan := make(chan int)

	// 2. Запускаем завод (Горутины)
	// Связываем их трубами
	go generateNumbers(numbersChan)
	go squareNumbers(numbersChan, squaredChan)
	go multiplyByTwo(squaredChan, finalChan)

	// 3. Главная функция работает как "Приемщик"
	// Мы читаем из последней трубы, пока она не закроется
	fmt.Println("🏁 Main: Жду результаты...")

	for res := range finalChan {
		fmt.Println(" ✅ Итог:", res)
	}

	// Нам не нужен WaitGroup!
	// range squaredChan сам подождет все данные и выйдет, когда канал закроют.
}
