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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Message{Message: "Hello, World!"})
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	writeJSON(w, http.StatusOK, Message{Message: fmt.Sprintf("Hello, %s!", name)})
}

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
