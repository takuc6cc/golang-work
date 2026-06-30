# 第2章: 変数・型・定数

## 変数の宣言

Go では変数を宣言する方法が3種類あります。

### パターン1: var キーワード（型を明示）

```go
var x int = 10
var name string = "Alice"
var isActive bool = true
```

### パターン2: var キーワード（型推論）

```go
var x = 10        // int と推論される
var name = "Alice" // string と推論される
```

### パターン3: 短縮宣言 `:=`（最もよく使う）

```go
x := 10
name := "Alice"
isActive := true
```

**Python との比較:**

```python
# Python: 型宣言なし、代入するだけ
x = 10
name = "Alice"
is_active = True
```

```go
// Go: 型が決まる（推論でも内部的には型あり）
x := 10           // int
name := "Alice"   // string
isActive := true  // bool
```

**重要: `:=` は関数の中でのみ使える。関数の外では `var` を使う。**

```go
package main

var globalVar = "グローバル"  // OK: パッケージレベルでは var

func main() {
    localVar := "ローカル"  // OK: 関数内では :=
    _ = localVar
}
```

### 複数変数の同時宣言

```go
var (
    x    int    = 1
    name string = "Alice"
    flag bool   = true
)

// または
a, b, c := 1, 2, 3
```

Python との比較:
```python
a, b, c = 1, 2, 3
```

---

## 基本的な型

| Go の型 | Python の型 | 説明 |
|---------|------------|------|
| `int` | `int` | 整数（64ビット環境では64ビット）|
| `int8` `int16` `int32` `int64` | — | 明示的なビット数 |
| `uint` | — | 符号なし整数 |
| `float32` `float64` | `float` | 浮動小数点数 |
| `string` | `str` | 文字列（UTF-8） |
| `bool` | `bool` | `true` / `false`（小文字！）|
| `byte` | — | `uint8` の別名、1バイト |
| `rune` | — | `int32` の別名、Unicode コードポイント |

**重要: Python の `True`/`False` と違い、Go は `true`/`false`（小文字）**

```go
flag := true   // OK
flag := True   // エラー！
```

### ゼロ値（Zero Value）

Go では変数を宣言すると自動的に「ゼロ値」が設定されます。
Python の `None` とは異なり、**型に応じた初期値**です。

```go
var i int     // 0
var f float64 // 0.0
var s string  // ""（空文字列）
var b bool    // false
```

Python との比較:
```python
# Python: 宣言と同時に代入しないと NameError
x = None  # 明示的に None を代入する
```

```go
// Go: 宣言だけで使える（ゼロ値が入る）
var count int
count++  // OK、count は 0 から始まる
```

---

## 型変換

Python では多くの場合、暗黙的な型変換が行われますが、Go では**明示的な変換が必要**です。

```python
# Python: 自動変換される場合もある
x = 3
y = 3.14
z = x + y  # 6.14（int + float → float）
```

```go
// Go: 明示的に変換しないとコンパイルエラー
x := 3
y := 3.14
// z := x + y  // エラー！型が違う

z := float64(x) + y  // OK: int を float64 に変換
```

よく使う型変換:

```go
i := 42
f := float64(i)    // int → float64
u := uint(f)       // float64 → uint
s := fmt.Sprintf("%d", i)  // int → string

// string → int は strconv パッケージを使う
import "strconv"
n, err := strconv.Atoi("42")  // "42" → 42
```

---

## 定数 const

```go
const Pi = 3.14159
const MaxRetry = 3
const AppName = "MyApp"
```

Python との比較:
```python
# Python: 慣習的に大文字で書くが、変更可能
PI = 3.14159
```

```go
// Go: const は本当に変更できない
const Pi = 3.14159
Pi = 3.14  // コンパイルエラー！
```

### iota: 連番定数

Python の `enum` に似た機能です。

