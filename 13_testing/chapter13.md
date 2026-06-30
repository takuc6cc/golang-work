# 第13章: テスト

## testing パッケージの基本

Go はテストが標準機能として組み込まれています。
テストファイルは `_test.go` で終わる名前にし、テスト関数は `Test` で始めます。

```go
// calc.go
package calc

func Add(a, b int) int {
    return a + b
}
```

```go
// calc_test.go
package calc

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}
```

```bash
go test ./...           # テスト実行
go test -v ./...        # 詳細出力
go test -cover ./...    # カバレッジ表示
go test -run TestAdd    # 特定のテストのみ
```

Python との比較:
```python
# test_calc.py (pytest)
from calc import add

def test_add():
    assert add(2, 3) == 5
```

```bash
pytest test_calc.py
pytest -v
pytest --cov=calc
```

---

## テーブル駆動テスト

Go の慣用的なテストパターンです。複数のケースをまとめてテストします。

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"正の数", 2, 3, 5},
        {"負の数", -1, -2, -3},
        {"ゼロ", 0, 0, 0},
        {"正と負", 5, -3, 2},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

Python との比較:
```python
import pytest

@pytest.mark.parametrize("a, b, expected", [
    (2, 3, 5),
    (-1, -2, -3),
    (0, 0, 0),
    (5, -3, 2),
])
def test_add(a, b, expected):
    assert add(a, b) == expected
```

---

## t のメソッド

```go
t.Error("失敗メッセージ")          // テスト失敗（実行は続く）
t.Errorf("失敗: got %d", got)     // フォーマット付き失敗
t.Fatal("致命的なエラー")          // テスト失敗（即座に停止）
t.Fatalf("致命的: %v", err)
t.Log("ログメッセージ")            // -v フラグ時のみ表示
t.Helper()                         // ヘルパー関数内でコールスタックを正しく表示
t.Skip("スキップ")                 // テストをスキップ
```

---

## インターフェースを使ったモック

インターフェースを使うと、テスト用のモック実装を作れます。
TODOアプリのハンドラーをテストする場合の例:

```go
// store.go
type Store interface {
    GetAll() ([]Todo, error)
    Create(todo Todo) (Todo, error)
}

// テスト用モック
type MockStore struct {
    todos []Todo
    err   error
}

func (m *MockStore) GetAll() ([]Todo, error) {
    return m.todos, m.err
}

func (m *MockStore) Create(todo Todo) (Todo, error) {
    todo.ID = len(m.todos) + 1
    m.todos = append(m.todos, todo)
    return todo, m.err
}

// テスト
func TestListTodos(t *testing.T) {
    mockStore := &MockStore{
        todos: []Todo{
            {ID: 1, Title: "タスク1"},
            {ID: 2, Title: "タスク2"},
        },
    }
    
    handler := NewHandler(mockStore)
    req := httptest.NewRequest("GET", "/todos", nil)
    rec := httptest.NewRecorder()
    
    handler.ListTodos(rec, req)
    
    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; want %d", rec.Code, http.StatusOK)
    }
}
```

---

## httptest パッケージ

HTTP ハンドラーのテストには `net/http/httptest` を使います。

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHelloHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/hello", nil)
    rec := httptest.NewRecorder()

    helloHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; want 200", rec.Code)
    }
    
    body := rec.Body.String()
    if !strings.Contains(body, "Hello") {
        t.Errorf("body = %q; want to contain 'Hello'", body)
    }
}
```

---

## カバレッジの確認

```bash
go test -cover ./...
# ok  myapp/handler  coverage: 85.7% of statements

# HTML レポートを生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## A Tour of Go

公式ドキュメント: https://pkg.go.dev/testing

---

## 練習問題

`exercises/calc_test/` ディレクトリに以下を作れ:

### 問題1: テーブル駆動テスト

`calc.go` に `Add`・`Subtract`・`Multiply`・`Divide` 関数を実装し、
`calc_test.go` でテーブル駆動テストを書け。

`Divide` は除数が 0 のときエラーを返すこと。

### 問題2: モックを使ったテスト

`Store` インターフェースのモックを実装し、
`FindAll` と `Create` を使う関数のユニットテストを書け。

### 問題3: カバレッジ80%以上

```bash
go test -cover ./exercises/calc_test/...
```

80% 以上を達成すること。

---

## 解答

`exercises/calc_test_solution/calc.go`:

```go
package calc

import "errors"

func Add(a, b int) int {
	return a + b
}

func Subtract(a, b int) int {
	return a - b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("ゼロ除算はできません")
	}
	return a / b, nil
}
```

`exercises/calc_test_solution/calc_test.go`:

```go
package calc

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"正の数", 2, 3, 5},
		{"負の数", -1, -2, -3},
		{"ゼロ", 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{5, 3, 2},
		{0, 5, -5},
		{-3, -2, -1},
	}
	for _, tt := range tests {
		if got := Subtract(tt.a, tt.b); got != tt.want {
			t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{3, 4, 12},
		{-2, 5, -10},
		{0, 100, 0},
	}
	for _, tt := range tests {
		if got := Multiply(tt.a, tt.b); got != tt.want {
			t.Errorf("Multiply(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestDivide(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		tests := []struct{ a, b, want int }{
			{10, 2, 5},
			{9, 3, 3},
			{7, 2, 3}, // 切り捨て
		}
		for _, tt := range tests {
			got, err := Divide(tt.a, tt.b)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if got != tt.want {
				t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		}
	})

	t.Run("ゼロ除算", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Error("エラーが返るはずですが、nil が返りました")
		}
	})
}
```

---

## まとめ

| 概念 | Python (pytest) | Go (testing) |
|------|-----------------|--------------|
| テストファイル | `test_foo.py` | `foo_test.go` |
| テスト関数 | `def test_foo():` | `func TestFoo(t *testing.T)` |
| アサーション | `assert x == y` | `if x != y { t.Errorf(...) }` |
| パラメータ化 | `@pytest.mark.parametrize` | テーブル駆動テスト |
| 実行コマンド | `pytest` | `go test ./...` |
| カバレッジ | `pytest --cov` | `go test -cover` |

次の章では TODO アプリを完成させます。→ [第14章](../14_todo_app/chapter14.md)
