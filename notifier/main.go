package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Словарь: "Тема": "Код, который должен быть в файлах"
var topics = map[string]string{
	"Переменные (var)":    "var ",
	"Циклы (for)":         "for ",
	"Функции (func)":      "func ",
	"Массивы/Слайсы ([])": "[]",
	"Карты (map)":         "map[",
	"Структуры (struct)":  "struct",
	"Указатели (*)":       "*",
	"Горутины (go)":       "go func",
	"Интерфейсы":          "interface",
}

func main() {
	// 1. Читаем ВЕСЬ код из всех .go файлов в одну кучу
	fullCode := readAllGoFiles()

	if len(fullCode) == 0 {
		fmt.Println("Код не найден! Напиши хоть что-нибудь.")
		return
	}

	// 2. Проверяем, какие темы встречаются в коде
	completed := 0
	total := len(topics)
	var doneList []string

	for name, keyword := range topics {
		if strings.Contains(fullCode, keyword) {
			completed++
			doneList = append(doneList, "✅ "+name)
		} else {
			doneList = append(doneList, "⬜ "+name)
		}
	}

	// 3. Считаем процент
	percent := (float64(completed) / float64(total)) * 100
	progressBar := drawProgressBar(completed, total)

	// 4. Формируем отчет
	// strings.Join собирает список тем в красивый столбик
	msgText := fmt.Sprintf(
		"🧠 **Анализ кода завершен!**\n\n"+
			"Я просканировал твои файлы.\n"+
			"Изучено тем: %d из %d\n"+
			"Прогресс: [%s] %.1f%%\n\n"+
			"**Детали:**\n%s\n\n"+
			"#golang #tracker",
		completed, total, progressBar, percent, strings.Join(doneList, "\n"),
	)

	sendToTelegram(msgText)
}

// Функция ходит по папкам и собирает весь текст из .go файлов
func readAllGoFiles() string {
	var allCode string
	// Walk ищет файлы во всех папках
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" { // Если файл заканчивается на .go
			data, _ := os.ReadFile(path)
			allCode += string(data) + "\n"
		}
		return nil
	})
	return allCode
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
	token := os.Getenv("TELEGRAM_TOKEN")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")
	if token == "" || chatId == "" {
		return
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	jsonBody := []byte(fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "parse_mode": "Markdown"}`, chatId, text))
	http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}
