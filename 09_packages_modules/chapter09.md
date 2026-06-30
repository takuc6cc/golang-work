# 第9章: パッケージ・モジュール

## モジュールとパッケージの概念

| 概念 | Python | Go |
|------|--------|----|
| プロジェクト単位 | なし（慣習） | **モジュール**（`go.mod`）|
| コード単位 | **パッケージ**（ディレクトリ）| **パッケージ**（ディレクトリ）|
| 依存関係管理 | `pip` + `requirements.txt` | `go get` + `go.mod` + `go.sum` |
| 仮想環境 | `venv` | 不要（モジュールが役割を担う）|

---

## モジュールの作成

```bash
# プロジェクトディレクトリを作って初期化
mkdir myapp
cd myapp
go mod init github.com/yourname/myapp
```

これで `go.mod` ファイルが作られます:

```
module github.com/yourname/myapp

go 1.22
```

Python との比較:
```bash
# Python
mkdir myapp && cd myapp
python -m venv venv
source venv/bin/activate
pip install requests
pip freeze > requirements.txt
```

```bash
# Go
mkdir myapp && cd myapp
go mod init github.com/yourname/myapp
go get github.com/some/package  # 依存を追加
```

---

## パッケージの分割

1つのディレクトリ = 1つのパッケージです。

```
myapp/
├── go.mod
├── main.go        ← package main
├── model/
│   └── user.go    ← package model
├── handler/
│   └── handler.go ← package handler
└── store/
    └── store.go   ← package store
```

### パッケージの公開・非公開

**大文字で始まる名前 = 公開（外部から使える）**
**小文字で始まる名前 = 非公開（パッケージ内だけで使える）**

```go
// model/user.go
package model

type User struct {       // 公開: 外部から使える
    ID    int
    Name  string
    email string        // 非公開: model パッケージ内だけで使える
}

func NewUser(name string) *User {  // 公開
    return &User{Name: name}
}

func (u *User) validate() bool {   // 非公開
    return u.Name != ""
}
```

Python との比較:
```python
class User:
    def __init__(self, name):
        self.name = name       # 公開（慣習）
        self._email = ""       # 非公開（慣習のみ、強制ではない）
        self.__id = 0          # より強い非公開（name mangling）
```

---

## import のルール

```go
import (
    "fmt"                           // 標準ライブラリ
    "net/http"                      // 標準ライブラリ（パス区切りに /）
    "github.com/yourname/myapp/model" // 自分のパッケージ
    "github.com/gin-gonic/gin"      // サードパーティ
)
```

使うときはパッケージ名でアクセスします:

```go
model.User{}
http.HandleFunc(...)
gin.Default()
```

### エイリアス

```go
import (
    myfmt "fmt"  // fmt を myfmt という名前で使う
    _ "github.com/lib/pq"  // 副作用のみ（init() 実行目的）
)
```

---

## よく使う標準ライブラリ

### fmt: フォーマット出力

```go
fmt.Println("Hello")
fmt.Printf("名前: %s\n", name)
s := fmt.Sprintf("名前: %s", name)
```

### strings: 文字列操作

```go
import "strings"

strings.Contains("Hello, World", "World")  // true
strings.HasPrefix("Hello", "He")           // true
strings.ToUpper("hello")                   // "HELLO"
strings.TrimSpace("  hello  ")             // "hello"
strings.Split("a,b,c", ",")               // ["a", "b", "c"]
strings.Join([]string{"a", "b"}, ", ")    // "a, b"
strings.Replace("aaa", "a", "b", 2)       // "bba"（2回だけ置換）
strings.ReplaceAll("aaa", "a", "b")       // "bbb"
```

Python との比較:
```python
"Hello, World".find("World") != -1   # Contains
"Hello".startswith("He")             # HasPrefix
"hello".upper()                      # ToUpper
"  hello  ".strip()                  # TrimSpace
"a,b,c".split(",")                   # Split
", ".join(["a", "b"])                # Join
```

### strconv: 型変換

```go
import "strconv"

n, err := strconv.Atoi("42")          // string → int
s := strconv.Itoa(42)                  // int → string
f, err := strconv.ParseFloat("3.14", 64) // string → float64
b, err := strconv.ParseBool("true")    // string → bool
```

Python との比較:
```python
int("42")       # Atoi
str(42)         # Itoa
float("3.14")   # ParseFloat
```

