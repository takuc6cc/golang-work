# 第14章: TODOアプリ実装

いよいよ最終章です。第1〜13章で学んだすべての知識を使って、
REST API を持つ TODO アプリを完成させます。

## アプリケーションの概要

- データはJSONファイルに永続化（`todos.json`）
- フレームワーク不使用（標準ライブラリのみ）
- Go 1.22 の `http.ServeMux`（パスパラメータ対応）を使用

### エンドポイント一覧

| メソッド | パス | 機能 |
|---------|------|------|
| `GET` | `/todos` | 一覧取得 |
| `POST` | `/todos` | 新規作成 |
| `GET` | `/todos/{id}` | 1件取得 |
| `PUT` | `/todos/{id}` | 更新 |
| `DELETE` | `/todos/{id}` | 削除 |
| `PATCH` | `/todos/{id}/complete` | 完了にする |

---

## ディレクトリ構造

```
14_todo_app/todo-app/
├── go.mod
├── main.go              ← エントリポイント・サーバー起動
├── model/
│   └── todo.go          ← Todo 構造体の定義
├── store/
│   ├── store.go         ← Store インターフェース
│   └── file_store.go    ← JSON ファイル実装
└── handler/
    ├── handler.go        ← HTTP ハンドラー
    └── handler_test.go  ← ハンドラーのテスト
```

---

## パート1: モデル定義

`model/todo.go`:

```go
package model

import "time"

type Todo struct {
    ID        int        `json:"id"`
    Title     string     `json:"title"`
    Done      bool       `json:"done"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func New(title string) Todo {
    return Todo{
        Title:     title,
        Done:      false,
        CreatedAt: time.Now(),
    }
}
```

---

## パート2: Store インターフェース

`store/store.go`:

```go
package store

import (
    "errors"
    "todo-app/model"
)

var ErrNotFound = errors.New("not found")

type Store interface {
    GetAll() ([]model.Todo, error)
    GetByID(id int) (model.Todo, error)
    Create(todo model.Todo) (model.Todo, error)
    Update(todo model.Todo) (model.Todo, error)
    Delete(id int) error
}
```

`store/file_store.go`:

```go
package store

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "sync"
    "time"
    "todo-app/model"
)

type FileStore struct {
    path   string
    mu     sync.Mutex
    nextID int
}

func NewFileStore(path string) *FileStore {
    return &FileStore{path: path, nextID: 1}
}

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
        return nil, fmt.Errorf("load unmarshal: %w", err)
    }
    // nextID を更新
    for _, t := range todos {
        if t.ID >= fs.nextID {
            fs.nextID = t.ID + 1
        }
    }
    return todos, nil
}

func (fs *FileStore) save(todos []model.Todo) error {
    data, err := json.MarshalIndent(todos, "", "  ")
    if err != nil {
        return fmt.Errorf("save marshal: %w", err)
    }
    return os.WriteFile(fs.path, data, 0644)
}

func (fs *FileStore) GetAll() ([]model.Todo, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    return fs.load()
}

func (fs *FileStore) GetByID(id int) (model.Todo, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    todos, err := fs.load()
    if err != nil {
        return model.Todo{}, err
    }
    for _, t := range todos {
        if t.ID == id {
            return t, nil
        }
    }
    return model.Todo{}, fmt.Errorf("ID %d: %w", id, ErrNotFound)
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

func (fs *FileStore) Update(todo model.Todo) (model.Todo, error) {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    todos, err := fs.load()
    if err != nil {
        return model.Todo{}, err
    }
    for i, t := range todos {
        if t.ID == todo.ID {
            now := time.Now()
            todo.CreatedAt = t.CreatedAt
            todo.UpdatedAt = &now
            todos[i] = todo
            return todo, fs.save(todos)
        }
    }
    return model.Todo{}, fmt.Errorf("ID %d: %w", todo.ID, ErrNotFound)
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
    return fmt.Errorf("ID %d: %w", id, ErrNotFound)
}
```

---

## パート3: HTTP ハンドラー

`handler/handler.go`:

```go
package handler

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"
    "todo-app/model"
    "todo-app/store"
)

