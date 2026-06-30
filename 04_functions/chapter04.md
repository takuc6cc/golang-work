# 第4章: 関数

## 基本的な関数

```go
func 関数名(引数名 型, ...) 戻り値の型 {
    // 処理
    return 値
}
```

```go
func add(a int, b int) int {
    return a + b
}

// 引数の型が同じ場合はまとめられる
func add(a, b int) int {
    return a + b
}
```

Python との比較:
```python
def add(a: int, b: int) -> int:
    return a + b
```

```go
func add(a, b int) int {
    return a + b
}
```

---

## 多値返却（Go最大の特徴の一つ）

Go では複数の値を返すことができます。
エラーハンドリングで特に重要です。

```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("ゼロ除算はできません")
    }
    return a / b, nil
}

// 呼び出し側
result, err := divide(10, 2)
if err != nil {
    fmt.Println("エラー:", err)
    return
}
fmt.Println(result) // 5
```

Python との比較:
```python
# Python: タプルで返す（型ヒントあり）
def divide(a: float, b: float) -> tuple[float, str | None]:
    if b == 0:
        return 0.0, "ゼロ除算はできません"
    return a / b, None

result, err = divide(10, 2)
if err:
    print("エラー:", err)
```

```go
// Go: 多値返却がネイティブ機能
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("ゼロ除算はできません")
    }
    return a / b, nil
}
result, err := divide(10, 2)
```

**不要な戻り値は `_` で無視する:**

```go
result, _ := divide(10, 2)  // エラーを無視（注意: 本番コードでは避ける）
```

---

## 名前付き戻り値

戻り値に名前をつけることができます。関数内で変数として使えます。

```go
func minMax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums {
        if n < min {
            min = n
        }
        if n > max {
            max = n
        }
    }
    return  // naked return: min と max が自動的に返る
}
```

---

## 可変長引数

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)        // 6
sum(1, 2, 3, 4, 5)  // 15

// スライスを展開して渡す
nums := []int{1, 2, 3}
sum(nums...)  // ... で展開
```

Python との比較:
```python
def sum_all(*nums: int) -> int:
    return sum(nums)

nums = [1, 2, 3]
sum_all(*nums)  # * でアンパック
```

---

## ファーストクラス関数

Go では関数を変数に代入したり、引数として渡せます。

```go
// 関数を変数に代入
add := func(a, b int) int {
    return a + b
}
result := add(3, 4)  // 7

// 関数を引数として渡す
func apply(f func(int, int) int, a, b int) int {
    return f(a, b)
}
apply(add, 3, 4)  // 7
```

Python との比較:
```python
add = lambda a, b: a + b
result = add(3, 4)

def apply(f, a, b):
    return f(a, b)
apply(add, 3, 4)
```

---

## クロージャ

関数が外側のスコープの変数を「捕捉（キャプチャ）」する機能です。

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

counter := makeCounter()
fmt.Println(counter())  // 1
fmt.Println(counter())  // 2
fmt.Println(counter())  // 3
```

Python との比較:
```python
def make_counter():
    count = 0
    def counter():
        nonlocal count
        count += 1
        return count
    return counter

counter = make_counter()
print(counter())  # 1
print(counter())  # 2
```

### フィボナッチ数列クロージャ（A Tour of Go の例題）

```go
func fibonacci() func() int {
    a, b := 0, 1
    return func() int {
        result := a
        a, b = b, a+b
        return result
    }
}

fib := fibonacci()
for i := 0; i < 10; i++ {
    fmt.Printf("%d ", fib())
}
// 0 1 1 2 3 5 8 13 21 34
```

---

## defer

`defer` は関数が**終了する直前**に実行されます。
ファイルのクローズやロックの解放によく使います。

```go
func readFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // 関数終了時に自動的に実行される

    // ファイル処理...
    return nil
}
```

Python との比較:
```python
# Python: with文でリソース管理
with open("file.txt") as f:
    # 処理
    pass  # ブロックを抜けると自動でclose
```

```go
// Go: defer でリソース管理
f, err := os.Open("file.txt")
if err != nil { return err }
defer f.Close()  // 忘れずにここに書く
// 以降の処理...
```

### defer の実行順序

複数の `defer` は**後入れ先出し（LIFO）**で実行されます。

```go
func main() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
}
// 出力:
// 3
// 2
// 1
```

---

## A Tour of Go

本章に対応するセクション:
- https://go.dev/tour/basics/4 （関数）
- https://go.dev/tour/moretypes/25 （クロージャ）
- https://go.dev/tour/flowcontrol/12 （defer）

---

## 練習問題

### 問題1: 四則演算関数

2つの整数を受け取り、和・差・積・商（float64）をすべて返す関数 `calculate` を実装せよ。

```go
sum, diff, product, quotient := calculate(10, 3)
// 13, 7, 30, 3.333...
```

### 問題2: 可変長合計

可変長引数を受け取り合計を返す `sum(...int) int` を実装せよ。
さらに、スライスを `...` で展開して渡せることも確認せよ。

### 問題3: defer の練習

`defer` を使って以下の出力が得られる関数を書け:

```
開始
処理中...
終了（defer で実行）
```

### 問題4: フィボナッチクロージャ

フィボナッチ数列を順番に返すクロージャを実装せよ。
呼び出すたびに次の値を返すこと。

```go
fib := fibonacci()
fmt.Println(fib()) // 0
fmt.Println(fib()) // 1
fmt.Println(fib()) // 1
fmt.Println(fib()) // 2
fmt.Println(fib()) // 3
```

---

## 解答

`exercises/calculator_solution/main.go`:

```go
package main

import (
	"fmt"
)

// 問題1: 四則演算
func calculate(a, b int) (int, int, int, float64) {
	return a + b, a - b, a * b, float64(a) / float64(b)
}

// 問題2: 可変長合計
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// 問題4: フィボナッチクロージャ
func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		result := a
		a, b = b, a+b
		return result
	}
}

// 問題3: defer
func process() {
	defer fmt.Println("終了（defer で実行）")
	fmt.Println("開始")
	fmt.Println("処理中...")
}

func main() {
	// 問題1
	s, d, p, q := calculate(10, 3)
	fmt.Printf("和:%d 差:%d 積:%d 商:%.3f\n", s, d, p, q)

	// 問題2
	fmt.Println(sum(1, 2, 3, 4, 5))
	nums := []int{10, 20, 30}
	fmt.Println(sum(nums...))

	// 問題3
	process()

	// 問題4
	fib := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", fib())
	}
	fmt.Println()
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| 関数定義 | `def f(a, b):` | `func f(a, b int) int {` |
| 多値返却 | タプルで代用 | ネイティブ機能 |
| 可変長引数 | `*args` | `...int` |
| ラムダ/無名関数 | `lambda x: x+1` | `func(x int) int { return x+1 }` |
| クロージャ | `nonlocal` で外部変数を参照 | 自動的にキャプチャ |
| リソース管理 | `with` | `defer` |

次の章ではデータ構造（スライス・マップ）を学びます。→ [第5章](../05_data_structures/chapter05.md)