### os: OS操作

```go
import "os"

os.ReadFile("path.txt")                 // ファイル読み込み
os.WriteFile("path.txt", data, 0644)   // ファイル書き込み
os.Remove("path.txt")                  // ファイル削除
os.Getenv("HOME")                      // 環境変数取得
os.Exit(1)                             // プロセス終了
```

### time: 日時操作

```go
import "time"

now := time.Now()                          // 現在時刻
formatted := now.Format("2006-01-02 15:04:05")  // フォーマット（Goの魔法の日時）
time.Sleep(2 * time.Second)               // 2秒待つ
oneHourLater := now.Add(time.Hour)        // 1時間後
diff := time.Since(start)                 // 経過時間
```

**注意: Go の時刻フォーマットは `2006-01-02 15:04:05` という特定の日時を使います。**
（この日時は Go 言語の誕生時刻に由来する）

Python との比較:
```python
from datetime import datetime, timedelta
now = datetime.now()
formatted = now.strftime("%Y-%m-%d %H:%M:%S")
one_hour_later = now + timedelta(hours=1)
```

---

## サードパーティパッケージの追加

```bash
# パッケージを追加
go get github.com/some/package

# go.mod と go.sum が自動更新される
# go.sum にはチェックサムが記録される（セキュリティ）

# 未使用パッケージを削除
go mod tidy
```

---

## 練習問題

### 問題1: パッケージ分割

以下の構造で `mypackage` モジュールを作れ:

```
exercises/mypackage/
├── go.mod          ← module mypackage
├── main.go         ← greet.Hello("Alice") を呼ぶ
└── greet/
    └── greet.go    ← package greet; func Hello(name string) string
```

### 問題2: strings パッケージ

`strings` パッケージを使って以下を実装せよ:

```go
// "Hello, World! Hello, Go!" から "Hello" の出現回数を数える
count := strings.Count("Hello, World! Hello, Go!", "Hello")  // 2

// メールアドレスの @ より前を取得する
parts := strings.SplitN("alice@example.com", "@", 2)
username := parts[0]  // "alice"
```

### 問題3: time パッケージ

以下を出力するプログラムを書け:

```
現在時刻: 2024-01-15 10:30:00
1時間後: 2024-01-15 11:30:00
経過時間: 0s
```

---

## 解答

`exercises/mypackage_solution/greet/greet.go`:

```go
package greet

import "fmt"

// Hello は名前を受け取り挨拶文を返す
func Hello(name string) string {
	return fmt.Sprintf("こんにちは、%s！", name)
}

// Goodbye は名前を受け取りお別れの言葉を返す
func Goodbye(name string) string {
	return fmt.Sprintf("さようなら、%s！", name)
}
```

`exercises/mypackage_solution/main.go`:

```go
package main

import (
	"fmt"
	"strings"
	"time"

	"mypackage/greet"
)

func main() {
	// 問題1: greet パッケージ
	fmt.Println(greet.Hello("Alice"))
	fmt.Println(greet.Goodbye("Bob"))

	// 問題2: strings パッケージ
	text := "Hello, World! Hello, Go!"
	count := strings.Count(text, "Hello")
	fmt.Printf("\"Hello\" の出現回数: %d\n", count)

	parts := strings.SplitN("alice@example.com", "@", 2)
	fmt.Printf("ユーザー名: %s\n", parts[0])

	// 問題3: time パッケージ
	now := time.Now()
	fmt.Println("現在時刻:", now.Format("2006-01-02 15:04:05"))
	fmt.Println("1時間後:", now.Add(time.Hour).Format("2006-01-02 15:04:05"))

	start := time.Now()
	// 何か処理...
	elapsed := time.Since(start)
	fmt.Printf("経過時間: %v\n", elapsed)
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| プロジェクト初期化 | `venv` + `pip` | `go mod init` |
| 依存追加 | `pip install pkg` | `go get pkg` |
| 依存ファイル | `requirements.txt` | `go.mod` + `go.sum` |
| 公開/非公開 | `_` 接頭辞（慣習）| 大文字/小文字（強制）|
| 文字列操作 | `str` メソッド | `strings` パッケージ |
| 型変換 | `int()`, `str()` | `strconv.Atoi`, `strconv.Itoa` |

次の章では並行処理を学びます。→ [第10章](../10_concurrency/chapter10.md)
