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

// Handler は HTTP ハンドラーを保持する
type Handler struct {
	store store.Store
}

// New は新しい Handler を作成する
func New(s store.Store) *Handler {
	return &Handler{store: s}
}

// RegisterRoutes は ServeMux にルートを登録する
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /todos", h.ListTodos)
	mux.HandleFunc("POST /todos", h.CreateTodo)
	mux.HandleFunc("GET /todos/{id}", h.GetTodo)
	mux.HandleFunc("PUT /todos/{id}", h.UpdateTodo)
	mux.HandleFunc("DELETE /todos/{id}", h.DeleteTodo)
	mux.HandleFunc("PATCH /todos/{id}/complete", h.CompleteTodo)
}

// writeJSON は JSON レスポンスを書く
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError はエラーレスポンスを書く
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// parseID はパスパラメータから ID を取り出す
func parseID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id: %s", idStr)
	}
	return id, nil
}

// ListTodos は GET /todos のハンドラー
func (h *Handler) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, todos)
}

// CreateTodo は POST /todos のハンドラー
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

// GetTodo は GET /todos/{id} のハンドラー
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

// UpdateTodo は PUT /todos/{id} のハンドラー
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
	if strings.TrimSpace(req.Title) == "" {
		writeError(w, http.StatusBadRequest, "title は必須です")
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

// DeleteTodo は DELETE /todos/{id} のハンドラー
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

// CompleteTodo は PATCH /todos/{id}/complete のハンドラー
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
