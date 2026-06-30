package main

import (
	"log"
	"net/http"
	"todo-app/handler"
	"todo-app/store"
)

func main() {
	// JSON ファイルストアを初期化
	s := store.NewFileStore("todos.json")

	// ハンドラーを初期化
	h := handler.New(s)

	// ルーターにルートを登録
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("TODO アプリを起動: http://localhost:8080")
	log.Println("エンドポイント:")
	log.Println("  GET    /todos              - 一覧取得")
	log.Println("  POST   /todos              - 新規作成")
	log.Println("  GET    /todos/{id}         - 1件取得")
	log.Println("  PUT    /todos/{id}         - 更新")
	log.Println("  DELETE /todos/{id}         - 削除")
	log.Println("  PATCH  /todos/{id}/complete - 完了にする")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
