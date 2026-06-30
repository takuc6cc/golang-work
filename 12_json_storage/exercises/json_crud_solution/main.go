package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

func saveTodos(path string, todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func loadTodos(path string) ([]Todo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Todo{}, nil
		}
		return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
	}
	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, fmt.Errorf("JSON解析エラー: %w", err)
	}
	return todos, nil
}

type Product struct {
	ID           int     `json:"id"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock,omitempty"`
	InternalCode string  `json:"-"`
}

type FileStore struct {
	path   string
	mu     sync.Mutex
	nextID int
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path, nextID: 1}
}

func (fs *FileStore) GetAll() ([]Todo, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return loadTodos(fs.path)
}

func (fs *FileStore) Create(title string) (Todo, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	todos, err := loadTodos(fs.path)
	if err != nil {
		return Todo{}, err
	}

	todo := Todo{
		ID:        fs.nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	fs.nextID++
	todos = append(todos, todo)
	return todo, saveTodos(fs.path, todos)
}

func (fs *FileStore) Delete(id int) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	todos, err := loadTodos(fs.path)
	if err != nil {
		return err
	}

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return saveTodos(fs.path, todos)
		}
	}
	return fmt.Errorf("ID %d: not found", id)
}

func main() {
	const path = "todos.json"
	defer os.Remove(path)

	todos := []Todo{
		{ID: 1, Title: "Goを学ぶ", Done: false, CreatedAt: time.Now()},
		{ID: 2, Title: "TODOアプリを作る", Done: false, CreatedAt: time.Now()},
	}
	saveTodos(path, todos)

	loaded, _ := loadTodos(path)
	fmt.Printf("読み込み: %d件\n", len(loaded))

	p1 := Product{ID: 1, Price: 1500.0, Stock: 0, InternalCode: "SKU-001"}
	p2 := Product{ID: 2, Price: 3000.0, Stock: 10, InternalCode: "SKU-002"}
	data, _ := json.MarshalIndent([]Product{p1, p2}, "", "  ")
	fmt.Println(string(data))

	store := NewFileStore("store_test.json")
	defer os.Remove("store_test.json")

	t1, _ := store.Create("タスク1")
	t2, _ := store.Create("タスク2")
	fmt.Printf("作成: ID=%d %s\n", t1.ID, t1.Title)
	fmt.Printf("作成: ID=%d %s\n", t2.ID, t2.Title)

	all, _ := store.GetAll()
	fmt.Printf("全件: %d件\n", len(all))

	store.Delete(t1.ID)
	all, _ = store.GetAll()
	fmt.Printf("削除後: %d件\n", len(all))
}
