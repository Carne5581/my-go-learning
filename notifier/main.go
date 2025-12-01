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
	"time"
)

// 🎯 СТРУКТУРА ОБУЧЕНИЯ
type Topic struct {
	Level       int
	Name        string
	Keywords    []string
	MinExamples int
	Found       int
}

// 🏆 ДОСТИЖЕНИЯ
type Achievement struct {
	ID          string
	Name        string
	Description string
	Icon        string
	Unlocked    bool
}

// 📊 СТАТИСТИКА ПОЛЬЗОВАТЕЛЯ
type UserStats struct {
	Username        string
	CurrentStreak   int
	LongestStreak   int
	TotalCommits    int
	Level           int
	CompletedTopics int
	LastCommitDate  string
	Achievements    []Achievement
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

// 🏆 Список всех достижений
var allAchievements = []Achievement{
	{ID: "first_commit", Name: "Первый шаг", Description: "Сделал первый коммит", Icon: "🎯"},
	{ID: "week_streak", Name: "Недельный марафон", Description: "7 дней подряд", Icon: "🔥"},
	{ID: "month_streak", Name: "Месячная преданность", Description: "30 дней подряд", Icon: "💪"},
	{ID: "level_3", Name: "Продвинутый новичок", Description: "Достиг 3 уровня", Icon: "⭐"},
	{ID: "level_5", Name: "ООП мастер", Description: "Достиг 5 уровня", Icon: "🎓"},
	{ID: "level_7", Name: "Go гуру", Description: "Достиг 7 уровня", Icon: "🧙‍♂️"},
	{ID: "maps_master", Name: "Властелин карт", Description: "Использовал maps 10+ раз", Icon: "🗺️"},
	{ID: "concurrency_king", Name: "Король параллелизма", Description: "Освоил горутины и каналы", Icon: "⚡"},
	{ID: "error_handler", Name: "Укротитель ошибок", Description: "Обработал 20+ ошибок", Icon: "🛡️"},
	{ID: "hundred_commits", Name: "Столетник", Description: "100 коммитов с Go кодом", Icon: "💯"},
}

func main() {
	fmt.Println("🔍 Сканирую Go файлы...")

	// Читаем статистику
	stats := loadStats()
	stats.TotalCommits++
	
	// Обновляем streak
	updateStreak(&stats)

	files := findGoFiles()
	if len(files) == 0 {
		fmt.Println("❌ Не найдено .go файлов")
		return
	}

	fmt.Printf("📂 Найдено файлов: %d\n", len(files))

	// Анализируем файлы
	for _, file := range files {
		analyzeFile(file)
	}

	// Считаем прогресс
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

	stats.Level = currentLevel
	stats.CompletedTopics = completed
	
	// Проверяем достижения
	newAchievements := checkAchievements(&stats)
	
	percent := (float64(completed) / float64(totalTopics)) * 100
	
	// Сохраняем статистику
	saveStats(stats)
	
	// Генерируем отчёт
	message := generateReport(stats, percent, nextTopic, completed, totalTopics, newAchievements)
	
	fmt.Println("\n" + message)
	
	// Отправляем в Telegram
	sendToTelegram(message)
	
	// Обновляем badges
	updateBadges(stats, percent)
	
	// Обновляем leaderboard
	updateLeaderboard(stats)
}

// 📊 Загрузка статистики
func loadStats() UserStats {
	data, err := os.ReadFile("stats.json")
	if err != nil {
		// Первый запуск
		return UserStats{
			Username:       getUsername(),
			CurrentStreak:  0,
			LongestStreak:  0,
			TotalCommits:   0,
			LastCommitDate: "",
			Achievements:   []Achievement{},
		}
	}
	
	var stats UserStats
	json.Unmarshal(data, &stats)
	return stats
}

// 💾 Сохранение статистики
func saveStats(stats UserStats) {
	data, _ := json.MarshalIndent(stats, "", "  ")
	os.WriteFile("stats.json", data, 0644)
}

// 🔥 Обновление streak
func updateStreak(stats *UserStats) {
	today := time.Now().Format("2006-01-02")
	
	if stats.LastCommitDate == "" {
		// Первый коммит
		stats.CurrentStreak = 1
		stats.LongestStreak = 1
	} else {
		lastDate, _ := time.Parse("2006-01-02", stats.LastCommitDate)
		daysDiff := int(time.Since(lastDate).Hours() / 24)
		
		if daysDiff == 1 {
			// Следующий день подряд
			stats.CurrentStreak++
			if stats.CurrentStreak > stats.LongestStreak {
				stats.LongestStreak = stats.CurrentStreak
			}
		} else if daysDiff > 1 {
			// Прервали streak
			stats.CurrentStreak = 1
		}
		// Если daysDiff == 0, то это коммит в тот же день (не меняем streak)
	}
	
	stats.LastCommitDate = today
}

// 🏆 Проверка достижений
func checkAchievements(stats *UserStats) []Achievement {
	var newAchievements []Achievement
	
	for _, achievement := range allAchievements {
		// Проверяем, не разблокировано ли уже
		alreadyUnlocked := false
		for _, unlocked := range stats.Achievements {
			if unlocked.ID == achievement.ID {
				alreadyUnlocked = true
				break
			}
		}
		
		if alreadyUnlocked {
			continue
		}
		
		// Проверяем условия
		unlocked := false
		
		switch achievement.ID {
		case "first_commit":
			unlocked = stats.TotalCommits >= 1
		case "week_streak":
			unlocked = stats.CurrentStreak >= 7
		case "month_streak":
			unlocked = stats.CurrentStreak >= 30
		case "level_3":
			unlocked = stats.Level >= 3
		case "level_5":
			unlocked = stats.Level >= 5
		case "level_7":
			unlocked = stats.Level >= 7
		case "maps_master":
			for _, topic := range syllabus {
				if topic.Name == "Maps (карты)" && topic.Found >= 10 {
					unlocked = true
				}
			}
		case "concurrency_king":
			goroutines := false
			channels := false
			for _, topic := range syllabus {
				if topic.Name == "Горутины" && topic.Found >= topic.MinExamples {
					goroutines = true
				}
				if topic.Name == "Каналы" && topic.Found >= topic.MinExamples {
					channels = true
				}
			}
			unlocked = goroutines && channels
		case "error_handler":
			for _, topic := range syllabus {
				if topic.Name == "Обработка ошибок" && topic.Found >= 20 {
					unlocked = true
				}
			}
		case "hundred_commits":
			unlocked = stats.TotalCommits >= 100
		}
		
		if unlocked {
			newAchievements = append(newAchievements, achievement)
			stats.Achievements = append(stats.Achievements, achievement)
		}
	}
	
	return newAchievements
}

// 👤 Получение username
func getUsername() string {
	// Пробуем получить из git config
	username := os.Getenv("GITHUB_ACTOR")
	if username == "" {
		username = "GoLearner"
	}
	return username
}

// 🔎 Поиск всех .go файлов
func findGoFiles() []string {
	var files []string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
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

// 📊 Анализ файла
func analyzeFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	
	code := string(data)
	code = removeComments(code)
	
	fmt.Printf("\n📄 Анализирую: %s\n", filename)
	
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

// 🧹 Удаление комментариев
func removeComments(code string) string {
	re1 := regexp.MustCompile(`//.*`)
	code = re1.ReplaceAllString(code, "")
	
	re2 := regexp.MustCompile(`(?s)/\*.*?\*/`)
	code = re2.ReplaceAllString(code, "")
	
	re3 := regexp.MustCompile(`"[^"]*"`)
	code = re3.ReplaceAllString(code, "")
	
	return code
}

// 📝 Генерация отчёта
func generateReport(stats UserStats, percent float64, nextTopic string, completed, total int, newAchievements []Achievement) string {
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
	
	levelName := getLevelName(stats.Level)
	
	// Streak сообщение
	streakMsg := ""
	if stats.CurrentStreak >= 7 {
		streakMsg = fmt.Sprintf("\n🔥 **Streak:** %d дней подряд! ", stats.CurrentStreak)
		if stats.CurrentStreak >= 30 {
			streakMsg += "Невероятно! 💪"
		} else if stats.CurrentStreak >= 14 {
			streakMsg += "Отлично! 👏"
		} else {
			streakMsg += "Продолжай! 🎯"
		}
	}
	
	// Новые достижения
	achievementMsg := ""
	if len(newAchievements) > 0 {
		achievementMsg = "\n\n🎉 **НОВЫЕ ДОСТИЖЕНИЯ:**\n"
		for _, ach := range newAchievements {
			achievementMsg += fmt.Sprintf("%s **%s** — %s\n", ach.Icon, ach.Name, ach.Description)
		}
	}
	
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
			topicList.WriteString(fmt.Sprintf("✅ %s (%d)\n", topic.Name, topic.Found))
		} else {
			topicList.WriteString(fmt.Sprintf("🔒 %s\n", topic.Name))
		}
	}
	topicList.WriteString("```")
	
	// Еженедельная статистика (если воскресенье)
	weeklyMsg := ""
	if time.Now().Weekday() == time.Sunday {
		weeklyMsg = fmt.Sprintf(
			"\n\n📅 **Недельная сводка:**\n"+
			"• Всего коммитов: %d\n"+
			"• Самый длинный streak: %d дней\n"+
			"• Достижений: %d/%d",
			stats.TotalCommits,
			stats.LongestStreak,
			len(stats.Achievements),
			len(allAchievements),
		)
	}
	
	return fmt.Sprintf(
		"🧙‍♂️ **GO LEARNING TRACKER**\n\n"+
			"👤 **%s** | Level %d — %s\n"+
			"📈 **Прогресс:** %s %.0f%% (%d/%d тем)\n"+
			"💻 **Всего коммитов:** %d%s%s\n\n"+
			"⚔️ **Следующая цель:** `%s`\n\n"+
			"📜 **Карта навыков:**\n%s%s\n\n"+
			"#golang #learninpublic #100daysofcode",
		stats.Username, stats.Level, levelName, bar, percent, completed, total,
		stats.TotalCommits, streakMsg, achievementMsg, nextTopic, topicList.String(), weeklyMsg,
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

// 🎨 Обновление badges в README
func updateBadges(stats UserStats, percent float64) {
	// Генерируем Shield.io badges
	levelBadge := fmt.Sprintf("![Level](https://img.shields.io/badge/Level-%d-blue)", stats.Level)
	progressBadge := fmt.Sprintf("![Progress](https://img.shields.io/badge/Progress-%.0f%%25-brightgreen)", percent)
	streakBadge := fmt.Sprintf("![Streak](https://img.shields.io/badge/Streak-%d_days-orange)", stats.CurrentStreak)
	commitsBadge := fmt.Sprintf("![Commits](https://img.shields.io/badge/Commits-%d-purple)", stats.TotalCommits)
	
	// Читаем README
	readmeContent, err := os.ReadFile("README.md")
	if err != nil {
		return
	}
	
	content := string(readmeContent)
	
	// Ищем секцию для замены
	badgesSection := fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		levelBadge, progressBadge, streakBadge, commitsBadge,
	)
	
	// Заменяем или добавляем badges после заголовка
	if strings.Contains(content, "![Level]") {
		// Заменяем существующие badges
		re := regexp.MustCompile(`!\[Level\].*\n!\[Progress\].*\n!\[Streak\].*\n!\[Commits\].*`)
		content = re.ReplaceAllString(content, badgesSection)
	} else {
		// Добавляем после первого заголовка
		lines := strings.Split(content, "\n")
		if len(lines) > 0 {
			lines = append(lines[:1], append([]string{"", badgesSection, ""}, lines[1:]...)...)
			content = strings.Join(lines, "\n")
		}
	}
	
	os.WriteFile("README.md", []byte(content), 0644)
	fmt.Println("✅ Badges обновлены в README.md")
}

// 📊 Обновление leaderboard
func updateLeaderboard(stats UserStats) {
	// Читаем существующий leaderboard
	type LeaderboardEntry struct {
		Username        string
		Level           int
		CompletedTopics int
		TotalCommits    int
		LongestStreak   int
	}
	
	var leaderboard []LeaderboardEntry
	data, err := os.ReadFile("LEADERBOARD.md")
	if err == nil {
		// Парсим существующий leaderboard (упрощённо)
		json.Unmarshal(data, &leaderboard)
	}
	
	// Обновляем/добавляем текущего пользователя
	found := false
	for i := range leaderboard {
		if leaderboard[i].Username == stats.Username {
			leaderboard[i].Level = stats.Level
			leaderboard[i].CompletedTopics = stats.CompletedTopics
			leaderboard[i].TotalCommits = stats.TotalCommits
			leaderboard[i].LongestStreak = stats.LongestStreak
			found = true
			break
		}
	}
	
	if !found {
		leaderboard = append(leaderboard, LeaderboardEntry{
			Username:        stats.Username,
			Level:           stats.Level,
			CompletedTopics: stats.CompletedTopics,
			TotalCommits:    stats.TotalCommits,
			LongestStreak:   stats.LongestStreak,
		})
	}
	
	// Сортируем по уровню, потом по коммитам
	// (упрощённая сортировка - в реальности нужно использовать sort.Slice)
	
	// Генерируем Markdown таблицу
	var mdContent strings.Builder
	mdContent.WriteString("# 🏆 Go Learning Leaderboard\n\n")
	mdContent.WriteString("Топ учеников Go со всего мира!\n\n")
	mdContent.WriteString("| 🏅 | Имя | Level | Темы | Коммиты | Longest Streak |\n")
	mdContent.WriteString("|---|-----|-------|------|---------|----------------|\n")
	
	for i, entry := range leaderboard {
		medal := "🥉"
		if i == 0 {
			medal = "🥇"
		} else if i == 1 {
			medal = "🥈"
		}
		
		mdContent.WriteString(fmt.Sprintf(
			"| %s | %s | %d | %d/16 | %d | %d days |\n",
			medal, entry.Username, entry.Level, entry.CompletedTopics,
			entry.TotalCommits, entry.LongestStreak,
		))
	}
	
	mdContent.WriteString("\n---\n*Обновлено автоматически*")
	
	os.WriteFile("LEADERBOARD.md", []byte(mdContent.String()), 0644)
	fmt.Println("✅ Leaderboard обновлён")
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
		fmt.Println("⚠️ Telegram токены не найдены (локальный тест)")
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
