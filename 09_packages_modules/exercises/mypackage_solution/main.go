package main

import (
	"fmt"
	"strings"
	"time"

	"mypackage/greet"
)

func main() {
	// 問題1: greet パッケージ
	fmt.Println(greet.Hello("Alice"))
	fmt.Println(greet.Goodbye("Bob"))

	// 問題2: strings パッケージ
	text := "Hello, World! Hello, Go!"
	count := strings.Count(text, "Hello")
	fmt.Printf("\"Hello\" の出現回数: %d\n", count)

	parts := strings.SplitN("alice@example.com", "@", 2)
	fmt.Printf("ユーザー名: %s\n", parts[0])

	// 問題3: time パッケージ
	now := time.Now()
	fmt.Println("現在時刻:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("1時間後:", now.Add(time.Hour).Format("2006-01-02 15:04:05"))
	start := time.Now()
	elapsed := time.Since(start)
	fmt.Printf("経過時間: %v\n", elapsed)
}
