package main

import "fmt"

func filterEven(nums []int) []int {
	result := []int{}
	for _, n := range nums {
		if n%2 == 0 {
			result = append(result, n)
		}
	}
	return result
}

func wordCount(words []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	return counts
}

func main() {
	// 問題1
	result := filterEven([]int{1, 2, 3, 4, 5, 6, 7, 8})
	fmt.Println(result)

	// 問題2
	words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	counts := wordCount(words)
	fmt.Println(counts)

	// 問題3: 電話帳
	phonebook := make(map[string]string)
	phonebook["田中太郎"] = "090-1234-5678"
	phonebook["山田花子"] = "080-8765-4321"
	phonebook["鈴木一郎"] = "070-1111-2222"

	if num, ok := phonebook["田中太郎"]; ok {
		fmt.Printf("田中太郎: %s\n", num)
	}

	delete(phonebook, "山田花子")

	fmt.Println("=== 電話帳 ===")
	for name, num := range phonebook {
		fmt.Printf("%s: %s\n", name, num)
	}
}
