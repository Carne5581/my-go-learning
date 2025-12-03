package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	XPReward    int // XP за изучение темы
	Found       int
}

// 🏆 ДОСТИЖЕНИЯ
type Achievement struct {
	ID          string
	Name        string
	Description string
	Icon        string
	XPReward    int
	Unlocked    bool
}

// 📊 СТАТИСТИКА ПОЛЬЗОВАТЕЛЯ
type UserStats struct {
	Username        string
	TotalXP         int
	CurrentStreak   int
	LongestStreak   int
	TotalCommits    int
	Level           int
	League          string
	CompletedTopics int
	LastCommitDate  string
	Achievements    []Achievement
	PenaltyDays     int // Дни без коммитов
}

// 🌍 LEADERBOARD ENTRY (для отправки на сервер)
type LeaderboardEntry struct {
	Username        string `json:"username"`
	TotalXP         int    `json:"total_xp"`
	Level           int    `json:"level"`
	League          string `json:"league"`
	CompletedTopics int    `json:"completed_topics"`
	CurrentStreak   int    `json:"current_streak"`
	LastUpdate      string `json:"last_update"`
}

var syllabus = []Topic{
	// LEVEL 1: Основы
	{Level: 1, Name: "Типы данных", Keywords: []string{"int", "float", "string", "bool"}, MinExamples: 3, XPReward: 50},
	{Level: 1, Name: "Переменные и константы", Keywords: []string{"var ", "const "}, MinExamples: 2, XPReward: 50},

	// LEVEL 2: Управление потоком
	{Level: 2, Name: "Условия (if/else)", Keywords: []string{"if ", "else"}, MinExamples: 2, XPReward: 75},
	{Level: 2, Name: "Циклы (for)", Keywords: []string{"for "}, MinExamples: 2, XPReward: 75},
	{Level: 2, Name: "Switch", Keywords: []string{"switch "}, MinExamples: 1, XPReward: 75},

	// LEVEL 3: Коллекции
	{Level: 3, Name: "Массивы и слайсы", Keywords: []string{"[]", "make([]", "append("}, MinExamples: 3, XPReward: 100},
	{Level: 3, Name: "Maps (карты)", Keywords: []string{"map[", "make(map"}, MinExamples: 2, XPReward: 100},

	// LEVEL 4: Функции
	{Level: 4, Name: "Функции", Keywords: []string{"func "}, MinExamples: 3, XPReward: 125},
	{Level: 4, Name: "Обработка ошибок", Keywords: []string{"error", "if err != nil"}, MinExamples: 2, XPReward: 125},

	// LEVEL 5: ООП в Go
	{Level: 5, Name: "Структуры", Keywords: []string{"type ", "struct"}, MinExamples: 2, XPReward: 150},
	{Level: 5, Name: "Методы", Keywords: []string{") func", "receiver"}, MinExamples: 2, XPReward: 150},
	{Level: 5, Name: "Интерфейсы", Keywords: []string{"interface"}, MinExamples: 1, XPReward: 150},

	// LEVEL 6: Concurrency
	{Level: 6, Name: "Горутины", Keywords: []string{"go func", "go "}, MinExamples: 1, XPReward: 200},
	{Level: 6, Name: "Каналы", Keywords: []string{"chan ", "<-"}, MinExamples: 2, XPReward: 200},

	// LEVEL 7: Продвинутое
	{Level: 7, Name: "HTTP сервер", Keywords: []string{"http.HandleFunc", "http.ListenAndServe"}, MinExamples: 1, XPReward: 250},
	{Level: 7, Name: "Тестирование", Keywords: []string{"func Test", "t.Error"}, MinExamples: 1, XPReward: 250},
}

// 🏆 Список всех достижений
var allAchievements = []Achievement{
	{ID: "first_commit", Name: "Первый шаг", Description: "Сделал первый коммит", Icon: "🎯", XPReward: 100},
	{ID: "week_streak", Name: "Огненная неделя", Description: "7 дней подряд", Icon: "🔥", XPReward: 300},
	{ID: "month_streak", Name: "Несгибаемый", Description: "30 дней подряд", Icon: "💪", XPReward: 1000},
	{ID: "level_3", Name: "Бронзовый воин", Description: "Достиг 3 уровня", Icon: "🥉", XPReward: 200},
	{ID: "level_5", Name: "Серебряный мастер", Description: "Достиг 5 уровня", Icon: "🥈", XPReward: 500},
	{ID: "level_7", Name: "Золотой гуру", Description: "Достиг 7 уровня", Icon: "🥇", XPReward: 1000},
	{ID: "maps_master", Name: "Картограф", Description: "Использовал maps 10+ раз", Icon: "🗺️", XPReward: 250},
	{ID: "concurrency_king", Name: "Повелитель потоков", Description: "Освоил горутины и каналы", Icon: "⚡", XPReward: 400},
	{ID: "error_handler", Name: "Страж ошибок", Description: "Обработал 20+ ошибок", Icon: "🛡️", XPReward: 300},
	{ID: "hundred_commits", Name: "Центурион", Description: "100 коммитов с Go кодом", Icon: "💯", XPReward: 2000},
}

