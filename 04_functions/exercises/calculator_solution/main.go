package main

import "fmt"

func calculate(a, b int) (int, int, int, float64) {
	return a + b, a - b, a * b, float64(a) / float64(b)
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		result := a
		a, b = b, a+b
		return result
	}
}

func process() {
	defer fmt.Println("終了（defer で実行）")
	fmt.Println("開始")
	fmt.Println("処理中...")
}

func main() {
	s, d, p, q := calculate(10, 3)
	fmt.Printf("和:%d 差:%d 積:%d 商:%.3f\n", s, d, p, q)

	fmt.Println(sum(1, 2, 3, 4, 5))
	nums := []int{10, 20, 30}
	fmt.Println(sum(nums...))

	process()

	fib := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", fib())
	}
	fmt.Println()
}
