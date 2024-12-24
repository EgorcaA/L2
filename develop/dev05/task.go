package anagrams

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"sort"
	"strings"
)

// FindAnagrams принимает массив слов и возвращает мапу множеств анаграмм.
func FindAnagrams(words []string) map[string][]string {
	anagrams := make(map[string][]string)
	seenWords := make(map[string]bool)

	// Преобразование всех слов в нижний регистр
	for _, word := range words {
		wordLower := strings.ToLower(word)
		sortedWord := sortString(wordLower)

		// Проверка, был ли уже добавлен аналогичный ключ
		if _, exists := anagrams[sortedWord]; !exists {
			anagrams[sortedWord] = []string{}
		}

		if !seenWords[wordLower] {
			anagrams[sortedWord] = append(anagrams[sortedWord], wordLower)
			seenWords[wordLower] = true
		}
	}

	// Формирование итоговой мапы
	result := make(map[string][]string)
	for _, group := range anagrams {
		if len(group) > 1 { // Исключение множеств с одним элементом
			sort.Strings(group) // Сортируем массив анаграмм
			result[group[0]] = group
		}
	}

	return result
}

// sortString сортирует символы строки по алфавиту.
func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}
