package main

import "fmt"

// --- ЭТО НАШИ ИНСТРУМЕНТЫ (ФУНКЦИИ) ---

// 1. Простая функция. Принимает рубли, возвращает доллары.
// func Имя(параметр Тип) ТипВозврата
func convertToUSD(rubles int) int {
	// Допустим, курс 100 (для простоты)
	usd := rubles / 100
	return usd // Возвращаем результат тому, кто вызвал
}

func convertToEUR(rubles2 int) int {
	eur := rubles2 / 110
	return eur
}

// 2. Функция с логикой (имитация запроса в банк)
// Возвращает два значения: курс и ошибку (в Go так принято!)
func getExchangeRate(currency string) (int, string) {
	if currency == "USD" {
		return 100, "OK" // Курс 100, ошибок нет
	} else if currency == "EUR" {
		return 110, "OK"
	} else {
		return 0, "Error: Unknown currency"
	}
}

// 3. ИНЖЕНЕРНЫЙ ЭКСПЕРИМЕНТ (Pass by Value)
// Мы пытаемся изменить баланс клиента внутри функции.
// Получится ли?
func tryToCheatBonus(balance int) {
	fmt.Println(" [Внутри функции] Пытаюсь накрутить бонус...")
	balance = balance + 1000000
	fmt.Println(" [Внутри функции] Ура, баланс теперь:", balance)
}

// --- ГЛАВНАЯ ПРОГРАММА (ТОЧКА ВХОДА) ---
func main() {
	fmt.Println("=== BANK SYSTEM v1.0 ===")

	// 1. Тест конвертации
	myMoney := 5000
	dollars := convertToUSD(myMoney)
	fmt.Printf("У меня было %d руб, стало %d баксов.\n", myMoney, dollars)

	// 2. Тест получения курса (возврат нескольких значение)
	rate, status := getExchangeRate("EUR")
	fmt.Printf("Курс Евро: %d, Статус системы: %s\n", rate, status)

	// А что если валюта неверная?
	rate2, status2 := getExchangeRate("TUGRIK")
	fmt.Printf("Курс Тугрика: %d. Статус: %s\n", rate2, status2)

	// 3. САМОЕ ВАЖНОЕ: Копирование значений
	fmt.Println("\n--- Тест безопасности памяти ---")
	clientBalance := 100
	fmt.Println("Мой баланс ДО функции:", clientBalance)

	// Вызываем функцию-хакера
	tryToCheatBonus(clientBalance)

	// Внимание на экран! Изменился ли баланс в main?
	fmt.Println("Мой баланс ПОСЛЕ функции:", clientBalance)

	weMoney := 4999
	euros := convertToEUR(weMoney)
	fmt.Printf("У меня было %d руб, стало %d евро.\n", weMoney, euros)
}
