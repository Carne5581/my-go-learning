package main

import (
	"fmt"
	"time"
)

// Имитация ответа от сервера
func serverResponse(serverName string, delay time.Duration, out chan<- string) {
	// Имитируем работу (сон)
	time.Sleep(delay)
	// Отправляем ответ в канал
	out <- serverName + ": Данные получены!"
}

func main() {
	fmt.Println("=== SERVER RACE SYSTEM v1.0 ===")

	// 1. Создаем каналы для ответов
	// Делаем их буферизированными (размер 1), чтобы горутина не зависла,
	// если мы уйдем по тайм-ауту и никто не прочитает её ответ.
	chan1 := make(chan string, 1)
	chan2 := make(chan string, 1)
	messages := make(chan string, 5)

	// 2. Запускаем горутины (Запросы к серверам)
	// Сервер 1: Быстрый (1 сек)
	go serverResponse("Server A (Fast)", 2*time.Second, chan1)

	// Сервер 2: Медленный (3 сек)
	go serverResponse("Server B (Slow)", 3*time.Second, chan2)

	// 3. SELECT - Ждем, кто первый
	fmt.Println("⏳ Жду ответа...")

	messages <- "Привет"
	messages <- "Как дела?"
	messages <- "Я учу Go"

	for i := 0; i < 4; i++ {
		select {
		case msg := <-messages:
			fmt.Println("Получено сообщение: ", msg)
		default:
			fmt.Println(" Сообщений нет, я пока сплю...")
		}
		time.Sleep(500 * time.Millisecond)
	}

	close(messages)

	select {
	case msg1 := <-chan1:
		fmt.Println("✅ Победа!", msg1)

	case msg2 := <-chan2:
		fmt.Println("✅ Победа!", msg2)

	// time.After создаем канал, в который прилетает сигнал через указанное время
	case <-time.After(2 * time.Second):
		fmt.Println("⛔️ ОШИБКА: Время вышло! (Timeout)")
	}

	fmt.Println("🏁 Программа завершена.")
}
