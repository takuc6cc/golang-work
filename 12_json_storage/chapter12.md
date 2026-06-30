# 第12章: JSON・ファイルストレージ

## encoding/json パッケージ

### 構造体 → JSON（Marshal / Encode）

```go
import "encoding/json"

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email,omitempty"` // 空なら省略
}

// Marshal: バイト列に変換
user := User{ID: 1, Name: "Alice", Email: "alice@example.com"}
data, err := json.Marshal(user)
// data = []byte(`{"id":1,"name":"Alice","email":"alice@example.com"}`)

// MarshalIndent: 整形出力
data, err = json.MarshalIndent(user, "", "  ")

// Encoder: io.Writer に直接書く（HTTP レスポンスに便利）
json.NewEncoder(os.Stdout).Encode(user)
```

Python との比較:
```python
import json

user = {"id": 1, "name": "Alice", "email": "alice@example.com"}
data = json.dumps(user)                    # Marshal
data = json.dumps(user, indent=2)          # MarshalIndent
print(json.dumps(user))                    # Encoder に近い
```

### JSON → 構造体（Unmarshal / Decode）

```go
jsonStr := `{"id":1,"name":"Alice","email":"alice@example.com"}`

var user User
err := json.Unmarshal([]byte(jsonStr), &user)
// user.Name → "Alice"

// Decoder: io.Reader から読む（HTTP リクエストに便利）
var user2 User
json.NewDecoder(r.Body).Decode(&user2)
```

Python との比較:
```python
data = '{"id": 1, "name": "Alice"}'
user = json.loads(data)      # Unmarshal
user = json.load(file)       # Decoder に近い
```

---

## 構造体タグ

```go
type Todo struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Done      bool      `json:"done"`
    CreatedAt time.Time `json:"created_at"`
    
    // omitempty: ゼロ値の場合に JSON から除外
    UpdatedAt *time.Time `json:"updated_at,omitempty"`
    
    // -: JSON に含めない
    InternalID string `json:"-"`
}
```

---

## ファイルへの読み書き

### os.ReadFile / os.WriteFile（推奨）

```go
import "os"

// 書き込み
data := []byte(`{"name": "Alice"}`)
err := os.WriteFile("data.json", data, 0644)

// 読み込み
data, err = os.ReadFile("data.json")
```

Python との比較:
```python
# Python
with open("data.json", "w") as f:
    json.dump(data, f)

with open("data.json", "r") as f:
    data = json.load(f)
```

```go
// Go
data, _ := json.Marshal(todos)
os.WriteFile("data.json", data, 0644)

data, _ = os.ReadFile("data.json")
json.Unmarshal(data, &todos)
```

### ファイルが存在しない場合の処理

```go
import (
    "errors"
    "os"
)

data, err := os.ReadFile("todos.json")
if err != nil {
    if errors.Is(err, os.ErrNotExist) {
        // ファイルが存在しない → 空で初期化
        return []Todo{}, nil
    }
    return nil, fmt.Errorf("ファイル読み込みエラー: %w", err)
}
```

---

## JSON ファイルストアの実装

TODOアプリで使うファイルストアを実装します。

```go
package store

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "sync"
    
    "myapp/model"
)

type FileStore struct {
    path   string
    mu     sync.Mutex
    nextID int
}

func NewFileStore(path string) *FileStore {
    return &FileStore{path: path, nextID: 1}
}

// ファイルからデータを読み込む（内部用）
func (fs *FileStore) load() ([]model.Todo, error) {
    data, err := os.ReadFile(fs.path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return []model.Todo{}, nil
        }
        return nil, fmt.Errorf("load: %w", err)
    }
    
    var todos []model.Todo
    if err := json.Unmarshal(data, &todos); err != nil {
        return nil, fmt.Errorf("load: %w", err)
    }
    return todos, nil
}

// ファイルにデータを保存する（内部用）
func (fs *FileStore) save(todos []model.Todo) error {
    data, err := json.MarshalIndent(todos, "", "  ")
    if err != nil {
        return fmt.Errorf("save: %w", err)
    }
    return os.WriteFile(fs.path, data, 0644)
}

func (fs *FileStore) GetAll() ([]model.Todo, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    return fs.load()
}

func (fs *FileStore) Create(todo model.Todo) (model.Todo, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    
    todos, err := fs.load()
    if err != nil {
        return model.Todo{}, err
    }
    
    todo.ID = fs.nextID
    fs.nextID++
    todos = append(todos, todo)
    
    return todo, fs.save(todos)
}
```

---

## time.Time の JSON シリアライズ

`time.Time` は自動的に RFC3339 形式（`"2024-01-15T10:30:00Z"`）で JSON 変換されます。

```go
type Todo struct {
    CreatedAt time.Time `json:"created_at"`
}

