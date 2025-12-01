package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Читаем файл README.md из корня проекта (выходим из папки notifier на уровень вверх)
	data, err := os.ReadFile("README.md")
	if err != nil {
		fmt.Println("Ошибка: не могу найти README.md", err)
		return
	}
	content := string(data)

	// Считаем задачи
	total := strings.Count(content, "- [ ]") + strings.Count(content, "- [x]")
	done := strings.Count(content, "- [x]")

	if total == 0 {
		fmt.Println("Нет задач в списке (чекбоксов)")
		return
	}

	// Считаем процент
	percent := (float64(done) / float64(total)) * 100
	progressBar := drawProgressBar(done, total)

	// Сообщение для Телеграма
	message := fmt.Sprintf(
		"🚀 **Прогресс обучения Go**\n\n"+
			"Готово задач: %d из %d\n"+
			"Прогресс: [%s] %.1f%%\n\n"+
			"#golang #learning",
		done, total, progressBar, percent,
	)

	sendToTelegram(message)
}

func drawProgressBar(done, total int) string {
	width := 10
	filled := int((float64(done) / float64(total)) * float64(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "▓"
		} else {
			bar += "░"
		}
	}
	return bar
}

func sendToTelegram(text string) {
	token := os.Getenv("8556952823:AAGEeEqJMh5Y3LARqYLG85zqNMZ5LpJ9RIk")
	chatId := os.Getenv("-1003378333826")

	if token == "" || chatId == "" {
		fmt.Println("Ошибка: Нет токена или ID чата в секретах!")
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	jsonBody := []byte(fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "parse_mode": "Markdown"}`, chatId, text))

	http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	fmt.Println("Сообщение отправлено в Telegram!")
}
