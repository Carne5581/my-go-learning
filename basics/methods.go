package main

import "fmt"

// 1. СТРУКТУРА
type Hero struct {
	Nickname string
	Health   int
	Weapon   string
}

// 2. МЕТОДЫ (Поведение)

// Метод БЕЗ звездочки (Valuer Receiver)
// Он работает с КОПИЕЙ. Подходит "чтения" данных.
// Синтаксис: func (переменная Тип) ИмяМетода()
func (h Hero) Introduce() {
	fmt.Printf("Привет! Я %s, у меня %d HP и оружие %s.\n", h.Nickname, h.Health, h.Weapon)
}

// Метод СО ЗВЕЗДОЧКОЙ (Pointer Receiver)
// Он работает с ОРИГИНАЛОМ. Нужен для "изменения" данных.
// h *Hero означает: "h - это адрес в памяти, где лежит герой"
func (h *Hero) TakeDamage(damage int) {
	fmt.Printf("%s получает удар на %d урона!\n", h.Nickname, damage)

	// Мы меняем поле реального объекта по адресу
	h.Health = h.Health - damage

	if h.Health < 0 {
		h.Health = 0 // Чтобы здоровье не ушло в минус
	}
}

func (h *Hero) HealHero(amount int) {
	fmt.Printf("Принял зелье здоровья, получил %d здоровья!\n", h.Health)

	h.Health = h.Health + amount

	if h.Health > 100 {
		h.Health = 100
		fmt.Println("(Здоровье полное!)")
	}
}

// Еще один метод с Pointer Receiver - Смена оружия
func (h *Hero) EquipWeapon(newWeapon string) {
	fmt.Printf("%s меняет оружие: %s -> %s\n", h.Nickname, h.Weapon, newWeapon)
	h.Weapon = newWeapon
}

func main() {
	fmt.Println("=== RPG METHODS SYSTEM ===")

	// Создаем героя
	player := Hero{
		Nickname: "Arthas",
		Health:   100,
		Weapon:   "Sword",
	}

	// 1. Вызываем простой метод (чтение)
	player.Introduce()

	// 2. Вызываем метод изменения (Pointer)
	// Go сам понимает, что нужно передать адрес, тебе не нужно писать (&player). TakeDamage
	// Это называется "Syntactic Sugar" (Синтаксический сахар)
	player.TakeDamage(80)
	player.HealHero(50)
	player.HealHero(1000)

	// Проверяем, изменилось ли здоровье в main
	fmt.Println("Здоровье в main после удара:", player.Health)
	// Должно быть 60! Если бы мы забыли звездочку в методе, осталось бы 100.

	// 3. Смена оружия
	player.EquipWeapon("Frostmourne")
	player.Introduce() // Проверяем статус
}
