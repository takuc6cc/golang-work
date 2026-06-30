package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// ヘルパー: JSON レスポンスを書く
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// 問題1: GET /hello で {"message": "Hello, World!"} を返す
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// ここにコードを書いてください
	_ = w
	_ = r
}

// 問題2: GET /greet?name=Alice で {"message": "Hello, Alice!"} を返す
func greetHandler(w http.ResponseWriter, r *http.Request) {
	// name パラメータを取得（なければ "World"）
	// ここにコードを書いてください
	_ = w
	_ = r
}

// 問題3: POST /echo でリクエストボディをそのまま返す
func echoHandler(w http.ResponseWriter, r *http.Request) {
	// ここにコードを書いてください
	_ = io.ReadAll
	_ = w
	_ = r
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)
	mux.HandleFunc("GET /greet", greetHandler)
	mux.HandleFunc("POST /echo", echoHandler)

	log.Println("サーバーを起動: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
