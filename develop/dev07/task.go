package cut

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Options задает параметры утилиты
type Options struct {
	fields    []int  // Колонки для выбора
	delimiter string // Разделитель
	separated bool   // Только строки с разделителем
}

func main() {
	// Определение флагов
	fields := flag.String("f", "", "выбрать поля (колонки), указать через запятую")
	delimiter := flag.String("d", "\t", "использовать другой разделитель (по умолчанию TAB)")
	separated := flag.Bool("s", false, "только строки с разделителем")

	flag.Parse()

	// Проверка, указаны ли колонки
	if *fields == "" {
		fmt.Fprintln(os.Stderr, "Error: -f flag is required")
		os.Exit(1)
	}

	// Разбор колонок
	fieldIndexes, err := parseFields(*fields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing fields: %v\n", err)
		os.Exit(1)
	}

	options := Options{
		fields:    fieldIndexes,
		delimiter: *delimiter,
		separated: *separated,
	}

	// Чтение из STDIN
	if err := cut(os.Stdin, os.Stdout, options); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFields(fields string) ([]int, error) {
	parts := strings.Split(fields, ",")
	var result []int
	for _, part := range parts {
		var index int
		_, err := fmt.Sscanf(part, "%d", &index)
		if err != nil || index <= 0 {
			return nil, fmt.Errorf("invalid field: %s", part)
		}
		result = append(result, index-1) // Преобразуем к 0-индексации
	}
	return result, nil
}

func cut(input io.Reader, output io.Writer, options Options) error {
	scanner := bufio.NewScanner(input)
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()

		// Разделение строки
		columns := strings.Split(line, options.delimiter)

		// Фильтрация строк без разделителя
		if options.separated && len(columns) < 2 {
			continue
		}

		// Выбор указанных колонок
		var selected []string
		for _, index := range options.fields {
			if index < len(columns) {
				selected = append(selected, columns[index])
			}
		}

		// Печать результата
		if len(selected) > 0 {
			fmt.Fprintln(writer, strings.Join(selected, options.delimiter))
		}
	}

	return scanner.Err()
}
