package main

import (
	"fmt"
)

func main() {
	num := 47

	results := make(chan int, num)
	jobs := make(chan int, num)

	for i := 0; i < num; i++ {
		jobs <- i
	}

	for cur_worker := 0; cur_worker < 100; cur_worker++ {
		go worker(results, jobs)
	}

	for i := 0; i < num; i++ {
		fmt.Println(i, <-results)
	}
}

func worker(results chan<- int, jobs <-chan int) {
	for i := range jobs {
		results <- fib(i)
	}
}

func fib(n int) int {
	if n <= 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
