package main

import (
	"fmt"
	"sync"
)

// 問題3: goroutine セーフカウンターを実装してください
// type SafeCounter struct { ... }
// func (c *SafeCounter) Increment() { ... }
// func (c *SafeCounter) Value() int { ... }

func main() {
	// 問題1: 3つの goroutine を起動して WaitGroup で待つ
	var wg sync.WaitGroup
	_ = wg
	// wg.Add(1)
	// go func() { defer wg.Done(); fmt.Println("...") }()
	// wg.Wait()

	// 問題2: channel を使ったパイプライン
	// generate goroutine: 1〜10 を channel に送信
	// square goroutine: 受け取って2乗を別の channel に送信
	// main: 結果を受け取って出力

	// 問題3: 100 goroutine × 100回インクリメント = 10000 を確認
	// counter := &SafeCounter{}
	// ...

	fmt.Println("練習問題を実装してください")
}
