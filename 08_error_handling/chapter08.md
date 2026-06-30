# 第8章: エラーハンドリング

## Go のエラー哲学

Python や PHP は例外（Exception）でエラーを扱います。
Go は**例外を持たず、エラーを戻り値として返します**。

```python
# Python: 例外を投げてキャッチする
try:
    result = divide(10, 0)
except ZeroDivisionError as e:
    print(f"エラー: {e}")
```

```go
// Go: エラーを戻り値で返す
result, err := divide(10, 0)
if err != nil {
    fmt.Println("エラー:", err)
    return
}
fmt.Println(result)
```

**Go 式の利点:**
- エラーが発生しうる箇所が明確（呼び出しコードを読めばわかる）
- エラーを無視すると `_ =` と明示しなければならない（うっかり無視しにくい）
- 制御フローが単純（例外の伝播を追わなくていい）

---

## error インターフェース

`error` は Go 組み込みのインターフェースです。

```go
type error interface {
    Error() string
}
```

`Error() string` メソッドを持つ型はすべて `error` として使えます。

---

## エラーの作り方

### errors.New: シンプルなエラー

```go
import "errors"

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("ゼロ除算はできません")
    }
    return a / b, nil
}
```

### fmt.Errorf: フォーマット付きエラー

```go
import "fmt"

func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("無効なID: %d", id)
    }
    // ...
}
```

---

## エラーのラッピング（%w）

エラーに文脈情報を付加しながら、元のエラーを保持します。

```go
func readConfig(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        // %w でラップすると元のエラーを取り出せる
        return fmt.Errorf("設定ファイルの読み込み失敗 (%s): %w", path, err)
    }
    _ = data
    return nil
}
```

### errors.Is: ラップされたエラーの比較

```go
var ErrNotFound = errors.New("見つかりません")

func findTodo(id int) error {
    return fmt.Errorf("findTodo: %w", ErrNotFound)
}

err := findTodo(99)
if errors.Is(err, ErrNotFound) {
    fmt.Println("TODO が見つかりません")
}
```

Python との比較:
```python
class NotFoundError(Exception):
    pass

try:
    raise NotFoundError("見つかりません")
except NotFoundError as e:
    print(e)
```

### errors.As: カスタムエラー型の取り出し

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validate(title string) error {
    if title == "" {
        return &ValidationError{Field: "title", Message: "空にできません"}
    }
    return nil
}

err := validate("")
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Printf("バリデーションエラー - フィールド: %s, メッセージ: %s\n", ve.Field, ve.Message)
}
```

---

## Sentinel エラー

パッケージレベルで定義したエラー変数を「番兵エラー（sentinel error）」と呼びます。

```go
var (
    ErrNotFound   = errors.New("not found")
    ErrPermission = errors.New("permission denied")
)

func GetTodo(id int) (*Todo, error) {
    if id > len(todos) {
        return nil, ErrNotFound
    }
    return &todos[id-1], nil
}

// 呼び出し側
todo, err := GetTodo(99)
if errors.Is(err, ErrNotFound) {
    // 404 Not Found を返す
}
```

---

## panic と recover

`panic` は回復不能なエラーに使います。通常のエラーハンドリングには**使いません**。

```go
func mustPositive(n int) int {
    if n <= 0 {
        panic(fmt.Sprintf("正の数が必要です: %d", n))
    }
    return n
}
```

`recover` は `panic` を捕捉します。ミドルウェアなどで使われます。

```go
func safeExecute(f func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("パニックを回復: %v", r)
        }
    }()
    f()
    return nil
}
```

Python との比較:
```python
# panic に近いのは RuntimeError を raise すること
raise RuntimeError("回復不能なエラー")

# recover に近いのは except
try:
    risky_function()
except Exception as e:
    print(f"回復: {e}")
```

**ルール: panic は「プログラムのバグ」に使い、「ユーザーの入力ミス」には使わない**

---

## エラーハンドリングのパターン

### パターン1: 早期リターン

```go
func process() error {
    a, err := stepA()
    if err != nil {
        return fmt.Errorf("stepA: %w", err)
    }

    b, err := stepB(a)
    if err != nil {
        return fmt.Errorf("stepB: %w", err)
    }

    return stepC(b)
}
```

### パターン2: エラーを無視（明示的に）

```go
// エラーを本当に無視してよい場合のみ（推奨しない）
_ = os.Remove("temp.txt")
```

---

## A Tour of Go

本章は標準ライブラリ中心のため、Tour には対応セクションが少ないです。
関連: https://go.dev/tour/methods/19 （エラー）

---

## 練習問題

### 問題1: エラーを返す関数

ファイルを開いて最初の行を返す関数 `readFirstLine(path string) (string, error)` を実装せよ。
`os.Open` のエラーを適切にラップして返すこと。

（ファイルは `os.Open` で開き、`bufio.NewScanner` で読む）

### 問題2: カスタムエラー型

フォームのバリデーション用カスタムエラー型を作れ:

```go
type ValidationError struct {
    Field   string
    Message string
}
```

`validate(title, description string) error` 関数を実装し、
`errors.As` でエラーの詳細を取り出せるようにせよ。

### 問題3: Sentinel エラー

TODO 操作用の sentinel エラーを定義せよ:
- `ErrNotFound`: IDが存在しない
- `ErrAlreadyDone`: すでに完了済み

`CompleteTodo(todos []Todo, id int) error` 関数を実装し、各エラーを返すようにせよ。

---

## 解答

`exercises/safediv_solution/main.go`:

```go
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// 問題2: カスタムエラー型
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

// 問題3: Sentinel エラー
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

// 問題1: ファイルを読む
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
	// 問題1のテスト（存在しないファイル）
	_, err := readFirstLine("nonexistent.txt")
	if err != nil {
		fmt.Println("エラー:", err)
	}

	// 問題2のテスト
	err = validate("", "説明")
	var ve *ValidationError
	if errors.As(err, &ve) {
		fmt.Printf("バリデーションエラー - %s: %s\n", ve.Field, ve.Message)
	}

	err = validate("有効なタイトル", "説明")
	fmt.Println("バリデーション結果:", err) // <nil>

	// 問題3のテスト
	todos := []Todo{
		{ID: 1, Title: "Goを学ぶ", Done: false},
		{ID: 2, Title: "テストを書く", Done: true},
	}

	err = CompleteTodo(todos, 1)
	fmt.Println("Complete 1:", err) // <nil>

	err = CompleteTodo(todos, 2)
	if errors.Is(err, ErrAlreadyDone) {
		fmt.Println("すでに完了済みです")
	}

	err = CompleteTodo(todos, 99)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("見つかりません")
	}
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| エラー表現 | 例外（Exception）| 戻り値（error型）|
| エラー作成 | `raise ValueError("...")` | `errors.New("...")` |
| エラーチェック | `try/except` | `if err != nil` |
| エラーに情報付加 | `raise X from Y` | `fmt.Errorf("...: %w", err)` |
| エラー型確認 | `isinstance(e, MyError)` | `errors.As(err, &target)` |
| 回復不能エラー | `raise RuntimeError` | `panic(...)` |

次の章ではパッケージとモジュールを学びます。→ [第9章](../09_packages_modules/chapter09.md)
