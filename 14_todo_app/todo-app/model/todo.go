package model

import "time"

// Todo は1件のタスクを表す
type Todo struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Done      bool       `json:"done"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// New は新しい Todo を作成する
func New(title string) Todo {
	return Todo{
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
}
