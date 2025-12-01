package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ---------------------------------------------------------
// 📜 ТВОЙ УЧЕБНЫЙ ПЛАН (SYLLABUS)
// Бот ищет "Keyword" в твоих файлах. Если находит — ставит галочку.
// ---------------------------------------------------------

type Topic struct {
	Name    string // Красивое название темы
	Keyword string // Код, который бот ищет в файлах
	IsDone  bool   // (не трогать) Статус выполнения
}

var syllabus = []Topic{
	{Name: "Переменные", Keyword: "var "},
	{Name: "Функции", Keyword: "func "},
	{Name: "Циклы", Keyword: "for "},
	{Name: "Условия", Keyword: "if "},
	{Name: "Массивы/Слайсы", Keyword: "[]"},
	{Name: "Карты (Maps)", Keyword: "map["},
	{Name: "Структуры", Keyword: "struct"},
	{Name: "Методы", Keyword: ") Set"}, // Поиск методов структур
	{Name: "Интерфейсы", Keyword: "interface"},
	{Name: "Горутины", Keyword: "go func"},
	{Name: "Каналы", Keyword: "chan "},
	{Name: "Обработка ошибок", Keyword: "if err !="},
}

// ---------------------------------------------------------

func main() {
	// 1. Читаем весь код из файлов
	fullCode := readAllGoFiles()

	if len(fullCode) == 0 {
		fmt.Println("Код не найден! Напиши хоть строчку.")
		return
	}

	// 2. Проверяем план
	completedCount := 0
	var nextTarget string = "Все изучено! 🎉"
	foundNext := false

	// Проходим по списку и ставим галочки
	for i := range syllabus {
		if strings.Contains(fullCode, syllabus[i].Keyword) {
			syllabus[i].IsDone = true
			completedCount++
		} else {
			// Запоминаем первую невыполненную задачу как цель
			if !foundNext {
				nextTarget = syllabus[i].Name
				foundNext = true
			}
		}
	}

	// 3. Считаем статистику
	total := len(syllabus)
	percent := (float64(completedCount) / float64(total)) * 100
	level := completedCount + 1 // Уровень героя = кол-во тем + 1

	// 4. Генерируем красивое сообщение
	message := generateFancyReport(level, percent, nextTarget, syllabus)

	// 5. Отправляем
	sendToTelegram(message)
}

func generateFancyReport(level int, percent float64, next string, topics []Topic) string {
	// Рисуем бар
	barWidth := 10
	filled := int((percent / 100) * float64(barWidth))
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "🟩"
		} else {
			bar += "⬜"
		}
	}

	// Собираем список достижений (показываем только последние 3 или важное)
	// Но для красоты выведем список: Сделано / Не сделано
	listBuilder := ""
	for _, t := range topics {
		if t.IsDone {
			listBuilder += "✅ " + t.Name + "\n"
		} else {
			listBuilder += "🔒 " + t.Name + "\n"
		}
	}

	return fmt.Sprintf(
		"🧙‍♂️ **GOLANG HERO REPORT**\n"+
			"👤 **Уровень:** %d (Novice)\n"+
			"📈 **Прогресс:** %s %.0f%%\n\n"+
			"⚔️ **Текущая цель:** `%s`\n\n"+
			"📜 **Карта навыков:**\n%s\n"+
			"#golang #levelup #buildinpublic",
		level, bar, percent, next, listBuilder,
	)
}

func readAllGoFiles() string {
	var allCode string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			data, _ := os.ReadFile(path)
			allCode += string(data) + "\n"
		}
		return nil
	})
	return allCode
}

// Структура для отправки JSON в Telegram (чтобы смайлики не ломались)
type TGMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func sendToTelegram(text string) {
	token := os.Getenv("TELEGRAM_TOKEN")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")
	if token == "" || chatId == "" {
		fmt.Println("Нет токена!")
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	
	msg := TGMessage{
		ChatID:    chatId,
		Text:      text,
		ParseMode: "Markdown",
	}

	jsonBody, _ := json.Marshal(msg)
	http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}
