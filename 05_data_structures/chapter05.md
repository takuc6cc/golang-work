# 第5章: データ構造（スライス・マップ）

## 配列（Array）

Go には固定長の配列がありますが、実際の開発ではほとんど使いません。
後述のスライスを使います。

```go
var a [3]int         // [0 0 0]
b := [3]int{1, 2, 3} // [1 2 3]
c := [...]int{1, 2, 3} // 長さ自動（[3]int と同じ）
```

**配列はスライスと違い、長さが型の一部です。** `[3]int` と `[4]int` は別の型です。

---

## スライス（Slice）

Python の `list` に相当する、Go で最もよく使うコレクション型です。
長さが動的に変わります。

```go
// 宣言と初期化
fruits := []string{"apple", "banana", "cherry"}

// 空のスライス
var empty []int        // nil スライス
empty2 := []int{}      // 空スライス（nilではない）

// make で作成（長さとキャパシティを指定）
s := make([]int, 3)       // [0 0 0]、len=3, cap=3
s2 := make([]int, 3, 10)  // [0 0 0]、len=3, cap=10
```

### append: 要素の追加

```go
s := []int{1, 2, 3}
s = append(s, 4)          // [1 2 3 4]
s = append(s, 5, 6, 7)    // [1 2 3 4 5 6 7]

// 別のスライスを展開して追加
other := []int{8, 9}
s = append(s, other...)   // [1 2 3 4 5 6 7 8 9]
```

Python との比較:
```python
s = [1, 2, 3]
s.append(4)          # [1, 2, 3, 4]
s.extend([5, 6, 7])  # [1, 2, 3, 4, 5, 6, 7]
```

**重要:** `append` は新しいスライスを返します。必ず `s = append(s, ...)` のように受け取る必要があります。

### スライス操作

```go
s := []int{0, 1, 2, 3, 4, 5}
s[1:4]   // [1 2 3]（インデックス1から3まで）
s[:3]    // [0 1 2]（先頭から3まで）
s[3:]    // [3 4 5]（インデックス3から最後まで）
s[:]     // [0 1 2 3 4 5]（全体のコピー参照）
```

Python との比較:
```python
s = [0, 1, 2, 3, 4, 5]
s[1:4]   # [1, 2, 3]
s[:3]    # [0, 1, 2]
s[3:]    # [3, 4, 5]
```

### len と cap

```go
s := make([]int, 3, 10)
fmt.Println(len(s))  // 3（現在の長さ）
fmt.Println(cap(s))  // 10（確保済みの容量）
```

### copy: スライスのコピー

```go
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)  // src を dst にコピー
```

Python との比較:
```python
src = [1, 2, 3]
dst = src.copy()   # または dst = src[:]
```

---

## マップ（Map）

Python の `dict` に相当します。キーと値のペアを保持します。

```go
// 宣言と初期化
scores := map[string]int{
    "Alice": 90,
    "Bob":   75,
    "Carol": 85,
}

// make で作成
m := make(map[string]int)
```

### 要素の追加・更新・取得

```go
m := make(map[string]int)
m["Alice"] = 90    // 追加
m["Alice"] = 95    // 更新

score := m["Alice"]  // 取得: 95
score = m["Dave"]    // 存在しないキー: ゼロ値 (0) が返る
```

### キーの存在確認

Go では存在しないキーにアクセスするとゼロ値が返るため、存在確認には2値取得を使います。

```go
score, ok := m["Alice"]
if ok {
    fmt.Printf("Alice のスコア: %d\n", score)
} else {
    fmt.Println("Alice は存在しない")
}
```

Python との比較:
```python
scores = {"Alice": 90}

# Python: キーエラーを避ける方法
if "Alice" in scores:
    print(scores["Alice"])

# または
score = scores.get("Alice", 0)  # デフォルト値
```

### delete: 要素の削除

```go
delete(m, "Alice")
```

Python との比較:
```python
del scores["Alice"]
scores.pop("Alice", None)  # エラーにしない場合
```

### range でイテレーション

```go
for key, value := range scores {
    fmt.Printf("%s: %d\n", key, value)
}

// キーのみ
for key := range scores {
    fmt.Println(key)
}
```

Python との比較:
```python
for key, value in scores.items():
    print(f"{key}: {value}")

for key in scores:
    print(key)
```

**注意: マップのイテレーション順序は保証されません（Python 3.7+ では挿入順）**

---

## 2次元スライス

```go
// 3行4列の2次元スライス
matrix := [][]int{
    {1, 2, 3, 4},
    {5, 6, 7, 8},
    {9, 10, 11, 12},
}
fmt.Println(matrix[1][2])  // 7
```

Python との比較:
```python
matrix = [
    [1, 2, 3, 4],
    [5, 6, 7, 8],
    [9, 10, 11, 12],
]
print(matrix[1][2])  # 7
```

---

## A Tour of Go

本章に対応するセクション: **More Types: Arrays, Slices, Maps (6-23)**
https://go.dev/tour/moretypes/6

---

## 練習問題

### 問題1: 偶数フィルター

整数スライスを受け取り、偶数のみを含む新しいスライスを返す関数 `filterEven` を実装せよ。

```go
result := filterEven([]int{1, 2, 3, 4, 5, 6, 7, 8})
// [2 4 6 8]
```

### 問題2: 単語カウント

文字列のスライスを受け取り、各単語の出現回数をマップで返す関数 `wordCount` を実装せよ。

```go
words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
counts := wordCount(words)
// map[apple:3 banana:2 cherry:1]
```

### 問題3: 電話帳

マップを使って簡単な電話帳を実装せよ。以下の操作を行うこと:

1. 連絡先を3件追加
2. 名前で検索して番号を表示
3. 1件削除
4. 全件表示

---

## 解答

`exercises/wordcount_solution/main.go`:

```go
package main

import "fmt"

// 問題1: 偶数フィルター
func filterEven(nums []int) []int {
	result := []int{}
	for _, n := range nums {
		if n%2 == 0 {
			result = append(result, n)
		}
	}
	return result
}

// 問題2: 単語カウント
func wordCount(words []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}
	return counts
}

func main() {
	// 問題1
	result := filterEven([]int{1, 2, 3, 4, 5, 6, 7, 8})
	fmt.Println(result)

	// 問題2
	words := []string{"apple", "banana", "apple", "cherry", "banana", "apple"}
	counts := wordCount(words)
	fmt.Println(counts)

	// 問題3: 電話帳
	phonebook := make(map[string]string)

	// 追加
	phonebook["田中太郎"] = "090-1234-5678"
	phonebook["山田花子"] = "080-8765-4321"
	phonebook["鈴木一郎"] = "070-1111-2222"

	// 検索
	if num, ok := phonebook["田中太郎"]; ok {
		fmt.Printf("田中太郎: %s\n", num)
	}

	// 削除
	delete(phonebook, "山田花子")

	// 全件表示
	fmt.Println("=== 電話帳 ===")
	for name, num := range phonebook {
		fmt.Printf("%s: %s\n", name, num)
	}
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| 動的配列 | `list` | `[]T`（スライス）|
| 追加 | `list.append(x)` | `s = append(s, x)` |
| 結合 | `list.extend(other)` | `s = append(s, other...)` |
| 辞書 | `dict` | `map[K]V` |
| キー存在確認 | `key in d` | `v, ok := m[key]` |
| 削除 | `del d[key]` | `delete(m, key)` |

次の章では構造体とメソッドを学びます。→ [第6章](../06_structs_methods/chapter06.md)
