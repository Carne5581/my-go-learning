package main

import "fmt"

func main() {
	fmt.Println("=== WAREHOUSE SYSTEM v1.0 ===")

	// --- ЗАДАНИЕ: СЧИТАЕМ СУММУ ЦЕН ---
	fmt.Println("\n--- ПОДСЧЕТ СУММЫ (ТВОЕ ЗАДАНИЕ) ---")

	// 1. Создаем СЛАЙС чисел (int), а не строк!
	// []int - значит "слайс целых чисел"
	// {} - позволяют сразу наполнить его данными
	prices := []int{100, 200, 300}

	fmt.Println("Цены товаров:", prices)

	// 2. Создаем переменную для общей суммы (Копилка)
	// Важно создать её ДО цикла, чтобы она накапливала результат.
	totalSum := 0

	// 3. Цикл
	// range возвращает два значения: (индекс, значение)
	// Индекс нам не нужен, поэтому ставим прочерк (_)
	// price - это конкретное число (сначала 100, потом 200, потом 300)
	for _, price := range prices {
		// Прибавляем текущую цену к общей сумме
		totalSum = totalSum + price
	}

	fmt.Println("Итоговая сумма:", totalSum)

	// --- ЧАСТЬ 2: СЛАЙС (Для повторения) ---
	fmt.Println("\n=== РАБОТА СО СЛАЙСАМИ (ТЕОРИЯ) ===")

	var boxes []string

	// Исправленная строка Printf (у тебя кавычка стояла не там)
	fmt.Printf("Старт: len=%d, cap=%d, Адрес=%p\n", len(boxes), cap(boxes), boxes)

	boxes = append(boxes, "Коробка 1")
	boxes = append(boxes, "Коробка 2")

	fmt.Printf("Добавили 2 коробки: %v\n", boxes)
	fmt.Printf("Состояние: len=%d, cap=%d, Адрес=%p\n", len(boxes), cap(boxes), boxes)

	// Переполнение
	boxes = append(boxes, "Коробка 3")
	fmt.Printf("НОВОЕ Состояние: len=%d, cap=%d, Адрес=%p\n", len(boxes), cap(boxes), boxes)
}
