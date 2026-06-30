package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validate(title, description string) error {
	if strings.TrimSpace(title) == "" {
		return &ValidationError{Field: "title", Message: "空にできません"}
	}
	if len(title) > 100 {
		return &ValidationError{Field: "title", Message: "100文字以内にしてください"}
	}
	if len(description) > 500 {
		return &ValidationError{Field: "description", Message: "500文字以内にしてください"}
	}
	return nil
}

var (
	ErrNotFound    = errors.New("not found")
	ErrAlreadyDone = errors.New("already done")
)

type Todo struct {
	ID    int
	Title string
	Done  bool
}

func CompleteTodo(todos []Todo, id int) error {
	for i, t := range todos {
		if t.ID == id {
			if t.Done {
				return fmt.Errorf("ID %d: %w", id, ErrAlreadyDone)
			}
			todos[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("ID %d: %w", id, ErrNotFound)
}

func readFirstLine(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("ファイルを開けません: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", nil
}

func main() {
	_, err := readFirstLine("nonexistent.txt")
	if err != nil {
		fmt.Println("エラー:", err)
	}

	err = validate("", "説明")
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("バリデーションエラー - %s: %s\n", ve.Field, ve.Message)
	}

	fmt.Println("バリデーション結果:", validate("有効なタイトル", "説明"))

	todos := []Todo{
		{ID: 1, Title: "Goを学ぶ", Done: false},
		{ID: 2, Title: "テストを書く", Done: true},
	}

	fmt.Println("Complete 1:", CompleteTodo(todos, 1))

	err = CompleteTodo(todos, 2)
	if errors.Is(err, ErrAlreadyDone) {
		fmt.Println("すでに完了済みです")
	}

	err = CompleteTodo(todos, 99)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("見つかりません")
	}
}
