package main

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

import (
	"fmt"
	"time"
)

// or объединяет один или более done-каналов в один.
func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		// Если нет каналов, просто завершаем
		if len(channels) == 0 {
			return
		}

		// Используем select для ожидания первого закрытого канала
		select {
		case <-channels[0]:
			// Если первый канал закрыт, завершаем горутину
			return
		case <-or(channels[1:]...): // Рекурсивно вызываем or для остальных каналов
			return
		}
	}()
	return out
}

func main() {
	// Функция, которая генерирует канал, который закрывается через заданное время
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	// Запускаем несколько каналов с разными временными задержками
	start := time.Now()
	<-or(
		sig(2*time.Second),
		sig(5*time.Minute),
		sig(10*time.Second),
		sig(1*time.Second),
		sig(1*time.Minute),
	)

	// Выводим время, через которое завершилось выполнение
	fmt.Printf("Done after %v\n", time.Since(start))
}
