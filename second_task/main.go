package main

import (
	"fmt"
	"sync"
)

// Merge принимает два read-only канала и возвращает выходной канал,
// в который последовательно (в любом порядке) будут отправлены все значения
// из обоих входных каналов.
//
// Выходной канал должен быть закрыт после того, как оба входных канала закроются.
// Merge не должен закрывать входные каналы
//
// Для проверки решения запустите тесты: go test -v
func Merge(ch1, ch2 <-chan int) <-chan int {
	res := make(chan int, 100)

	// Реализовал через горутины для равномерного получения из двух каналов
	// Хорошо работает, если данные поступают неравномерно,
	// плохо, если чтение результата происходит медленно и буфер переполнится
	go func() {
		defer close(res)
		var wg sync.WaitGroup
		wg.Add(2)

		mergeChan := func(ch <-chan int) {
			for val := range ch {
				res <- val
			}
			wg.Done()
		}
		mergeChan(ch1)
		mergeChan(ch2)
		wg.Wait()
	}()

	return res
}

func main() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		defer close(a)
		a <- 4
		a <- 1
	}()

	go func() {
		defer close(b)
		b <- 2
		b <- 4
	}()

	for v := range Merge(a, b) {
		fmt.Println(v)
	}
}
