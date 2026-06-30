package main

import (
	"fmt"
	"strings"
	"time"

	"mypackage/greet"
)

func main() {
	// 問題1: greet パッケージを使って挨拶する
	fmt.Println(greet.Hello("Alice"))
	fmt.Println(greet.Goodbye("Bob"))

	// 問題2: strings パッケージを使った操作
	// "Hello, World! Hello, Go!" の "Hello" の出現回数を数える
	text := "Hello, World! Hello, Go!"
	_ = text
	// strings.Count を使ってください

	// "alice@example.com" からユーザー名を取り出す
	email := "alice@example.com"
	_ = email
	// strings.SplitN を使ってください

	// 問題3: time パッケージ
	// 現在時刻、1時間後、経過時間を表示する
	_ = time.Now()
}
