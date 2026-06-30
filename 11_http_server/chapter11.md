# 第11章: HTTP サーバー

## net/http パッケージ

Go の標準ライブラリ `net/http` だけで完全な HTTP サーバーを作れます。
フレームワーク（Gin、Echo など）不要です。

### 最小構成の HTTP サーバー

```go
package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/hello", helloHandler)
    http.ListenAndServe(":8080", nil)
}
```

Python（Flask）との比較:
```python
from flask import Flask
app = Flask(__name__)

@app.route("/hello")
def hello():
    return "Hello, World!"

app.run(port=8080)
```

---

## ResponseWriter と Request

### http.ResponseWriter: レスポンスを書く

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // ヘッダーを設定（Write より前に呼ぶ）
    w.Header().Set("Content-Type", "application/json")
    
    // ステータスコードを設定（WriteHeader より前にヘッダーを設定する）
    w.WriteHeader(http.StatusCreated)  // 201
    
    // ボディを書く
    fmt.Fprintln(w, `{"message": "created"}`)
}
```

### *http.Request: リクエストを読む

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // メソッド確認
    fmt.Println(r.Method)  // "GET", "POST", etc.
    
    // URL パス
    fmt.Println(r.URL.Path)  // "/todos/123"
    
    // クエリパラメータ
    name := r.URL.Query().Get("name")  // ?name=Alice → "Alice"
    
    // リクエストボディ（POST/PUT）
    body, err := io.ReadAll(r.Body)
    defer r.Body.Close()
}
```

---

## JSON レスポンス

```go
import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Message string `json:"message"`
    Count   int    `json:"count"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    resp := Response{Message: "Hello", Count: 42}
    json.NewEncoder(w).Encode(resp)
}
```

### JSON リクエストを受け取る

```go
type CreateTodoRequest struct {
    Title string `json:"title"`
}

func createHandler(w http.ResponseWriter, r *http.Request) {
    var req CreateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "リクエストが不正です", http.StatusBadRequest)
        return
    }
    
    // req.Title を使って処理...
}
```

---

## HTTP メソッドの振り分け

Go 1.22 以降では `http.ServeMux` がメソッドとパスの組み合わせに対応しました。

```go
mux := http.NewServeMux()

// メソッド + パスの組み合わせ
mux.HandleFunc("GET /todos", listTodosHandler)
mux.HandleFunc("POST /todos", createTodoHandler)
mux.HandleFunc("GET /todos/{id}", getTodoHandler)
mux.HandleFunc("PUT /todos/{id}", updateTodoHandler)
mux.HandleFunc("DELETE /todos/{id}", deleteTodoHandler)

http.ListenAndServe(":8080", mux)
```

### パスパラメータの取得（Go 1.22+）

```go
func getTodoHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")  // /todos/123 → "123"
    // ...
}
```

---

## ステータスコードの定数

```go
http.StatusOK                  // 200
http.StatusCreated             // 201
http.StatusNoContent           // 204
http.StatusBadRequest          // 400
http.StatusUnauthorized        // 401
http.StatusForbidden           // 403
http.StatusNotFound            // 404
http.StatusInternalServerError // 500
```

---

## ミドルウェア

ミドルウェアは `http.Handler` をラップする関数です。

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
    })
}

// 使い方
mux := http.NewServeMux()
mux.HandleFunc("GET /hello", helloHandler)

// ミドルウェアで包む
http.ListenAndServe(":8080", loggingMiddleware(mux))
```

Python（Flask デコレータ）との比較:
```python
from functools import wraps
import time

def log_time(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        start = time.time()
        result = f(*args, **kwargs)
        print(f"{time.time() - start:.3f}s")
        return result
    return wrapper
```

---

## A Tour of Go

本章に対応する公式ドキュメント:
https://pkg.go.dev/net/http

---

## 練習問題

### 問題1: Hello サーバー

`GET /hello` で `{"message": "Hello, World!"}` を返すサーバーを作れ。

```bash
go run main.go &
curl http://localhost:8080/hello
# {"message":"Hello, World!"}
```

### 問題2: クエリパラメータ

`GET /greet?name=Alice` でクエリパラメータを読み取り、
`{"message": "Hello, Alice!"}` を返すエンドポイントを追加せよ。
`name` が指定されない場合は `"World"` にする。

### 問題3: エコーサーバー

`POST /echo` でリクエストボディの JSON をそのまま返すエンドポイントを作れ。

```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"key": "value"}'
# {"key":"value"}
```

### 問題4: ロギングミドルウェア

リクエストのメソッド・パス・処理時間を出力するミドルウェアを実装せよ。

```
GET /hello 1.234ms
POST /echo 0.567ms
```

---

## 解答

`exercises/hello_server_solution/main.go`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// 問題1: Hello ハンドラー
func helloHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Message{Message: "Hello, World!"})
}

// 問題2: Greet ハンドラー
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	writeJSON(w, http.StatusOK, Message{Message: fmt.Sprintf("Hello, %s!", name)})
}

// 問題3: Echo ハンドラー
func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "読み取りエラー", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// 問題4: ロギングミドルウェア
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("GET /greet", greetHandler)
	mux.HandleFunc("POST /echo", echoHandler)

	log.Println("サーバーを起動: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", loggingMiddleware(mux)))
}
```

---

## まとめ

| 概念 | Python (Flask) | Go (net/http) |
|------|----------------|---------------|
| ルート定義 | `@app.route("/")` | `mux.HandleFunc("GET /", h)` |
| レスポンス | `return jsonify(data)` | `json.NewEncoder(w).Encode(data)` |
| ステータスコード | `return data, 201` | `w.WriteHeader(http.StatusCreated)` |
| リクエストボディ | `request.json` | `json.NewDecoder(r.Body).Decode(&v)` |
| クエリパラメータ | `request.args.get("key")` | `r.URL.Query().Get("key")` |

次の章では JSON とファイルストレージを学びます。→ [第12章](../12_json_storage/chapter12.md)
