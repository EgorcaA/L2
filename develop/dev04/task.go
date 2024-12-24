package main

/*
=== Утилита sort ===
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры):
на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:
    -k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел);
    -n — сортировать по числовому значению;
    -r — сортировать в обратном порядке;
    -u — не выводить повторяющиеся строки.
Дополнительно
Реализовать поддержку утилитой следующих ключей:
    -M — сортировать по названию месяца;
    -b — игнорировать хвостовые пробелы;
    -c — проверять отсортированы ли данные;
    -h — сортировать по числовому значению с учетом суффиксов.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// SortOptions содержит флаги сортировки
type SortOptions struct {
	column       int
	numeric      bool
	reverse      bool
	unique       bool
	month        bool
	ignoreSpaces bool
	checkSorted  bool
	humanNumeric bool
}

func main() {
	// Парсинг аргументов
	column := flag.Int("k", 0, "колонка для сортировки (начиная с 1)")
	numeric := flag.Bool("n", false, "сортировка по числовому значению")
	reverse := flag.Bool("r", false, "обратный порядок сортировки")
	unique := flag.Bool("u", false, "удаление повторяющихся строк")
	month := flag.Bool("M", false, "сортировка по названию месяца")
	ignoreSpaces := flag.Bool("b", false, "игнорировать хвостовые пробелы")
	checkSorted := flag.Bool("c", false, "проверить, отсортированы ли данные")
	humanNumeric := flag.Bool("h", false, "сортировка с учетом числовых суффиксов")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: sort [OPTIONS] <input-file> [output-file]")
		os.Exit(1)
	}

	inputFile := flag.Arg(0)
	outputFile := ""
	if len(flag.Args()) > 1 {
		outputFile = flag.Arg(1)
	}

	options := SortOptions{
		column:       *column - 1,
		numeric:      *numeric,
		reverse:      *reverse,
		unique:       *unique,
		month:        *month,
		ignoreSpaces: *ignoreSpaces,
		checkSorted:  *checkSorted,
		humanNumeric: *humanNumeric,
	}

	lines, err := readLines(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	if options.checkSorted {
		if isSorted(lines, options) {
			fmt.Println("File is sorted.")
		} else {
			fmt.Println("File is not sorted.")
		}
		os.Exit(0)
	}

	sortedLines, err := sortLines(lines, options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during sorting: %v\n", err)
		os.Exit(1)
	}

	if outputFile == "" {
		for _, line := range sortedLines {
			fmt.Println(line)
		}
	} else {
		if err := writeLines(outputFile, sortedLines); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
	}
}

// readLines читает строки из файла
func readLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines записывает строки в файл
func writeLines(fileName string, lines []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// sortLines сортирует строки с учетом заданных опций
func sortLines(lines []string, options SortOptions) ([]string, error) {
	if options.ignoreSpaces {
		for i, line := range lines {
			lines[i] = strings.TrimRight(line, " ")
		}
	}

	if options.unique {
		lines = unique(lines)
	}

	sort.Slice(lines, func(i, j int) bool {
		left, right := lines[i], lines[j]

		if options.column >= 0 {
			leftParts := strings.Fields(left)
			rightParts := strings.Fields(right)
			if len(leftParts) > options.column {
				left = leftParts[options.column]
			} else {
				left = ""
			}
			if len(rightParts) > options.column {
				right = rightParts[options.column]
			} else {
				right = ""
			}
		}

		if options.numeric {
			leftNum, err1 := strconv.ParseFloat(left, 64)
			rightNum, err2 := strconv.ParseFloat(right, 64)
			if err1 == nil && err2 == nil {
				if options.reverse {
					return leftNum > rightNum
				}
				return leftNum < rightNum
			}
		}

		if options.reverse {
			return left > right
		}
		return left < right
	})

	return lines, nil
}

// unique удаляет повторяющиеся строки
func unique(lines []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, line := range lines {
		if _, exists := seen[line]; !exists {
			result = append(result, line)
			seen[line] = struct{}{}
		}
	}
	return result
}

// isSorted проверяет, отсортированы ли строки
func isSorted(lines []string, options SortOptions) bool {
	for i := 1; i < len(lines); i++ {
		if options.reverse {
			if lines[i-1] < lines[i] {
				return false
			}
		} else {
			if lines[i-1] > lines[i] {
				return false
			}
		}
	}
	return true
}
