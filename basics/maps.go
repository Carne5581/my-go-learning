package main

import "fmt"

func main() {
	fmt.Println("=== CRYPTO PORTFOLIO (MAPS) ===")

	// 1. СОЗДАНИЕ MAP
	// map[ключ]значение
	// Ключи - строки (названия), Значения - числа (баланс)
	portfolio := map[string]float64{
		"BTC":  0.5,
		"ETH":  10.0,
		"USDT": 1000,
	}

	// Добавляем новую валюту просто приваиванием
	portfolio["TON"] = 500.0

	fmt.Println("Мой портфель:", portfolio)
	// ВАЖНО: Map в Go не гарантирует порядок! При запуске порядок может меняться.

	// 2. ПОИСК ЗНАЧЕНИЯ (Инеженерный подход)
	fmt.Println("\n--- ПРОВЕРКА БАЛАНСА ----")

	// Попробуем достать то, что есть
	btcBalance := portfolio["BTC"]
	fmt.Printf("У меня есть %.2f BTC\n", btcBalance)

	// Попробуем достать то, чего НЕТ
	// Go не выдаст ошибку! Он вернет "0" (пустое значение).
	dogeBalance := portfolio["DOGE"]
	fmt.Printf("У меня есть %.2f DOGE (Хотя я их не покупал!)\n", dogeBalance)

	// КАК ПРОВЕРИТЬ, ЧТО ВАЛЮТА РЕАЛЬНО СУЩЕСТВУЕТ?
	// Используем синтаксис: val, ok := map[key]

	val, exists := portfolio["DOGE"]
	if exists {
		fmt.Println("DOGE найден, баланс:", val)
	} else {
		fmt.Println("DOGE в портфеле НЕ НАЙДЕН. (Это честная проверка)")
	}

	// 3. УДАЛЕНИЕ
	fmt.Println("\n--- ПРОДАЛ ETH ---")
	delete(portfolio, "ETH") // Встроенная функция delete
	fmt.Println("Портфель после продажи:", portfolio)

	// 4. ЦИКЛ ПО КАРТЕ
	fmt.Println("\n--- ПОЛНЫЙ ОТЧЕТ ---")
	for coin, balance := range portfolio {
		fmt.Printf("Валюта: %s | Баланс: %.2f\n", coin, balance)
	}
}
