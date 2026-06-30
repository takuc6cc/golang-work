package main

import (
	"fmt"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// 問題1: JSON ファイルへの保存・読み込み関数を実装してください
// func saveTodos(path string, todos []Todo) error { ... }
// func loadTodos(path string) ([]Todo, error) { ... }

// 問題2: 構造体タグを設定した Product 構造体を定義してください
// type Product struct { ... }

// 問題3: sync.Mutex を使った goroutine セーフな FileStore を実装してください
// type FileStore struct { ... }
// func (fs *FileStore) GetAll() ([]Todo, error) { ... }
// func (fs *FileStore) Create(title string) (Todo, error) { ... }
// func (fs *FileStore) Delete(id int) error { ... }

func main() {
	// 問題1のテスト
	todos := []Todo{
		{ID: 1, Title: "Goを学ぶ", Done: false, CreatedAt: time.Now()},
	}
	_ = todos
	// saveTodos("todos.json", todos)
	// loaded, _ := loadTodos("todos.json")
	// fmt.Printf("読み込み: %d件\n", len(loaded))

	fmt.Println("練習問題を実装してください")
}
