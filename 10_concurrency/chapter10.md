# 第10章: 並行処理（基礎）

## goroutine とは

goroutine は Go の軽量スレッドです。`go` キーワードを関数の前に付けるだけで起動できます。

```go
func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

// 通常の呼び出し（同期）
sayHello("Alice")

// goroutine として起動（非同期）
go sayHello("Bob")
```

Python との比較:
```python
import threading

def say_hello(name):
    print(f"Hello, {name}!")

# 通常の呼び出し
say_hello("Alice")

# スレッドとして起動
t = threading.Thread(target=say_hello, args=("Bob",))
t.start()
t.join()
```

**重要:** `main` 関数が終わると、起動した goroutine も強制終了します。
goroutine の完了を待つ仕組みが必要です。

---

## sync.WaitGroup: goroutine の完了を待つ

```go
import "sync"

var wg sync.WaitGroup

wg.Add(3)  // 待つ goroutine の数を設定

go func() {
    defer wg.Done()  // goroutine 終了時に呼ぶ
    fmt.Println("goroutine 1")
}()

go func() {
    defer wg.Done()
    fmt.Println("goroutine 2")
}()

go func() {
    defer wg.Done()
    fmt.Println("goroutine 3")
}()

wg.Wait()  // すべての goroutine が Done を呼ぶまで待つ
fmt.Println("すべて完了")
```

Python との比較:
```python
import threading

threads = []
for i in range(3):
    t = threading.Thread(target=lambda i=i: print(f"thread {i+1}"))
    threads.append(t)
    t.start()

for t in threads:
    t.join()
print("すべて完了")
```

---

## channel: goroutine 間のデータ送受信

channel は goroutine 間でデータを安全に渡すためのパイプです。

```go
// channel の作成
ch := make(chan int)      // バッファなし
ch2 := make(chan int, 5)  // バッファサイズ5

// 送信
ch <- 42

// 受信
value := <-ch

// close で送信終了を知らせる
close(ch)
```

### バッファなし channel（同期）

```go
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i  // 受信側が受け取るまでブロック
    }
    close(ch)
}

func main() {
    ch := make(chan int)
    go producer(ch)

    for n := range ch {  // close されるまで受信
        fmt.Println(n)
    }
}
```

### バッファあり channel（非同期）

```go
ch := make(chan string, 3)
ch <- "first"    // バッファに入る（ブロックしない）
ch <- "second"
ch <- "third"
// ch <- "fourth"  // バッファが満杯でブロック

fmt.Println(<-ch)  // "first"
fmt.Println(<-ch)  // "second"
```

Python との比較:
```python
import queue
import threading

q = queue.Queue()

def producer():
    for i in range(5):
        q.put(i)

threading.Thread(target=producer).start()
while True:
    item = q.get()
    print(item)
    if item == 4:
        break
```

---

## sync.Mutex: データ競合の防止

複数の goroutine が同じデータに同時にアクセスすると「データ競合」が発生します。
`sync.Mutex` でロックして保護します。

```go
import "sync"

type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

Python との比較:
```python
import threading

class SafeCounter:
    def __init__(self):
        self._lock = threading.Lock()
        self._count = 0
    
    def increment(self):
        with self._lock:
            self._count += 1
```

---

## HTTP サーバーと goroutine

`net/http` の HTTP サーバーは**リクエストごとに自動的に goroutine を起動**します。
そのため、Handler 関数は複数の goroutine から同時に呼ばれる可能性があります。

```go
type Handler struct {
    store Store
    mu    sync.RWMutex  // 読み書きロック
}

func (h *Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
    h.mu.RLock()  // 読み込みロック
    defer h.mu.RUnlock()
    // ...
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    h.mu.Lock()  // 書き込みロック
    defer h.mu.Unlock()
    // ...
}
```

---

## A Tour of Go

本章に対応するセクション: **Concurrency (1-11)**
https://go.dev/tour/concurrency/1

---

## 練習問題

### 問題1: WaitGroup の基本

3つの goroutine を起動し、それぞれ異なるメッセージを出力させよ。
`WaitGroup` で全て完了するのを待つこと。

### 問題2: channel でパイプライン

以下のパイプラインを実装せよ:

1. `generate` goroutine: 1〜10 の整数を channel に送信
2. `square` goroutine: channel から受け取り、2乗を別の channel に送信
3. main: 結果を受け取って出力

```
1 → 1
4 → 4
9 → 9
...
```

### 問題3: goroutine セーフカウンター

`sync.Mutex` を使った goroutine セーフなカウンターを実装せよ。
100個の goroutine を起動し、それぞれ100回インクリメントして、
最終的に `10000` になることを確認せよ。

---

## 解答

`exercises/goroutine_basic_solution/main.go`:

```go
package main

import (
	"fmt"
	"sync"
)

// 問題3: goroutine セーフカウンター
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// 問題2: パイプライン
func generate(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		ch <- i
	}
	close(ch)
}

func square(in <-chan int, out chan<- int) {
	for n := range in {
		out <- n * n
	}
	close(out)
}

func main() {
	// 問題1: WaitGroup
	fmt.Println("=== WaitGroup ===")
	var wg sync.WaitGroup
	messages := []string{"Hello from goroutine 1", "Hello from goroutine 2", "Hello from goroutine 3"}
	for _, msg := range messages {
		wg.Add(1)
		go func(m string) {
			defer wg.Done()
			fmt.Println(m)
		}(msg)
	}
	wg.Wait()
	fmt.Println("すべての goroutine 完了")

	// 問題2: パイプライン
	fmt.Println("\n=== パイプライン ===")
	naturals := make(chan int)
	squares := make(chan int)

	go generate(naturals)
	go square(naturals, squares)

	for n := range squares {
		fmt.Printf("%d\n", n)
	}

	// 問題3: goroutine セーフカウンター
	fmt.Println("\n=== goroutine セーフカウンター ===")
	counter := &SafeCounter{}
	var wg2 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}
	wg2.Wait()
	fmt.Printf("最終カウント: %d（期待値: 10000）\n", counter.Value())
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| 軽量スレッド | `threading.Thread` | `goroutine`（`go func()`）|
| 完了待ち | `thread.join()` | `wg.Wait()` |
| データ転送 | `queue.Queue` | `channel`（`chan`）|
| 排他制御 | `threading.Lock` | `sync.Mutex` |
| goroutine 数 | 数百〜数千（重い）| 数百万でも軽い |

次の章では HTTP サーバーを学びます。→ [第11章](../11_http_server/chapter11.md)
