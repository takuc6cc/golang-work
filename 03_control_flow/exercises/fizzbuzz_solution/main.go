package main

import "fmt"

func main() {
	// 問題1: FizzBuzz
	fmt.Println("=== FizzBuzz ===")
	for i := 1; i <= 100; i++ {
		switch {
		case i%15 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}

	// 問題2: 素数
	fmt.Println("=== 素数 ===")
	for n := 2; n <= 100; n++ {
		isPrime := true
		for i := 2; i*i <= n; i++ {
			if n%i == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			fmt.Printf("%d ", n)
		}
	}
	fmt.Println()

	// 問題3: 最大値
	fmt.Println("=== 最大値 ===")
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	max := nums[0]
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	fmt.Printf("最大値: %d\n", max)
}
