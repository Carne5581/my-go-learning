package main

import "fmt"

// --- 1. КОНТРАКТ (ИНТЕРФЕЙС) ---
// Мы объявляем: "Любой, кто хочет называться DamageTaker,
// ОБЯЗАН иметь метод TakeDamage(int)"
type DamageTaker interface {
	TakeDamage(amount int)
}

// --- 2. РАЗНЫЕ СТРУКТУРЫ ---

// Структура 1: Герой
type Hero struct {
	Name   string
	Health int
}

// Реализуем метод TakeDamage для Героя
// (Теперь Герой автоматически считается DamageTaker!)
func (h *Hero) TakeDamage(amount int) {
	h.Health -= amount
	fmt.Printf("🧍 Герой %s получил по лицу! Осталось HP: %d\n", h.Name, h.Health)
}

// Структура 2: Дракон
type Dragon struct {
	HP       int
	IsFlying bool
}

// Структура 3: Стена
type Wall struct {
	Durability int
}

// Реализуем метод TakeDamage для Дракона
// У дракона логика другая (он толстокожий)
func (d *Dragon) TakeDamage(amount int) {
	d.HP -= amount
	if d.HP > 0 {
		fmt.Printf("🐲 Дракон ревет! У него осталось %d HP\n", d.HP)
	} else {
		fmt.Println("🐲 Дракон пал на землю!")
		d.IsFlying = false
	}
}

// Метод TakeDamage для Стены
func (h *Wall) TakeDamage(amount int) {
	h.Durability -= amount
	if h.Durability <= 0 {
		fmt.Printf("🧱 Стена разрушена в крошку!")
	}
}

// --- 3. УНИВЕРСАЛЬНАЯ ФУНКЦИЯ (ПОЛИМОРФИЗМ) ---
// Внимание! Эта функция принимает HE Hero и HE Dragon.
// Она принимает ИНТЕРФЕЙС.
// Ей можно скормить кого угодно, кто умеет получать урон.
func AttackSomething(target DamageTaker, damage int) {
	// Мы даже не знаем, кто такой target.
	// Но мы точно знаем, что у него есть метод TakeDamage.
	target.TakeDamage(damage)
}

func main() {
	fmt.Println("=== BATTLEFIELD v2.0 (INTERFACES) ===")

	// Создаем конкретные объекты
	arthas := &Hero{Name: "Arthas", Health: 100} // & нужен, т.к. методы используют *Hero
	smaug := &Dragon{HP: 500, IsFlying: true}
	stoneWall := &Wall{Durability: 20}

	// АТАКА!
	// Смотри: мы используем одну и ту же функцию AttackSomething

	fmt.Println("\n--- Бьем Героя ---")
	AttackSomething(arthas, 30)

	fmt.Println("\n--- Бьем Дракона ---")
	AttackSomething(smaug, 500)

	fmt.Println("\n--- Бьем Стену ---")
	AttackSomething(stoneWall, 50)

	// В Go можно создать слайс из интерфейсов!
	// Сборная солянка: список тех, кого можно бить
	fmt.Println("\n--- МАССОВАЯ АТАКА ---")
	enemies := []DamageTaker{arthas, smaug}

	for _, enemy := range enemies {
		AttackSomething(enemy, 50)
	}
}
