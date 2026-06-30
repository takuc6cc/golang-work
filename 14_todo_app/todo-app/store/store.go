package store

import (
	"errors"
	"todo-app/model"
)

// ErrNotFound は指定した ID の Todo が存在しない場合のエラー
var ErrNotFound = errors.New("not found")

// Store は Todo の永続化を抽象化するインターフェース
type Store interface {
	GetAll() ([]model.Todo, error)
	GetByID(id int) (model.Todo, error)
	Create(todo model.Todo) (model.Todo, error)
	Update(todo model.Todo) (model.Todo, error)
	Delete(id int) error
}