func main() {
	fmt.Println("🔍 Начинаю анализ кода...")

	// Читаем статистику
	stats := loadStats()

	// Проверяем штрафы за пропуски
	applyPenalties(&stats)

	stats.TotalCommits++

	// Обновляем streak
	updateStreak(&stats)

	files := findGoFiles()
	if len(files) == 0 {
		fmt.Println("❌ Не найдено .go файлов")
		return
	}

	fmt.Printf("📂 Найдено файлов: %d\n", len(files))

	// Сбрасываем счётчики перед новым анализом
	for i := range syllabus {
		syllabus[i].Found = 0
	}

	// Анализируем файлы
	for _, file := range files {
		analyzeFile(file)
	}

	// Считаем прогресс и начисляем XP
	completed := 0
	totalTopics := len(syllabus)
	currentLevel := 1
	var nextTopic string
	xpGained := 0

	// Загружаем предыдущее состояние
	prevCompleted := loadPreviousState()

	for i := range syllabus {
		if syllabus[i].Found >= syllabus[i].MinExamples {
			// Проверяем, была ли тема изучена ранее
			wasCompleted := false
			for j := range prevCompleted {
				if prevCompleted[j] == syllabus[i].Name {
					wasCompleted = true
					break
				}
			}

			// Начисляем XP только за НОВЫЕ темы
			if !wasCompleted {
				xpGained += syllabus[i].XPReward
				fmt.Printf("✨ Новая тема изучена: %s (+%d XP)\n", syllabus[i].Name, syllabus[i].XPReward)
			}

			completed++
			if syllabus[i].Level > currentLevel {
				currentLevel = syllabus[i].Level
			}
		} else if nextTopic == "" {
			nextTopic = syllabus[i].Name
		}
	}

	// ВАЖНО: Проверяем если темы УДАЛЕНЫ (рефакторинг/удаление файлов)
	if completed < len(prevCompleted) {
		// Темы были удалены, но XP НЕ отнимаем (это честно заработано)
		fmt.Printf("⚠️ Внимание: %d тем больше не обнаружено в коде\n", len(prevCompleted)-completed)
		fmt.Println("💡 XP сохранён (рефакторинг не наказывается)")
	}

	if nextTopic == "" {
		nextTopic = "Все темы изучены! 🎉"
	}

	// Начисляем XP за streak
	if stats.CurrentStreak > 0 {
		streakXP := stats.CurrentStreak * 20
		xpGained += streakXP
		fmt.Printf("🔥 Streak бонус: +%d XP (%d дней)\n", streakXP, stats.CurrentStreak)
	}

	stats.TotalXP += xpGained
	stats.Level = currentLevel
	stats.CompletedTopics = completed

	// Определяем лигу
	stats.League = determineLeague(stats.Level, stats.TotalXP)

	// Проверяем достижения
	newAchievements := checkAchievements(&stats)

	// Начисляем XP за новые достижения
	for _, ach := range newAchievements {
		stats.TotalXP += ach.XPReward
		fmt.Printf("🏆 Достижение разблокировано: %s (+%d XP)\n", ach.Name, ach.XPReward)
	}

	percent := (float64(completed) / float64(totalTopics)) * 100

	// Сохраняем текущее состояние
	saveCurrentState(completed, &stats)

	// Сохраняем статистику
	saveStats(stats)

	// Генерируем отчёт
	message := generateReport(stats, percent, nextTopic, completed, totalTopics, newAchievements, xpGained)

	fmt.Println("\n" + message)

	// Отправляем в Telegram
	sendToTelegram(message)

	// Обновляем badges
	updateBadges(stats, percent)

	// Отправляем на центральный leaderboard
	sendToLeaderboard(stats)

	fmt.Println("\n✅ Анализ завершён!")
}

