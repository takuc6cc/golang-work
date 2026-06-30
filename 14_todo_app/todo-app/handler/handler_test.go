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

// mockStore はテスト用のインメモリ Store
type mockStore struct {
	todos  []model.Todo
	err    error
	nextID int
}

func newMockStore(todos ...model.Todo) *mockStore {
	ms := &mockStore{nextID: 1}
	for i, t := range todos {
		t.ID = i + 1
		ms.todos = append(ms.todos, t)
		ms.nextID = i + 2
	}
	return ms
}

func (m *mockStore) GetAll() ([]model.Todo, error) {
	return m.todos, m.err
}

func (m *mockStore) GetByID(id int) (model.Todo, error) {
	if m.err != nil {
		return model.Todo{}, m.err
	}
	for _, t := range m.todos {
		if t.ID == id {
			return t, nil
		}
	}
	return model.Todo{}, store.ErrNotFound
}

func (m *mockStore) Create(todo model.Todo) (model.Todo, error) {
	if m.err != nil {
		return model.Todo{}, m.err
	}
	todo.ID = m.nextID
	m.nextID++
	m.todos = append(m.todos, todo)
	return todo, nil
}

func (m *mockStore) Update(todo model.Todo) (model.Todo, error) {
	if m.err != nil {
		return model.Todo{}, m.err
	}
	for i, t := range m.todos {
		if t.ID == todo.ID {
			m.todos[i] = todo
			return todo, nil
		}
	}
	return model.Todo{}, store.ErrNotFound
}

func (m *mockStore) Delete(id int) error {
	if m.err != nil {
		return m.err
	}
	for i, t := range m.todos {
		if t.ID == id {
			m.todos = append(m.todos[:i], m.todos[i+1:]...)
			return nil
		}
	}
	return store.ErrNotFound
}

func TestListTodos(t *testing.T) {
	ms := newMockStore(
		model.Todo{Title: "タスク1"},
		model.Todo{Title: "タスク2"},
	)
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
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"正常", `{"title": "新しいタスク"}`, http.StatusCreated},
		{"titleなし", `{"title": ""}`, http.StatusBadRequest},
		{"不正JSON", `{invalid}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := newMockStore()
			h := New(ms)

			req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h.CreateTodo(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetTodo(t *testing.T) {
	ms := newMockStore(model.Todo{Title: "タスク1"})
	h := New(ms)

	t.Run("存在するID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos/1", nil)
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()
		h.GetTodo(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("存在しないID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos/999", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()
		h.GetTodo(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d; want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestDeleteTodo(t *testing.T) {
	ms := newMockStore(model.Todo{Title: "タスク1"})
	h := New(ms)

	req := httptest.NewRequest("DELETE", "/todos/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	h.DeleteTodo(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d; want %d", rec.Code, http.StatusNoContent)
	}

	if len(ms.todos) != 0 {
		t.Errorf("todos 件数 = %d; want 0", len(ms.todos))
	}
}

func TestCompleteTodo(t *testing.T) {
	ms := newMockStore(model.Todo{Title: "タスク1", Done: false})
	h := New(ms)

	req := httptest.NewRequest("PATCH", "/todos/1/complete", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()
	h.CompleteTodo(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
	}

	var todo model.Todo
	json.NewDecoder(rec.Body).Decode(&todo)
	if !todo.Done {
		t.Error("Done = false; want true")
	}
}
