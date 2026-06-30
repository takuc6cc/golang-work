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

// FileStore は JSON ファイルに Todo を永続化する Store 実装
type FileStore struct {
	path   string
	mu     sync.Mutex
	nextID int
}

// NewFileStore は新しい FileStore を作成する
func NewFileStore(path string) *FileStore {
	return &FileStore{path: path, nextID: 1}
}

// load はファイルから Todo リストを読み込む（内部用）
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
	// nextID を最大 ID + 1 に更新
	for _, t := range todos {
		if t.ID >= fs.nextID {
			fs.nextID = t.ID + 1
		}
	}
	return todos, nil
}

// save は Todo リストをファイルに保存する（内部用）
func (fs *FileStore) save(todos []model.Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return fmt.Errorf("save marshal: %w", err)
	}
	return os.WriteFile(fs.path, data, 0644)
}

// GetAll は全件取得する
func (fs *FileStore) GetAll() ([]model.Todo, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.load()
}

// GetByID は指定した ID の Todo を取得する
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

// Create は新しい Todo を作成して保存する
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

// Update は指定した ID の Todo を更新する
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

// Delete は指定した ID の Todo を削除する
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
