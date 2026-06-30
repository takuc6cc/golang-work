package main

import (
	"fmt"
	"sync"
)

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func generate(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
}

func square(in <-chan int, out chan<- int) {
	for n := range in {
		out <- n * n
	}
	close(out)
}

func main() {
	// 問題1: WaitGroup
	fmt.Println("=== WaitGroup ===")
	var wg sync.WaitGroup
	messages := []string{"Hello from goroutine 1", "Hello from goroutine 2", "Hello from goroutine 3"}
	for _, msg := range messages {
		wg.Add(1)
		go func(m string) {
			defer wg.Done()
			fmt.Println(m)
		}(msg)
	}
	wg.Wait()
	fmt.Println("すべての goroutine 完了")

	// 問題2: パイプライン
	fmt.Println("\n=== パイプライン ===")
	naturals := make(chan int)
	squares := make(chan int)
	go generate(naturals)
	go square(naturals, squares)
	for n := range squares {
		fmt.Printf("%d\n", n)
	}

	// 問題3: goroutine セーフカウンター
	fmt.Println("\n=== goroutine セーフカウンター ===")
	counter := &SafeCounter{}
	var wg2 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}
	wg2.Wait()
	fmt.Printf("最終カウント: %d（期待値: 10000）\n", counter.Value())
}
