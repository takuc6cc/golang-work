# 第7章: インターフェース

## インターフェースとは

インターフェースは「このメソッドを持っていれば使える」という**契約**を定義します。

Go のインターフェースの最大の特徴は**暗黙的な実装**です。
Python の `ABC（抽象基底クラス）`と違い、実装側が `implements` や継承を宣言する必要はありません。

```go
// インターフェースの定義
type Speaker interface {
    Speak() string
}

// Dog は Speaker を実装している（宣言不要）
type Dog struct {
    Name string
}
func (d Dog) Speak() string {
    return d.Name + "：ワン！"
}

// Cat も Speaker を実装している
type Cat struct {
    Name string
}
func (c Cat) Speak() string {
    return c.Name + "：ニャー！"
}

// Speaker インターフェース型の変数に代入できる
var s Speaker = Dog{Name: "ポチ"}
fmt.Println(s.Speak())  // ポチ：ワン！

s = Cat{Name: "タマ"}
fmt.Println(s.Speak())  // タマ：ニャー！
```

Python との比較:
```python
from abc import ABC, abstractmethod

class Speaker(ABC):
    @abstractmethod
    def speak(self) -> str:
        pass

class Dog(Speaker):  # 明示的に継承が必要
    def __init__(self, name):
        self.name = name
    
    def speak(self) -> str:
        return f"{self.name}：ワン！"
```

---

## インターフェースを使う利点

インターフェースを引数に受け取る関数を作ると、異なる型を統一的に扱えます。

```go
func makeNoise(speakers []Speaker) {
    for _, s := range speakers {
        fmt.Println(s.Speak())
    }
}

animals := []Speaker{
    Dog{Name: "ポチ"},
    Cat{Name: "タマ"},
    Dog{Name: "ハチ"},
}
makeNoise(animals)
```

---

## 標準ライブラリのインターフェース

Go 標準ライブラリにはよく使われるインターフェースが多数あります。

### fmt.Stringer

`String() string` メソッドを持つ型は `fmt.Println` などで自動的に使われます。

```go
type Point struct {
    X, Y int
}

func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

pt := Point{X: 3, Y: 4}
fmt.Println(pt)  // (3, 4)
```

### error インターフェース

```go
type error interface {
    Error() string
}
```

`Error() string` メソッドを持つ型はすべて `error` として扱えます（詳細は第8章）。

### io.Reader / io.Writer

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

ファイル、HTTP レスポンス、バッファなど「読める/書けるもの」を統一的に扱えます。

---

## 空インターフェース

`interface{}` または Go 1.18+ の `any` は、あらゆる型を格納できます。
Python の `object` 型に相当します。

```go
var v interface{}  // または: var v any

v = 42
v = "Hello"
v = []int{1, 2, 3}
```

主に型が不明な場合に使いますが、型安全性が失われるので多用は避けます。

---

## 型アサーション

インターフェース型から具体的な型を取り出します。

```go
var i interface{} = "Hello"

// 型アサーション（失敗するとパニック）
s := i.(string)
fmt.Println(s)  // Hello

// 安全な型アサーション（失敗しても ok で確認できる）
s, ok := i.(string)
if ok {
    fmt.Println("string:", s)
} else {
    fmt.Println("string ではない")
}
```

Python との比較:
```python
v = "Hello"
if isinstance(v, str):
    print("string:", v)
```

---

## 型スイッチ

複数の型を分岐処理する場合は型スイッチが便利です。

```go
func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("整数: %d\n", v)
    case string:
        fmt.Printf("文字列: %s\n", v)
    case bool:
        fmt.Printf("真偽値: %t\n", v)
    default:
        fmt.Printf("不明な型: %T\n", v)
    }
}

describe(42)        // 整数: 42
describe("Hello")   // 文字列: Hello
describe(true)      // 真偽値: true
```

Python との比較:
```python
def describe(v):
    match v:
        case int():
            print(f"整数: {v}")
        case str():
            print(f"文字列: {v}")
        case bool():
            print(f"真偽値: {v}")
        case _:
            print(f"不明な型: {type(v)}")
```

---

## TODOアプリへの応用: Storage インターフェース

インターフェースを使うと、テスト用のモックと本番用の実装を切り替えやすくなります。

```go
// インターフェースを定義
type Store interface {
    GetAll() ([]Todo, error)
    GetByID(id int) (Todo, error)
    Create(todo Todo) (Todo, error)
    Update(todo Todo) (Todo, error)
    Delete(id int) error
}

// 本番用: JSON ファイルに保存
type FileStore struct {
    path string
}

// テスト用: メモリに保存
type MemoryStore struct {
    todos []Todo
}

// Handler はインターフェースだけを知っている（実装に依存しない）
type Handler struct {
    store Store
}
```

---

## A Tour of Go

本章に対応するセクション: **Methods and Interfaces (8-18)**
https://go.dev/tour/methods/9

---

## 練習問題

### 問題1: Shape インターフェース

`Shape` インターフェース（`Area() float64`、`Perimeter() float64`）を定義し、
`Circle` と `Rectangle` で実装せよ。

```go
func printShape(s Shape) {
    fmt.Printf("面積: %.2f, 周囲: %.2f\n", s.Area(), s.Perimeter())
}
```

### 問題2: 型スイッチ

任意の値を受け取り、型と値を表示する関数 `describe(v interface{})` を実装せよ。

```
describe(42)     → 整数: 42
describe("Hi")   → 文字列: Hi
describe(3.14)   → 浮動小数点数: 3.14
describe(true)   → 真偽値: true
```

### 問題3: Stringer 実装

第6章で作った `Todo` 型に `String()` メソッドを実装し、
`fmt.Println(todo)` で `[x] 1: タスクのタイトル` と表示されるようにせよ。

---

## 解答

`exercises/shape_solution/main.go`:

```go
package main

import (
	"fmt"
	"math"
)

// 問題1: Shape インターフェース
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func printShape(s Shape) {
	fmt.Printf("面積: %.2f, 周囲: %.2f\n", s.Area(), s.Perimeter())
}

// 問題2: 型スイッチ
func describe(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Printf("整数: %d\n", val)
	case string:
		fmt.Printf("文字列: %s\n", val)
	case float64:
		fmt.Printf("浮動小数点数: %.2f\n", val)
	case bool:
		fmt.Printf("真偽値: %t\n", val)
	default:
		fmt.Printf("不明な型: %T\n", val)
	}
}

// 問題3: Todo
type Todo struct {
	ID    int
	Title string
	Done  bool
}

func (t Todo) String() string {
	mark := " "
	if t.Done {
		mark = "x"
	}
	return fmt.Sprintf("[%s] %d: %s", mark, t.ID, t.Title)
}

func main() {
	// 問題1
	c := Circle{Radius: 5}
	r := Rectangle{Width: 4, Height: 6}
	printShape(c)
	printShape(r)

	// 問題2
	describe(42)
	describe("Hello")
	describe(3.14)
	describe(true)

	// 問題3
	todo := Todo{ID: 1, Title: "Goを学ぶ", Done: true}
	fmt.Println(todo)
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| インターフェース定義 | `class Foo(ABC):` | `type Foo interface { ... }` |
| 実装の宣言 | `class Bar(Foo):` 必要 | 不要（自動的に判定）|
| 型チェック | `isinstance(v, Foo)` | `v, ok := i.(Foo)` |
| 型スイッチ | `match type(v):` | `switch v := i.(type)` |
| 何でも入る型 | `object` / `Any` | `interface{}` / `any` |

次の章ではエラーハンドリングを学びます。→ [第8章](../08_error_handling/chapter08.md)
