package unpack

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// UnpackString распаковывает строку с повторяющимися символами.
func UnpackString(input string) (string, error) {
	if input == "" {
		return "", nil // Пустая строка возвращает пустую строку.
	}

	var result strings.Builder
	var prevRune rune
	escaped := false

	for i, r := range input {
		if unicode.IsDigit(r) && !escaped {
			if prevRune == 0 {
				return "", errors.New("строка начинается с цифры")
			}

			// Чтение многозначного числа
			numStr := string(r)
			for j := i + 1; j < len(input); j++ {
				if unicode.IsDigit(rune(input[j])) {
					numStr += string(input[j])
					i++
				} else {
					break
				}
			}

			repeatCount, err := strconv.Atoi(numStr)
			if err != nil {
				return "", err
			}
			result.WriteString(strings.Repeat(string(prevRune), repeatCount-1))
			prevRune = 0
		} else if r == '\\' && !escaped {
			// Обработка начала escape-последовательности
			escaped = true
		} else {
			if escaped {
				// Экранированный символ
				if prevRune != 0 {
					result.WriteRune(prevRune)
				}
				prevRune = r
				escaped = false
			} else {
				// Обычный символ
				if prevRune != 0 {
					result.WriteRune(prevRune)
				}
				prevRune = r
			}
		}
	}

	// Записываем последний символ, если он есть
	if prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}