```go
type Weekday int

const (
    Sunday Weekday = iota  // 0
    Monday                 // 1
    Tuesday                // 2
    Wednesday              // 3
    Thursday               // 4
    Friday                 // 5
    Saturday               // 6
)
```

Python との比較:
```python
from enum import Enum

class Weekday(Enum):
    SUNDAY = 0
    MONDAY = 1
    TUESDAY = 2
```

---

## fmt パッケージ: 出力

```go
fmt.Println("Hello")           // 改行あり
fmt.Print("Hello")             // 改行なし
fmt.Printf("名前: %s, 年齢: %d\n", name, age)  // 書式指定
s := fmt.Sprintf("名前: %s", name)  // 文字列として取得
```

### よく使うフォーマット指定子

| 指定子 | 意味 | 例 |
|--------|------|-----|
| `%d` | 整数 | `42` |
| `%f` | 浮動小数点数 | `3.140000` |
| `%.2f` | 小数点以下2桁 | `3.14` |
| `%s` | 文字列 | `Hello` |
| `%v` | 値（型に応じた表示） | `{Alice 30}` |
| `%T` | 型名 | `int`, `string` |
| `%t` | bool | `true` |
| `%p` | ポインタアドレス | `0xc000...` |

Python との比較:
```python
name = "Alice"
age = 30
print(f"名前: {name}, 年齢: {age}")
print("名前: {}, 年齢: {}".format(name, age))
```

```go
name := "Alice"
age := 30
fmt.Printf("名前: %s, 年齢: %d\n", name, age)
```

---

## A Tour of Go

本章に対応するセクション: **Basics: Packages, Variables, Types (1-17)**
https://go.dev/tour/basics/1

---

## 練習問題

### 問題1: 型変換と計算

`exercises/basics/main.go` に以下を実装してください。

```
摂氏 100 度は華氏 212.0 度です。
```

変換式: `F = C × 9/5 + 32`

条件:
- 摂氏を `int` 型の変数 `celsius` に入れる
- 計算に `float64` への型変換を使う
- `fmt.Printf` で `%.1f` を使って小数点以下1桁で表示する

### 問題2: iota で曜日定数

月曜〜日曜を `iota` で定義し、各曜日名と数値を出力せよ。

出力例:
```
月曜日: 0
火曜日: 1
水曜日: 2
木曜日: 3
金曜日: 4
土曜日: 5
日曜日: 6
```

### 問題3: fmt.Printf の練習

以下の情報を指定の形式で出力せよ。

```
名前: 山田花子
年齢: 25歳
身長: 162.5cm
Go好き: true
```

---

## 解答

`exercises/basics_solution/main.go`:

```go
package main

import "fmt"

// 問題2用の型と定数
type Weekday int

const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	// 問題1: 摂氏→華氏変換
	celsius := 100
	fahrenheit := float64(celsius)*9/5 + 32
	fmt.Printf("摂氏 %d 度は華氏 %.1f 度です。\n", celsius, fahrenheit)

	// 問題2: 曜日
	weekdays := []string{"月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日", "日曜日"}
	for i, day := range weekdays {
		fmt.Printf("%s: %d\n", day, i)
	}

	// 問題3: fmt.Printf
	name := "山田花子"
	age := 25
	height := 162.5
	likesGo := true
	fmt.Printf("名前: %s\n", name)
	fmt.Printf("年齢: %d歳\n", age)
	fmt.Printf("身長: %.1fcm\n", height)
	fmt.Printf("Go好き: %t\n", likesGo)
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| 変数宣言 | `x = 10` | `x := 10` または `var x int = 10` |
| 型推論 | 常に | `:=` または `var x = 10` |
| ゼロ値 | `None` を代入 | 自動的に初期値が入る |
| 型変換 | 多くは暗黙的 | 必ず明示的に `int(x)` |
| 定数 | 慣習のみ | `const` で強制 |
| 真偽値 | `True`/`False` | `true`/`false` |

次の章では制御フローを学びます。→ [第3章](../03_control_flow/chapter03.md)