type Handler struct {
    store store.Store
}

func New(s store.Store) *Handler {
    return &Handler{store: s}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("GET /todos", h.ListTodos)
    mux.HandleFunc("POST /todos", h.CreateTodo)
    mux.HandleFunc("GET /todos/{id}", h.GetTodo)
    mux.HandleFunc("PUT /todos/{id}", h.UpdateTodo)
    mux.HandleFunc("DELETE /todos/{id}", h.DeleteTodo)
    mux.HandleFunc("PATCH /todos/{id}/complete", h.CompleteTodo)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

func parseID(r *http.Request) (int, error) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        return 0, fmt.Errorf("invalid id: %s", idStr)
    }
    return id, nil
}

func (h *Handler) ListTodos(w http.ResponseWriter, r *http.Request) {
    todos, err := h.store.GetAll()
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, todos)
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Title string `json:"title"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "リクエストが不正です")
        return
    }
    if strings.TrimSpace(req.Title) == "" {
        writeError(w, http.StatusBadRequest, "title は必須です")
        return
    }

    todo, err := h.store.Create(model.New(req.Title))
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    writeJSON(w, http.StatusCreated, todo)
}

func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    todo, err := h.store.GetByID(id)
    if err != nil {
        if errors.Is(err, store.ErrNotFound) {
            writeError(w, http.StatusNotFound, "TODO が見つかりません")
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, todo)
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    var req struct {
        Title string `json:"title"`
        Done  bool   `json:"done"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "リクエストが不正です")
        return
    }

    todo := model.Todo{ID: id, Title: req.Title, Done: req.Done}
    updated, err := h.store.Update(todo)
    if err != nil {
        if errors.Is(err, store.ErrNotFound) {
            writeError(w, http.StatusNotFound, "TODO が見つかりません")
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    if err := h.store.Delete(id); err != nil {
        if errors.Is(err, store.ErrNotFound) {
            writeError(w, http.StatusNotFound, "TODO が見つかりません")
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
    id, err := parseID(r)
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }
    todo, err := h.store.GetByID(id)
    if err != nil {
        if errors.Is(err, store.ErrNotFound) {
            writeError(w, http.StatusNotFound, "TODO が見つかりません")
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    todo.Done = true
    now := time.Now()
    todo.UpdatedAt = &now
    updated, err := h.store.Update(todo)
    if err != nil {
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }
    writeJSON(w, http.StatusOK, updated)
}
```

---

## パート4: main.go

`main.go`:

```go
package main

import (
    "log"
    "net/http"
    "todo-app/handler"
    "todo-app/store"
)

func main() {
    s := store.NewFileStore("todos.json")
    h := handler.New(s)

    mux := http.NewServeMux()
    h.RegisterRoutes(mux)

    log.Println("TODO アプリを起動: http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

---

## パート5: テスト

`handler/handler_test.go`:

```go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "todo-app/model"
    "todo-app/store"
)

// モック Store
type mockStore struct {
    todos []model.Todo
    err   error
}

func (m *mockStore) GetAll() ([]model.Todo, error) {
    return m.todos, m.err
}

func (m *mockStore) GetByID(id int) (model.Todo, error) {
    for _, t := range m.todos {
        if t.ID == id {
            return t, nil
        }
    }
    return model.Todo{}, store.ErrNotFound
}

func (m *mockStore) Create(todo model.Todo) (model.Todo, error) {
    todo.ID = len(m.todos) + 1
    m.todos = append(m.todos, todo)
    return todo, m.err
}

func (m *mockStore) Update(todo model.Todo) (model.Todo, error) {
    for i, t := range m.todos {
        if t.ID == todo.ID {
            m.todos[i] = todo
            return todo, nil
        }
    }
    return model.Todo{}, store.ErrNotFound
}

func (m *mockStore) Delete(id int) error {
    for i, t := range m.todos {
        if t.ID == id {
            m.todos = append(m.todos[:i], m.todos[i+1:]...)
            return nil
        }
    }
    return store.ErrNotFound
}

func TestListTodos(t *testing.T) {
    ms := &mockStore{
        todos: []model.Todo{
            {ID: 1, Title: "タスク1"},
            {ID: 2, Title: "タスク2"},
        },
    }
    h := New(ms)

    req := httptest.NewRequest("GET", "/todos", nil)
    rec := httptest.NewRecorder()
    h.ListTodos(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
    }

    var todos []model.Todo
    json.NewDecoder(rec.Body).Decode(&todos)
    if len(todos) != 2 {
        t.Errorf("len(todos) = %d; want 2", len(todos))
    }
}

func TestCreateTodo(t *testing.T) {
    ms := &mockStore{}
    h := New(ms)

    body := bytes.NewBufferString(`{"title": "新しいタスク"}`)
    req := httptest.NewRequest("POST", "/todos", body)
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()
    h.CreateTodo(rec, req)

    if rec.Code != http.StatusCreated {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusCreated)
    }

    var todo model.Todo
    json.NewDecoder(rec.Body).Decode(&todo)
    if todo.Title != "新しいタスク" {
        t.Errorf("title = %q; want '新しいタスク'", todo.Title)
    }
}
```

---

## パート6: 動作確認

```bash
# プロジェクトに移動
cd 14_todo_app/todo-app

# 起動
go run main.go
# 2024/01/15 10:30:00 TODO アプリを起動: http://localhost:8080
```

別ターミナルで:

```bash
# 一覧取得（空）
curl http://localhost:8080/todos
# []

# 作成
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Goを学ぶ"}'
# {"id":1,"title":"Goを学ぶ","done":false,"created_at":"..."}

# もう1件作成
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "TODOアプリを作る"}'

# 一覧取得
curl http://localhost:8080/todos

# 1件取得
curl http://localhost:8080/todos/1

# 完了にする
curl -X PATCH http://localhost:8080/todos/1/complete

# 更新
curl -X PUT http://localhost:8080/todos/2 \
  -H "Content-Type: application/json" \
  -d '{"title": "TODOアプリを完成させる", "done": false}'

# 削除
curl -X DELETE http://localhost:8080/todos/1
# （204 No Content が返る）

# テスト実行
go test ./...
```

---

## パート7: 発展課題（オプション）

TODOアプリが完成したら、以下に挑戦してみましょう。

### 発展1: HTML UI を追加

`static/index.html` に JavaScript で書いたフロントエンドを追加し、
ブラウザから操作できるようにする。

```go
// main.go に追加
mux.Handle("GET /", http.FileServer(http.Dir("static")))
```

### 発展2: ルーターライブラリを使う

```bash
go get github.com/go-chi/chi/v5
```

`chi` を使うとパスパラメータや middleware がより書きやすくなります。

### 発展3: データベースを使う

SQLite を使う場合:

```bash
go get modernc.org/sqlite
```

`FileStore` の代わりに `SQLiteStore` を実装し、`Store` インターフェースに沿って差し替えます。
ハンドラーのコードはまったく変更不要です（インターフェースの力！）

### 発展4: Docker コンテナ化

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o todo-app .

FROM alpine:3.19
COPY --from=builder /app/todo-app /usr/local/bin/
CMD ["todo-app"]
```

---

## おめでとうございます！

第1章から第14章まで完走し、Go で REST API を持つ TODO アプリを完成させました。

**習得したスキル:**
- Go の基本文法（変数・型・制御フロー・関数）
- データ構造（スライス・マップ・構造体）
- インターフェース設計
- エラーハンドリング
- パッケージとモジュール管理
- goroutine による並行処理
- `net/http` での HTTP サーバー構築
- JSON のシリアライズ・デシリアライズ
- ファイルへの永続化
- テスト（テーブル駆動テスト・モック）

**次のステップ:**
- [A Tour of Go](https://go.dev/tour/) で残りのセクションを学ぶ
- [Effective Go](https://go.dev/doc/effective_go) でイディオムを学ぶ
- [Go by Example](https://gobyexample.com/) で実践的なパターンを学ぶ
