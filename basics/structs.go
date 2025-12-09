package main

import "fmt"

// 1. ОПРЕДЕЛЕНИЕ СТРУКТУРЫ (ЧЕРТЕЖ)
// Мы создаем новый тип данных "Hero"
// Обрати внимание: Имена полей с Большой буквы (чтобы их видели другие пакеты)
type Hero struct {
	Nickname string
	Level    int
	Health   int
	IsAlive  bool
	Weapon   string
}

func main() {
	fmt.Println("=== RPG GAME SYSTEM v1.0 ===")

	// 2. СОЗДАНИЕ ОБЪЕКТА (Инстанс)
	// Мы строим дом по чертежу
	player := Hero{
		Nickname: "Arthas",
		Level:    80,
		Health:   100,
		IsAlive:  true,
		Weapon:   "Frostmourne",
	}

	fmt.Printf("Герой создан: %+v\n", player)
	// %+v - спец. формат, показывает имена полей (очень удобно для отладки!)

	// 3. ДОСТУП К ПОЛЯМ (Через точку)
	fmt.Println("\n--- ИНФО О ГЕРОЕ ---")
	fmt.Println("Имя:", player.Nickname)
	fmt.Println("HP:", player.Health)

	// 4. ИЗМЕНЕНИЕ ПОЛЕЙ
	fmt.Println("\n--- БИТВА НАЧАЛАСЬ! ---")
	fmt.Println("Герой получает удар на 40 урона!")

	player.Health = player.Health - 40

	fmt.Printf("Текущее здоровье: %d\n", player.Health)

	// Провека на жизнь
	if player.Health <= 0 {
		player.IsAlive = false
		fmt.Println("Герой пал...")
	} else {
		fmt.Println("Герой ранен, но жив!")
	}

	player.Weapon = "Axe"

	fmt.Println("\nГерой сменил оружие на:", player.Weapon)

	// 5. КОПИРОВАНИЕ СТРУКТУР (Важный момент!)
	fmt.Println("\n--- МАГИЯ КОПИРОВАНИЯ ---")
	// Создаем "Тень" героя
	// В Go привание структуры = КОПИРОВАНИЕ всех полей
	shadow := player

	shadow.Nickname = "Shadow Arthas"
	shadow.Health = 1000 // Тень сильнее

	fmt.Println("Оригинал:", player.Nickname, player.Health)
	fmt.Println("Тень: ", shadow.Nickname, shadow.Health)

	// Вывод: Изменение копии (shadow) НЕ ПОВЛИЯЛО на оригинал (player).
	// Это то самое Pass by Value, о котором мы говорили.
}
