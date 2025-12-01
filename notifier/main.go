package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 🎯 СТРУКТУРА ОБУЧЕНИЯ (12 тем, разбиты по уровням)
type Topic struct {
	Level       int      // Уровень сложности (1-7)
	Name        string   // Название темы
	Keywords    []string // Ключевые слова для поиска
	MinExamples int      // Минимум примеров, чтобы засчитать тему
	Found       int      // Сколько раз нашли в коде
}

var syllabus = []Topic{
	// LEVEL 1: Основы
	{Level: 1, Name: "Типы данных", Keywords: []string{"int", "float", "string", "bool"}, MinExamples: 3},
	{Level: 1, Name: "Переменные и константы", Keywords: []string{"var ", "const "}, MinExamples: 2},
	
	// LEVEL 2: Управление потоком
	{Level: 2, Name: "Условия (if/else)", Keywords: []string{"if ", "else"}, MinExamples: 2},
	{Level: 2, Name: "Циклы (for)", Keywords: []string{"for "}, MinExamples: 2},
	{Level: 2, Name: "Switch", Keywords: []string{"switch "}, MinExamples: 1},
	
	// LEVEL 3: Коллекции
	{Level: 3, Name: "Массивы и слайсы", Keywords: []string{"[]", "make([]", "append("}, MinExamples: 3},
	{Level: 3, Name: "Maps (карты)", Keywords: []string{"map[", "make(map"}, MinExamples: 2},
	
	// LEVEL 4: Функции
	{Level: 4, Name: "Функции", Keywords: []string{"func "}, MinExamples: 3},
	{Level: 4, Name: "Обработка ошибок", Keywords: []string{"error", "if err != nil"}, MinExamples: 2},
	
	// LEVEL 5: ООП в Go
	{Level: 5, Name: "Структуры", Keywords: []string{"type ", "struct"}, MinExamples: 2},
	{Level: 5, Name: "Методы", Keywords: []string{") func", "receiver"}, MinExamples: 2},
	{Level: 5, Name: "Интерфейсы", Keywords: []string{"interface"}, MinExamples: 1},
	
	// LEVEL 6: Concurrency
	{Level: 6, Name: "Горутины", Keywords: []string{"go func", "go "}, MinExamples: 1},
	{Level: 6, Name: "Каналы", Keywords: []string{"chan ", "<-"}, MinExamples: 2},
	
	// LEVEL 7: Продвинутое
	{Level: 7, Name: "HTTP сервер", Keywords: []string{"http.HandleFunc", "http.ListenAndServe"}, MinExamples: 1},
	{Level: 7, Name: "Тестирование", Keywords: []string{"func Test", "t.Error"}, MinExamples: 1},
}

func main() {
	fmt.Println("🔍 Сканирую Go файлы...")

	files := findGoFiles()
	if len(files) == 0 {
		fmt.Println("❌ Не найдено .go файлов")
		return
	}

	fmt.Printf("📂 Найдено файлов: %d\n", len(files))

	// Анализируем каждый файл
	for _, file := range files {
		analyzeFile(file)
	}

	// Считаем статистику
	completed := 0
	totalTopics := len(syllabus)
	currentLevel := 1
	var nextTopic string

	for i := range syllabus {
		if syllabus[i].Found >= syllabus[i].MinExamples {
			completed++
			if syllabus[i].Level > currentLevel {
				currentLevel = syllabus[i].Level
			}
		} else if nextTopic == "" {
			nextTopic = syllabus[i].Name
		}
	}

	if nextTopic == "" {
		nextTopic = "Все темы изучены! 🎉"
	}

	percent := (float64(completed) / float64(totalTopics)) * 100
	
	// Генерируем отчёт
	message := generateReport(currentLevel, percent, nextTopic, completed, totalTopics)
	
	fmt.Println("\n" + message)
	
	// Отправляем в Telegram
	sendToTelegram(message)
}

