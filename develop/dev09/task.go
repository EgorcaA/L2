package main

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	for {
		// Вывод приглашения
		fmt.Print("myshell> ")

		// Чтение ввода пользователя
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		// Если введено \quit, выходим из программы
		if input == "\\quit" {
			break
		}

		// Парсим команду
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		// Обработка команд
		switch args[0] {
		case "cd":
			if len(args) < 2 {
				fmt.Println("cd: missing argument")
			} else {
				if err := os.Chdir(args[1]); err != nil {
					fmt.Println("cd:", err)
				}
			}
		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println("pwd:", err)
			} else {
				fmt.Println(dir)
			}
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))
		case "kill":
			if len(args) < 2 {
				fmt.Println("kill: missing argument")
			} else {
				pid := args[1]
				cmd := exec.Command("kill", pid)
				if err := cmd.Run(); err != nil {
					fmt.Println("kill:", err)
				}
			}
		case "ps":
			cmd := exec.Command("ps", "-aux")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("ps:", err)
			}
		default:
			// Обработка пайпа
			commands := strings.Split(input, "|")
			if len(commands) > 1 {
				handlePipe(commands)
			} else {
				handleSimpleCommand(args)
			}
		}
	}
}

// Функция для обработки команд с пайпами
func handlePipe(commands []string) {
	var cmdArr []*exec.Cmd
	for _, command := range commands {
		args := strings.Fields(strings.TrimSpace(command))
		if len(args) > 0 {
			cmdArr = append(cmdArr, exec.Command(args[0], args[1:]...))
		}
	}

	// Связываем пайпы
	for i := 0; i < len(cmdArr)-1; i++ {
		pr, pw := pipe()
		cmdArr[i].Stdout = pw
		cmdArr[i+1].Stdin = pr
	}

	// Выполняем все команды
	for _, cmd := range cmdArr {
		if err := cmd.Start(); err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// Ждем завершения всех команд
	for _, cmd := range cmdArr {
		cmd.Wait()
	}
}

// Функция для выполнения одной команды
func handleSimpleCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

// Функция для создания pipe
func pipe() (*os.File, *os.File) {
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("Error creating pipe:", err)
	}
	return r, w
}
