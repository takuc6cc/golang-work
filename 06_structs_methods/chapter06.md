# 第6章: 構造体・メソッド・ポインタ

## 構造体（Struct）

Python の `class` や `dataclass` に相当しますが、継承はありません。

```go
type Person struct {
    Name string
    Age  int
    Email string
}
```

Python との比較:
```python
from dataclasses import dataclass

@dataclass
class Person:
    name: str
    age: int
    email: str
```

### 構造体の初期化

```go
// フィールド名を指定（推奨）
p1 := Person{
    Name:  "Alice",
    Age:   30,
    Email: "alice@example.com",
}

// 順番通りに指定（フィールドが増えると壊れやすいので非推奨）
p2 := Person{"Bob", 25, "bob@example.com"}

// ゼロ値で初期化
var p3 Person  // Person{Name:"", Age:0, Email:""}

// new: ポインタを返す
p4 := new(Person)  // *Person 型、&Person{} と同じ
```

### フィールドへのアクセス

```go
p := Person{Name: "Alice", Age: 30}
fmt.Println(p.Name)  // Alice
p.Age = 31
```

---

## ポインタ

Go ではポインタを使うことで、変数の参照（メモリアドレス）を渡せます。
Python では基本的にすべてが参照渡しですが、Go の基本型や構造体は**値渡し**です。

```go
x := 42
p := &x        // p は x のアドレス（ポインタ）
fmt.Println(p)  // 0xc000014090（メモリアドレス）
fmt.Println(*p) // 42（デリファレンス: ポインタが指す値）

*p = 100       // ポインタ経由で値を変更
fmt.Println(x)  // 100（x が変わった）
```

### 値渡し vs ポインタ渡し

```go
// 値渡し: 関数内での変更は呼び出し元に影響しない
func double(n int) {
    n = n * 2  // コピーを変更しているので元の変数は変わらない
}

// ポインタ渡し: 関数内での変更が呼び出し元に影響する
func doublePtr(n *int) {
    *n = *n * 2  // ポインタ経由で元の変数を変更
}

x := 5
double(x)
fmt.Println(x)  // 5（変わっていない）

doublePtr(&x)
fmt.Println(x)  // 10（変わった！）
```

Python との比較:
```python
# Python: intは不変オブジェクトなので常に値コピー的な動作
def double(n):
    n = n * 2  # ローカル変数を変更するだけ

x = 5
double(x)
print(x)  # 5（変わっていない）

# リストは参照渡し的な動作
def append_to(lst):
    lst.append(4)

nums = [1, 2, 3]
append_to(nums)
print(nums)  # [1, 2, 3, 4]
```

---

## メソッド

構造体に関数を関連付けることができます。これをメソッドと呼びます。

```go
// レシーバ（メソッドの対象）を func の後に書く
func (p Person) Greet() string {
    return fmt.Sprintf("こんにちは、%sです！", p.Name)
}

p := Person{Name: "Alice"}
fmt.Println(p.Greet())  // こんにちは、Aliceです！
```

Python との比較:
```python
class Person:
    def __init__(self, name: str):
        self.name = name
    
    def greet(self) -> str:
        return f"こんにちは、{self.name}です！"
```

### 値レシーバ vs ポインタレシーバ

```go
type Counter struct {
    Count int
}

// 値レシーバ: コピーに対して操作（元の値は変わらない）
func (c Counter) Value() int {
    return c.Count
}

// ポインタレシーバ: 元の値を変更できる
func (c *Counter) Increment() {
    c.Count++
}

c := Counter{Count: 0}
c.Increment()  // Count が 1 に
c.Increment()  // Count が 2 に
fmt.Println(c.Value())  // 2
```

**使い分けの目安:**
- 構造体の状態を変更するメソッド → **ポインタレシーバ**
- 構造体の値を読むだけのメソッド → 値レシーバでもいいが、一貫性のためポインタレシーバを使うことが多い

---

## Stringer インターフェース

`fmt.Stringer` インターフェースを実装すると、`fmt.Println` での表示をカスタマイズできます。
Python の `__str__` メソッドに相当します。

```go
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

todo := Todo{ID: 1, Title: "Go を学ぶ", Done: true}
fmt.Println(todo)  // [x] 1: Go を学ぶ
```