// 🔎 Поиск всех .go файлов
func findGoFiles() []string {
	var files []string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// Игнорируем папку notifier (чтобы бот не анализировал сам себя)
		if strings.Contains(path, "notifier") {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// 📊 Анализ одного файла
func analyzeFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	
	code := string(data)
	
	// Убираем комментарии и строки (чтобы не считать случайные совпадения)
	code = removeComments(code)
	
	fmt.Printf("\n📄 Анализирую: %s\n", filename)
	
	// Проверяем каждую тему
	for i := range syllabus {
		for _, keyword := range syllabus[i].Keywords {
			count := strings.Count(code, keyword)
			syllabus[i].Found += count
			if count > 0 {
				fmt.Printf("  ✓ Найдено '%s': %d раз\n", keyword, count)
			}
		}
	}
}

// 🧹 Удаляем комментарии из кода (чтобы не считать слова в комментариях)
func removeComments(code string) string {
	// Убираем // комментарии
	re1 := regexp.MustCompile(`//.*`)
	code = re1.ReplaceAllString(code, "")
	
	// Убираем /* */ комментарии
	re2 := regexp.MustCompile(`(?s)/\*.*?\*/`)
	code = re2.ReplaceAllString(code, "")
	
	// Убираем строки (чтобы "if" в строке не считался)
	re3 := regexp.MustCompile(`"[^"]*"`)
	code = re3.ReplaceAllString(code, "")
	
	return code
}

// 📝 Генерация красивого отчёта
func generateReport(level int, percent float64, nextTopic string, completed, total int) string {
	// Прогресс бар
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
	
	// Уровень опыта
	levelName := getLevelName(level)
	
	// Список тем
	var topicList strings.Builder
	topicList.WriteString("```\n")
	
	currentLvl := 0
	for _, topic := range syllabus {
		if topic.Level != currentLvl {
			currentLvl = topic.Level
			topicList.WriteString(fmt.Sprintf("\n🎯 Level %d:\n", currentLvl))
		}
		
		if topic.Found >= topic.MinExamples {
			topicList.WriteString(fmt.Sprintf("✅ %s (%d примеров)\n", topic.Name, topic.Found))
		} else {
			topicList.WriteString(fmt.Sprintf("🔒 %s (нужно %d)\n", topic.Name, topic.MinExamples))
		}
	}
	topicList.WriteString("```")
	
	return fmt.Sprintf(
		"🧙‍♂️ **GO LEARNING TRACKER**\n\n"+
			"👤 **Уровень:** %d — %s\n"+
			"📈 **Прогресс:** %s %.0f%% (%d/%d тем)\n\n"+
			"⚔️ **Следующая цель:** `%s`\n\n"+
			"📜 **Карта навыков:**\n%s\n\n"+
			"#golang #learninpublic #100daysofcode",
		level, levelName, bar, percent, completed, total, nextTopic, topicList.String(),
	)
}

// 🏆 Название уровня
func getLevelName(level int) string {
	names := map[int]string{
		1: "Новичок 🌱",
		2: "Ученик 📚",
		3: "Практикант 🔧",
		4: "Разработчик 💻",
		5: "Мастер 🎯",
		6: "Эксперт ⚡",
		7: "Гуру 🧙‍♂️",
	}
	if name, ok := names[level]; ok {
		return name
	}
	return "Новичок"
}

// 📤 Отправка в Telegram
type TGMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func sendToTelegram(text string) {
	token := os.Getenv("TELEGRAM_TOKEN")
	chatId := os.Getenv("TELEGRAM_CHAT_ID")
	
	if token == "" || chatId == "" {
		fmt.Println("⚠️ Telegram токены не найдены (это нормально для локального теста)")
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	msg := TGMessage{
		ChatID:    chatId,
		Text:      text,
		ParseMode: "Markdown",
	}
	
	jsonBody, _ := json.Marshal(msg)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	
	if err != nil {
		fmt.Printf("❌ Ошибка отправки: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Println("✅ Отчёт отправлен в Telegram!")
	} else {
		fmt.Printf("⚠️ Telegram ответил: %d\n", resp.StatusCode)
	}
}
