package main

import (
	"fmt"
	"math/rand" // Пакет для случайных чисел
	"time"
	// Пакет для работы со временем (паузы)
)

// --- ФУНКЦИЯ: Получить текущую цену биткоина ---
// Она генерирует случайное число от 0 до 10000
func getRandomPrice() int {
	// rand.Intn(1000) дает число от 0 до 999.
	// Мы прибавляем 9000, чтобы цена "гуляла" от 9000 до 10000 (для примера)
	// Ты можешь поиграть с этими числами.
	minPrice := 5000
	maxPrice := 10000
	price := rand.Intn(maxPrice-minPrice) + minPrice
	return price
}

// --- ФУНКЦИЯ: Покупка ---
// Принимает: текущие доллары, цену.
// Возваращает: сколько осталось долларов, сколько стало биткоинов.
func buyBitcoin(usdBalance float64, price int) (float64, float64) {
	// На все деньги покупаем битки
	btcBought := usdBalance / float64(price)
	return 0.0, btcBought // Долларов стало 0, биткоинов прибавилось
}

// --- ФУНКЦИЯ: Продажа ---
// Принимает: текущие биткоины, цену.
// Возвращает: сколько стало долларов, сколько осталось биткоинов.
func sellBitcoin(btcBalance float64, price int) (float64, float64) {
	usdReceived := btcBalance * float64(price)
	return usdReceived, 0.0 // Долларов прибавилось, биткоинов стало 0
}

func main() {
	// Сетевая настройка: чтобы случайные числа были разными при каждом запуске
	// (В новых версиях Go это не обязательно, но знать полезно)
	// rand.Seed(time.Now().UnixNano())

	fmt.Println("=== CRYPTO BOT v1.0 запущен ===")

	// 1. Стартовые условия
	var walletUSD float64 = 1000.0 // Начальный капитал
	var walletBTC float64 = 0.0    // Биткоинов пока нет
	var buySignal int = 8000
	var sellSignal int = 8500

	// 2. Запускаем цикл на 10 дней
	for day := 1; day <= 10; day++ {
		// Получаем цену на сегодня
		currentPrice := getRandomPrice()
		fmt.Printf("\n[ДЕНЬ %d] Цена Bitcoin: $%d\n", day, currentPrice)

		// ЛОГИКА РОБОТА:
		// Стратегия:
		// - Если цена ниже 6000 -> ПОКУПАЙ! (если есть на что)
		// - Если цена выше 9000 -> ПРОДАВАЙ! (если есть что продать)

		if currentPrice < buySignal {
			if walletUSD > 0 {
				fmt.Println(" -> Низкая цена! ПОКУПАЕМ на всю котлету!")
				// ВНИМАНИЕ: Мы перезаписываем переменные результатами функции!
				walletUSD, walletBTC = buyBitcoin(walletUSD, currentPrice)
			} else {
				fmt.Println(" -> Хотел купить, но долларов нет (мы в позиции).")
			}

		} else if currentPrice > sellSignal {
			if walletBTC > 0 {
				fmt.Println(" -> Высокая цена! ФИКСИРУЕМ ПРИБЫЛЬ!")
				walletUSD, walletBTC = sellBitcoin(walletBTC, currentPrice)
			} else {
				fmt.Println(" -> Хотел продать, но биткоинов нет.")
			}

		} else {
			fmt.Println(" - > Цена ни то ни сё. ХОДЛИМ (ждем).")
		}

		// Выводим текущий баланс
		fmt.Printf(" Баланс: $%.2f | BTC: %.4f\n", walletUSD, walletBTC)

		// Пауза 1 секунда, чтобы было эпично (раскомментируй, если хочешь)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n=== ИТОГИ ТОРГОВЛИ ===")
	// Если остались битки, продаем их по последней цене, чтобы посчитать итог
	finalTotal := walletUSD
	if walletBTC > 0 {
		// Представим, что финальная цена средняя
		finalTotal = walletBTC * 8000.0
		fmt.Println("Продаем остатки BTC по курсу 8000...")
	}

	fmt.Printf("Начали с: $1000.00\n")
	fmt.Printf("Закончили с: $%.2f\n", finalTotal)

	profit := finalTotal - 1000
	if profit > 0 {
		fmt.Printf("ПРИБЫЛЬ: +$%.2f (Успех!)\n", profit)
	} else {
		fmt.Printf("УБЫТТОК: $%.2f (Рынок жесток...)\n", profit)
	}
}