todo := Todo{CreatedAt: time.Now()}
data, _ := json.Marshal(todo)
// {"created_at":"2024-01-15T10:30:00.123456789+09:00"}
```

---

## A Tour of Go

本章は標準ライブラリ中心のため、公式ドキュメントを参照してください:
- https://pkg.go.dev/encoding/json
- https://pkg.go.dev/os

---

## 練習問題

### 問題1: JSON の読み書き

`[]Todo` を JSON ファイルに保存・読み込みする関数を書け。

```go
func saveTodos(path string, todos []Todo) error { ... }
func loadTodos(path string) ([]Todo, error) { ... }
```

### 問題2: 構造体タグの練習

以下の条件で構造体タグを設定した `Product` 構造体を定義せよ:
- `id` フィールドは `"id"` としてシリアライズ
- `price` フィールドは `"price"` としてシリアライズ
- `stock` フィールドは 0 の場合に JSON から除外（`omitempty`）
- `internal_code` フィールドは JSON に含めない

### 問題3: goroutine セーフな FileStore

`sync.Mutex` で保護した `FileStore` を実装せよ。
`GetAll`・`Create`・`Delete` メソッドを実装すること。

---

## 解答

`exercises/json_crud_solution/main.go`:

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Todo モデル
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// 問題1: JSON の読み書き
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

// 問題2: 構造体タグ
type Product struct {
	ID           int     `json:"id"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock,omitempty"`
	InternalCode string  `json:"-"`
}

// 問題3: FileStore
type FileStore struct {
	path   string
	mu     sync.Mutex
	nextID int
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path, nextID: 1}
}

func (fs *FileStore) load() ([]Todo, error) {
	return loadTodos(fs.path)
}

func (fs *FileStore) save(todos []Todo) error {
	return saveTodos(fs.path, todos)
}

func (fs *FileStore) GetAll() ([]Todo, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.load()
}

func (fs *FileStore) Create(title string) (Todo, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	todos, err := fs.load()
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
	return todo, fs.save(todos)
}

func (fs *FileStore) Delete(id int) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	todos, err := fs.load()
	if err != nil {
		return err
	}

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return fs.save(todos)
		}
	}
	return fmt.Errorf("ID %d: not found", id)
}

func main() {
	const path = "todos.json"
	defer os.Remove(path) // テスト後に削除

	// 問題1のテスト
	todos := []Todo{
		{ID: 1, Title: "Goを学ぶ", Done: false, CreatedAt: time.Now()},
		{ID: 2, Title: "TODOアプリを作る", Done: false, CreatedAt: time.Now()},
	}
	if err := saveTodos(path, todos); err != nil {
		fmt.Println("エラー:", err)
		return
	}

	loaded, err := loadTodos(path)
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}
	fmt.Printf("読み込み: %d件\n", len(loaded))

	// 問題2のテスト
	p1 := Product{ID: 1, Price: 1500.0, Stock: 0, InternalCode: "SKU-001"}
	p2 := Product{ID: 2, Price: 3000.0, Stock: 10, InternalCode: "SKU-002"}
	data, _ := json.MarshalIndent([]Product{p1, p2}, "", "  ")
	fmt.Println(string(data))
	// InternalCode は含まれない、Stock=0 は含まれない

	// 問題3のテスト
	store := NewFileStore("store_test.json")
	defer os.Remove("store_test.json")

	t1, _ := store.Create("タスク1")
	t2, _ := store.Create("タスク2")
	fmt.Printf("作成: %+v\n", t1)
	fmt.Printf("作成: %+v\n", t2)

	all, _ := store.GetAll()
	fmt.Printf("全件: %d件\n", len(all))

	store.Delete(t1.ID)
	all, _ = store.GetAll()
	fmt.Printf("削除後: %d件\n", len(all))
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| 構造体→JSON | `json.dumps(dict)` | `json.Marshal(struct)` |
| JSON→構造体 | `json.loads(str)` | `json.Unmarshal(data, &struct)` |
| ファイル書き込み | `open + json.dump` | `os.WriteFile + json.Marshal` |
| ファイル読み込み | `open + json.load` | `os.ReadFile + json.Unmarshal` |
| フィールド名変更 | `@dataclass` + カスタム or dict | 構造体タグ `json:"name"` |

次の章ではテストを学びます。→ [第13章](../13_testing/chapter13.md)