// 📊 Загрузка статистики
func loadStats() UserStats {
	data, err := os.ReadFile("stats.json")
	if err != nil {
		return UserStats{
			Username:       getUsername(),
			TotalXP:        0,
			CurrentStreak:  0,
			LongestStreak:  0,
			TotalCommits:   0,
			League:         "🥉 Bronze",
			LastCommitDate: "",
			Achievements:   []Achievement{},
			PenaltyDays:    0,
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

// 📝 Загрузка предыдущего состояния
func loadPreviousState() []string {
	data, err := os.ReadFile(".completed_topics")
	if err != nil {
		return []string{}
	}

	var topics []string
	json.Unmarshal(data, &topics)
	return topics
}

// 💾 Сохранение текущего состояния
func saveCurrentState(completed int, stats *UserStats) {
	var completedTopics []string
	for _, topic := range syllabus {
		if topic.Found >= topic.MinExamples {
			completedTopics = append(completedTopics, topic.Name)
		}
	}

	data, _ := json.Marshal(completedTopics)
	os.WriteFile(".completed_topics", data, 0644)
}

// ⚠️ Применение штрафов за пропуски
func applyPenalties(stats *UserStats) {
	if stats.LastCommitDate == "" {
		return
	}

	lastDate, _ := time.Parse("2006-01-02", stats.LastCommitDate)
	daysSince := int(time.Since(lastDate).Hours() / 24)

	if daysSince > 1 {
		stats.PenaltyDays = daysSince - 1
		penalty := stats.PenaltyDays * 30 // 30 XP за каждый пропущенный день
		stats.TotalXP -= penalty

		if stats.TotalXP < 0 {
			stats.TotalXP = 0
		}

		fmt.Printf("⚠️ Штраф: -%d XP за %d дней без коммитов\n", penalty, stats.PenaltyDays)
	} else {
		stats.PenaltyDays = 0
	}
}

// 🔥 Обновление streak
func updateStreak(stats *UserStats) {
	today := time.Now().Format("2006-01-02")

	if stats.LastCommitDate == "" {
		stats.CurrentStreak = 1
		stats.LongestStreak = 1
	} else {
		lastDate, _ := time.Parse("2006-01-02", stats.LastCommitDate)
		daysDiff := int(time.Since(lastDate).Hours() / 24)

		if daysDiff == 1 {
			stats.CurrentStreak++
			if stats.CurrentStreak > stats.LongestStreak {
				stats.LongestStreak = stats.CurrentStreak
			}
		} else if daysDiff > 1 {
			stats.CurrentStreak = 1
		}
	}

	stats.LastCommitDate = today
}

// 🏆 Определение лиги
func determineLeague(level, xp int) string {
	if level >= 7 || xp >= 3000 {
		return "💎 Diamond"
	} else if level >= 5 || xp >= 2000 {
		return "🥇 Gold"
	} else if level >= 3 || xp >= 1000 {
		return "🥈 Silver"
	}
	return "🥉 Bronze"
}

// 🏆 Проверка достижений
func checkAchievements(stats *UserStats) []Achievement {
	var newAchievements []Achievement

	for _, achievement := range allAchievements {
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
		if strings.Contains(path, "notifier") || strings.Contains(path, ".git") {
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
				fmt.Printf("  ✓ '%s': %d раз\n", keyword, count)
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
func generateReport(stats UserStats, percent float64, nextTopic string, completed, total int, newAchievements []Achievement, xpGained int) string {
	barWidth := 10
	filled := int((percent / 100) * float64(barWidth))
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "▰"
		} else {
			bar += "▱"
		}
	}

	levelName := getLevelName(stats.Level)

	// Основной отчёт
	var report strings.Builder
	report.WriteString("🎮 GO LEARNING TRACKER\n\n")

	// Информация о пользователе
	report.WriteString(fmt.Sprintf("👤 %s\n", stats.Username))
	report.WriteString(fmt.Sprintf("⚡ Level %d · %s · %d XP", stats.Level, levelName, stats.TotalXP))
	if xpGained > 0 {
		report.WriteString(fmt.Sprintf(" *(+%d)*", xpGained))
	}
	report.WriteString("\n")
	report.WriteString(fmt.Sprintf("🛡 %s\n\n", stats.League))

	// Прогресс бар
	report.WriteString(fmt.Sprintf("%s %.0f%%\n", bar, percent))
	report.WriteString(fmt.Sprintf("%d/%d тем · %d коммитов\n", completed, total, stats.TotalCommits))

	// Streak (если >= 3 дней)
	if stats.CurrentStreak >= 3 {
		report.WriteString(fmt.Sprintf("\n🔥 Огненная серия: %d дней подряд", stats.CurrentStreak))
		if stats.CurrentStreak >= 30 {
			report.WriteString(" — Легенда!")
		} else if stats.CurrentStreak >= 14 {
			report.WriteString(" — Невероятно!")
		} else if stats.CurrentStreak >= 7 {
			report.WriteString(" — Отлично!")
		}
		report.WriteString("\n")
	}

	// Штрафы
	if stats.PenaltyDays > 0 {
		report.WriteString(fmt.Sprintf("\n⚠️ Потеря концентрации: -%d XP (%d дней без практики)\n", stats.PenaltyDays*30, stats.PenaltyDays))
	}

	// Новые достижения
	if len(newAchievements) > 0 {
		report.WriteString("\n🎉 Новое достижение разблокировано!\n")
		for _, ach := range newAchievements {
			report.WriteString(fmt.Sprintf("%s %s *(+%d XP)*\n", ach.Icon, ach.Name, ach.XPReward))
		}
	}

	// Следующая цель
	report.WriteString(fmt.Sprintf("\n🎯 Следующая цель: %s\n", nextTopic))

	// Изученные навыки (только текущий и следующий уровень)
	report.WriteString("\nИзучено:\n")

	showLevels := []int{stats.Level}
	if stats.Level < 7 {
		showLevels = append(showLevels, stats.Level+1)
	}

	shownCount := 0
	maxShow := 5 // Показываем максимум 5 тем

	for _, lvl := range showLevels {
		for _, topic := range syllabus {
			if topic.Level == lvl && shownCount < maxShow {
				if topic.Found >= topic.MinExamples {
					report.WriteString(fmt.Sprintf("  ✓ %s\n", topic.Name))
				} else {
					report.WriteString(fmt.Sprintf("  → %s\n", topic.Name))
				}
				shownCount++
			}
		}
	}

	report.WriteString("\n#golang #buildinpublic\n")

	return report.String()
}

// 🏆 Название уровня (Фэнтези стиль)
func getLevelName(level int) string {
	names := map[int]string{
		1: "Новобранец 🌱",
		2: "Подмастерье ⚔️",
		3: "Искатель 🗡️",
		4: "Следопыт 🏹",
		5: "Чародей 🔮",
		6: "Архимаг ⚡",
		7: "Великий Магистр 👑",
	}
	if name, ok := names[level]; ok {
		return name
	}
	return "Новобранец 🌱"
}

// 🎨 Обновление badges
func updateBadges(stats UserStats, percent float64) {
	levelBadge := fmt.Sprintf("![Level](https://img.shields.io/badge/Level-%d-blue)", stats.Level)
	progressBadge := fmt.Sprintf("![Progress](https://img.shields.io/badge/Progress-%.0f%%25-brightgreen)", percent)
	streakBadge := fmt.Sprintf("![Streak](https://img.shields.io/badge/Streak-%d_days-orange)", stats.CurrentStreak)
	xpBadge := fmt.Sprintf("![XP](https://img.shields.io/badge/XP-%d-purple)", stats.TotalXP)
	leagueBadge := fmt.Sprintf("![League](https://img.shields.io/badge/League-%s-gold)", strings.ReplaceAll(stats.League, " ", "_"))

	readmeContent, err := os.ReadFile("README.md")
	if err != nil {
		return
	}

	content := string(readmeContent)

	badgesSection := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		levelBadge, progressBadge, streakBadge, xpBadge, leagueBadge,
	)

	if strings.Contains(content, "![Level]") {
		re := regexp.MustCompile(`!\[Level\].*\n!\[Progress\].*\n!\[Streak\].*\n!\[XP\].*\n!\[League\].*`)
		content = re.ReplaceAllString(content, badgesSection)
	} else {
		lines := strings.Split(content, "\n")
		if len(lines) > 0 {
			lines = append(lines[:1], append([]string{"", badgesSection, ""}, lines[1:]...)...)
			content = strings.Join(lines, "\n")
		}
	}

	os.WriteFile("README.md", []byte(content), 0644)
	fmt.Println("✅ Badges обновлены")
}

// 🌍 Отправка на центральный leaderboard
func sendToLeaderboard(stats UserStats) {
	webhookURL := os.Getenv("LEADERBOARD_WEBHOOK")
	if webhookURL == "" {
		fmt.Println("⚠️ LEADERBOARD_WEBHOOK не настроен (пропускаю)")
		return
	}

	entry := LeaderboardEntry{
		Username:        stats.Username,
		TotalXP:         stats.TotalXP,
		Level:           stats.Level,
		League:          stats.League,
		CompletedTopics: stats.CompletedTopics,
		CurrentStreak:   stats.CurrentStreak,
		LastUpdate:      time.Now().Format("2006-01-02 15:04:05"),
	}

	jsonData, _ := json.Marshal(entry)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("⚠️ Ошибка отправки на leaderboard: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Данные отправлены на leaderboard")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("⚠️ Leaderboard ответил %d: %s\n", resp.StatusCode, string(body))
	}
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
		fmt.Println("⚠️ Telegram токены не найдены")
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
