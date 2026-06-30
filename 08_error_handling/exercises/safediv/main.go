package main

import "fmt"

// 問題2: カスタムエラー型を定義してください
// type ValidationError struct { ... }
// func (e *ValidationError) Error() string { ... }
// func validate(title, description string) error { ... }

// 問題3: Sentinel エラーを定義してください
// var ErrNotFound = ...
// var ErrAlreadyDone = ...

// type Todo struct { ... }
// func CompleteTodo(todos []Todo, id int) error { ... }

// 問題1: ファイルの最初の行を返す関数を実装してください
// func readFirstLine(path string) (string, error) { ... }

func main() {
	// 問題1のテスト（存在しないファイルでエラーを確認）
	// _, err := readFirstLine("nonexistent.txt")
	// if err != nil { fmt.Println("エラー:", err) }

	// 問題2のテスト
	// err = validate("", "説明")
	// ...

	// 問題3のテスト
	// todos := []Todo{ ... }
	// err = CompleteTodo(todos, 1)
	// ...

	fmt.Println("練習問題を実装してください")
}
