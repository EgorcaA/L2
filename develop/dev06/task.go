package main

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Options задает флаги фильтрации
type Options struct {
	after       int  // -A
	before      int  // -B
	context     int  // -C
	count       bool // -c
	ignoreCase  bool // -i
	invert      bool // -v
	fixed       bool // -F
	lineNumbers bool // -n
}

func main() {
	// Парсинг аргументов
	after := flag.Int("A", 0, "печатать +N строк после совпадения")
	before := flag.Int("B", 0, "печатать +N строк до совпадения")
	context := flag.Int("C", 0, "печатать ±N строк вокруг совпадения")
	count := flag.Bool("c", false, "количество строк")
	ignoreCase := flag.Bool("i", false, "игнорировать регистр")
	invert := flag.Bool("v", false, "исключать совпадения")
	fixed := flag.Bool("F", false, "точное совпадение со строкой")
	lineNumbers := flag.Bool("n", false, "печатать номер строки")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: grep [OPTIONS] <pattern> <file>")
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	fileName := flag.Arg(1)

	options := Options{
		after:       *after,
		before:      *before,
		context:     *context,
		count:       *count,
		ignoreCase:  *ignoreCase,
		invert:      *invert,
		fixed:       *fixed,
		lineNumbers: *lineNumbers,
	}

	if err := grep(pattern, fileName, options); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func grep(pattern, fileName string, options Options) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if options.context > 0 {
		options.after = options.context
		options.before = options.context
	}

	matches := make(map[int]bool)

	for i, line := range lines {
		lineToMatch := line
		patternToMatch := pattern

		if options.ignoreCase {
			lineToMatch = strings.ToLower(line)
			patternToMatch = strings.ToLower(pattern)
		}

		match := false
		if options.fixed {
			match = lineToMatch == patternToMatch
		} else {
			match = strings.Contains(lineToMatch, patternToMatch)
		}

		if options.invert {
			match = !match
		}

		if match {
			matches[i] = true
		}
	}

	if options.count {
		fmt.Println(len(matches))
		return nil
	}

	// Формирование вывода с учетом флагов -A, -B, -C
	output := make(map[int]bool)
	for lineNum := range matches {
		start := max(0, lineNum-options.before)
		end := min(len(lines)-1, lineNum+options.after)
		for i := start; i <= end; i++ {
			output[i] = true
		}
	}

	// Печать результата
	for i := 0; i < len(lines); i++ {
		if output[i] {
			if options.lineNumbers {
				fmt.Printf("%d: %s\n", i+1, lines[i])
			} else {
				fmt.Println(lines[i])
			}
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func grepLines(pattern string, lines []string, options Options) ([]string, error) {
	if options.context > 0 {
		options.after = options.context
		options.before = options.context
	}

	matches := make(map[int]bool)

	for i, line := range lines {
		lineToMatch := line
		patternToMatch := pattern

		if options.ignoreCase {
			lineToMatch = strings.ToLower(line)
			patternToMatch = strings.ToLower(pattern)
		}

		match := false
		if options.fixed {
			match = lineToMatch == patternToMatch
		} else {
			match = strings.Contains(lineToMatch, patternToMatch)
		}

		if options.invert {
			match = !match
		}

		if match {
			matches[i] = true
		}
	}

	if options.count {
		return []string{fmt.Sprintf("%d", len(matches))}, nil
	}

	output := make(map[int]bool)
	for lineNum := range matches {
		start := max(0, lineNum-options.before)
		end := min(len(lines)-1, lineNum+options.after)
		for i := start; i <= end; i++ {
			output[i] = true
		}
	}

	var result []string
	for i := 0; i < len(lines); i++ {
		if output[i] {
			if options.lineNumbers {
				result = append(result, fmt.Sprintf("%d: %s", i+1, lines[i]))
			} else {
				result = append(result, lines[i])
			}
		}
	}

	return result, nil
}
