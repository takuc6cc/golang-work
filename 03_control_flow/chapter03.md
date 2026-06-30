# 第3章: 制御フロー

## if 文

基本的な書き方は Python と似ていますが、**条件式を括弧で囲まない**のが特徴です。

```go
x := 10
if x > 5 {
    fmt.Println("5より大きい")
} else if x == 5 {
    fmt.Println("5と等しい")
} else {
    fmt.Println("5より小さい")
}
```

Python との比較:
```python
x = 10
if x > 5:
    print("5より大きい")
elif x == 5:
    print("5と等しい")
else:
    print("5より小さい")
```

### 初期化ステートメント付き if

Go の if 文には「初期化ステートメント」を書ける独自の構文があります。
エラーハンドリングでよく使います。

```go
if err := someFunction(); err != nil {
    fmt.Println("エラー:", err)
    return
}
```

`err` は `if` ブロックのスコープ内だけで有効です。

---

## for 文

**Go には `while` がありません。** すべてのループを `for` で書きます。
Python ユーザーにとってここが最初の戸惑いポイントです。

### パターン1: C言語スタイル（カウンタ付き）

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

Python との比較:
```python
for i in range(5):
    print(i)
```

### パターン2: while 相当（条件のみ）

```go
i := 0
for i < 5 {
    fmt.Println(i)
    i++
}
```

Python との比較:
```python
i = 0
while i < 5:
    print(i)
    i += 1
```

### パターン3: 無限ループ

```go
for {
    // 永遠に繰り返す
    // break で抜ける
}
```

Python との比較:
```python
while True:
    pass
```

### パターン4: range を使ったループ（最もよく使う）

```go
// スライス（Pythonのリスト）
fruits := []string{"apple", "banana", "cherry"}
for i, fruit := range fruits {
    fmt.Printf("%d: %s\n", i, fruit)
}

// インデックスが不要な場合は _ で無視
for _, fruit := range fruits {
    fmt.Println(fruit)
}

// マップ（Pythonのdict）
scores := map[string]int{"Alice": 90, "Bob": 75}
for name, score := range scores {
    fmt.Printf("%s: %d\n", name, score)
}

// 文字列
for i, ch := range "Hello" {
    fmt.Printf("%d: %c\n", i, ch)
}
```

Python との比較:
```python
fruits = ["apple", "banana", "cherry"]
for i, fruit in enumerate(fruits):
    print(f"{i}: {fruit}")

scores = {"Alice": 90, "Bob": 75}
for name, score in scores.items():
    print(f"{name}: {score}")
```

### break と continue

Python と同じように使えます。

```go
for i := 0; i < 10; i++ {
    if i == 3 {
        continue  // 3 をスキップ
    }
    if i == 7 {
        break  // 7 で終了
    }
    fmt.Println(i)
}
// 0 1 2 4 5 6
```

---

## switch 文

Python には `match` 文（3.10+）がありますが、Go の `switch` はより強力で古くから使えます。

**重要: Go の switch は各 case の末尾に `break` を書く必要がありません。**
Python/PHP と逆で、**自動的に break** されます。

```go
day := "Monday"
switch day {
case "Saturday", "Sunday":
    fmt.Println("週末")
case "Monday":
    fmt.Println("月曜日")
default:
    fmt.Println("平日")
}
```

Python との比較:
```python
day = "Monday"
match day:
    case "Saturday" | "Sunday":
        print("週末")
    case "Monday":
        print("月曜日")
    case _:
        print("平日")
```

### fallthrough

通常は自動的に break されますが、`fallthrough` で次のケースを実行できます（まれにしか使わない）。

```go
n := 2
switch n {
case 1:
    fmt.Println("1")
    fallthrough
case 2:
    fmt.Println("2")  // ここが実行される
    fallthrough
case 3:
    fmt.Println("3")  // これも実行される
}
// 2
// 3
```

### 条件なし switch（if-else チェーンの代替）

```go
score := 85
switch {
case score >= 90:
    fmt.Println("A")
case score >= 80:
    fmt.Println("B")
case score >= 70:
    fmt.Println("C")
default:
    fmt.Println("D")
}
```

### 初期化ステートメント付き switch

```go
switch os := runtime.GOOS; os {
case "darwin":
    fmt.Println("macOS")
case "linux":
    fmt.Println("Linux")
default:
    fmt.Printf("%s\n", os)
}
```

---

## A Tour of Go

本章に対応するセクション: **Basics: Flow Control (1-13)**
https://go.dev/tour/flowcontrol/1

---

## 練習問題

### 問題1: FizzBuzz

1〜100 の数値に対して:
- 3の倍数なら `Fizz`
- 5の倍数なら `Buzz`
- 両方の倍数なら `FizzBuzz`
- それ以外は数値をそのまま出力

`for` ループと `switch` を使って実装してください。

`exercises/fizzbuzz/main.go` に書いてください。

### 問題2: 素数を探せ

1〜100 の素数をすべて出力せよ。

出力例:
```
2 3 5 7 11 13 17 ...
```

ヒント: 外側の `for` で各数値を、内側の `for` で割り切れるか調べる。

### 問題3: スライスの最大値

スライス `nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}` から最大値を求めて出力せよ。

`for range` を使うこと。

---

## 解答

`exercises/fizzbuzz_solution/main.go`:

```go
package main

import "fmt"

func main() {
	// 問題1: FizzBuzz
	fmt.Println("=== FizzBuzz ===")
	for i := 1; i <= 100; i++ {
		switch {
		case i%15 == 0:
			fmt.Println("FizzBuzz")
		case i%3 == 0:
			fmt.Println("Fizz")
		case i%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(i)
		}
	}

	// 問題2: 素数
	fmt.Println("=== 素数 ===")
	for n := 2; n <= 100; n++ {
		isPrime := true
		for i := 2; i*i <= n; i++ {
			if n%i == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			fmt.Printf("%d ", n)
		}
	}
	fmt.Println()

	// 問題3: 最大値
	fmt.Println("=== 最大値 ===")
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	max := nums[0]
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	fmt.Printf("最大値: %d\n", max)
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| if | `if x > 5:` | `if x > 5 {` |
| elif/else if | `elif` | `else if` |
| for カウンタ | `for i in range(5):` | `for i := 0; i < 5; i++ {` |
| while | `while cond:` | `for cond {` |
| for each | `for x in list:` | `for _, x := range list {` |
| switch | `match` (3.10+) | `switch { case ...: }` |
| break 必要？ | 不要（match は） | 不要（自動break）|

次の章では関数を学びます。→ [第4章](../04_functions/chapter04.md)