Python との比較:
```python
class Todo:
    def __init__(self, id, title, done=False):
        self.id = id
        self.title = title
        self.done = done
    
    def __str__(self):
        mark = "x" if self.done else " "
        return f"[{mark}] {self.id}: {self.title}"
```

---

## 埋め込み（Embedding）

Go には継承がありませんが、「埋め込み」で構造体を組み合わせられます。

```go
type Animal struct {
    Name string
}

func (a Animal) Speak() string {
    return a.Name + " が鳴いています"
}

type Dog struct {
    Animal        // Animal を埋め込む
    Breed string
}

d := Dog{
    Animal: Animal{Name: "ポチ"},
    Breed:  "柴犬",
}
fmt.Println(d.Speak())  // ポチ が鳴いています（Animal のメソッドが使える）
fmt.Println(d.Name)     // ポチ（Animal のフィールドに直接アクセス）
```

Python との比較:
```python
class Animal:
    def __init__(self, name):
        self.name = name
    
    def speak(self):
        return f"{self.name} が鳴いています"

class Dog(Animal):  # 継承
    def __init__(self, name, breed):
        super().__init__(name)
        self.breed = breed
```

---

## A Tour of Go

本章に対応するセクション: **Methods and Interfaces (1-7)**
https://go.dev/tour/methods/1

---

## 練習問題

### 問題1: Person 構造体

`Person` 構造体（Name, Age フィールド）を定義し、
`Greet() string` メソッドを実装せよ。

出力例:
```
こんにちは！私はAliceです。30歳です。
```

### 問題2: BankAccount

`BankAccount` 構造体（Balance フィールド）を作り、以下のメソッドをポインタレシーバで実装せよ:
- `Deposit(amount int)`: 入金
- `Withdraw(amount int) error`: 出金（残高不足の場合はエラー）
- `GetBalance() int`: 残高確認

### 問題3: Todo 構造体

TODOアプリの準備として以下を実装せよ:
- `Todo` 構造体（ID int, Title string, Done bool）
- `Complete()` メソッド（Done を true にする）
- `String()` メソッド（`[x] 1: タイトル` 形式）

---

## 解答

`exercises/todo_model_solution/main.go`:

```go
package main

import (
	"errors"
	"fmt"
)

// 問題1
type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	return fmt.Sprintf("こんにちは！私は%sです。%d歳です。", p.Name, p.Age)
}

// 問題2
type BankAccount struct {
	Balance int
}

func (b *BankAccount) Deposit(amount int) {
	b.Balance += amount
}

func (b *BankAccount) Withdraw(amount int) error {
	if amount > b.Balance {
		return errors.New("残高不足です")
	}
	b.Balance -= amount
	return nil
}

func (b *BankAccount) GetBalance() int {
	return b.Balance
}

// 問題3
type Todo struct {
	ID    int
	Title string
	Done  bool
}

func (t *Todo) Complete() {
	t.Done = true
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
	p := Person{Name: "Alice", Age: 30}
	fmt.Println(p.Greet())

	// 問題2
	account := &BankAccount{Balance: 1000}
	account.Deposit(500)
	fmt.Printf("残高: %d円\n", account.GetBalance())

	if err := account.Withdraw(2000); err != nil {
		fmt.Println("エラー:", err)
	}
	account.Withdraw(300)
	fmt.Printf("残高: %d円\n", account.GetBalance())

	// 問題3
	todo := Todo{ID: 1, Title: "Goを学ぶ"}
	fmt.Println(todo)
	todo.Complete()
	fmt.Println(todo)
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| クラス | `class Foo:` | `type Foo struct { ... }` |
| インスタンス生成 | `Foo()` | `Foo{}` または `&Foo{}` |
| メソッド | `def method(self):` | `func (f Foo) Method() {` |
| 状態変更メソッド | `self.x = ...` | ポインタレシーバ `(f *Foo)` |
| `__str__` | `def __str__(self):` | `func (f Foo) String() string {` |
| 継承 | `class Dog(Animal):` | 埋め込み `type Dog struct { Animal }` |

次の章ではインターフェースを学びます。→ [第7章](../07_interfaces/chapter07.md)
