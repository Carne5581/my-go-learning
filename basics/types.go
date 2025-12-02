package main

import (
	"fmt"
	"unsafe" // Пакет для "низкоуровневых" дел, показывает размер в памяти
)

func main() {
	fmt.Println("=== ИЗУЧАЕМ ТИПЫ ДАННЫХ ===")

	// 1. ЦЕЛЫЕ ЧИСЛА (Integers)
	// int - стандартный тип, его размер зависит от процессора (обычно 64 бита = 8 байт)
	var age int = 25
	fmt.Printf("Переменная 'age': тип=%T, значение=%d, занимает байт=%d\n", age, age, unsafe.Sizeof(age))

	// int8 - очень маленький тип (только от -128 до 127). Экономит память.
	var smallNum int8 = 100
	fmt.Printf("Переменная 'smallNum': тип=%T, значение=%d, занимает байт=%d\n", smallNum, smallNum, unsafe.Sizeof(smallNum))

	// 2. ДРОБНЫЕ ЧИСЛА (Float)
	var pi float64 = 3.1415926535
	fmt.Printf("Число Pi: тип=%T, значение=%.2f (обрезали при выводе)\n", pi, pi)

	// 3. СТРОКИ (String)
	// Внимание: Английская буква = 1 байт. Русская буква = 2 байта (UTF-8).
	nameEn := "Marsel"
	nameRu := "Марсель"

	fmt.Println("\n--- Сравнение строк ---")
	fmt.Printf("Имя '%s', длина: %d байт\n", nameEn, len(nameEn))
	fmt.Printf("Имя '%s', длина: %d байт (Видишь? Букв 7, а байт 14)\n", nameRu, len(nameRu))

	// 4. ЛОГИЧЕСКИЙ ТИП (Boolean) - Правда или Ложь
	isProgrammer := true
	fmt.Printf("\nТы программист? %t, занимает байт=%d\n", isProgrammer, unsafe.Sizeof(isProgrammer))
}
