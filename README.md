# Go言語学習カリキュラム

Python/PHPエンジニアがGo言語を習得し、**TODOアプリを自力で作れるレベル**を目指すカリキュラムです。

## 対象者

- Python または PHP の基本的な読み書きができる
- Go言語は未経験
- Web アプリ（CRUD + HTTP API）を作ることが最終目標

## 学習の進め方

1. 各章の `chapterXX.md` を読む
2. 章末の練習問題を自力で解く（`exercises/問題名/main.go` に書く）
3. 解答と照らし合わせる（`exercises/問題名_solution/main.go`）
4. 第14章で TODO アプリを完成させる

[A Tour of Go](https://go.dev/tour/)と並行して学ぶと効果的です。

---

## 目次

| 章 | タイトル | 学習時間目安 |
|----|----------|-------------|
| [第1章](01_introduction/chapter01.md) | 入門・環境構築 | 1〜2時間 |
| [第2章](02_basics/chapter02.md) | 変数・型・定数 | 2〜3時間 |
| [第3章](03_control_flow/chapter03.md) | 制御フロー | 2〜3時間 |
| [第4章](04_functions/chapter04.md) | 関数 | 3〜4時間 |
| [第5章](05_data_structures/chapter05.md) | データ構造（スライス・マップ） | 3〜4時間 |
| [第6章](06_structs_methods/chapter06.md) | 構造体・メソッド・ポインタ | 4〜5時間 |
| [第7章](07_interfaces/chapter07.md) | インターフェース | 4〜5時間 |
| [第8章](08_error_handling/chapter08.md) | エラーハンドリング | 3〜4時間 |
| [第9章](09_packages_modules/chapter09.md) | パッケージ・モジュール | 2〜3時間 |
| [第10章](10_concurrency/chapter10.md) | 並行処理（基礎） | 3〜4時間 |
| [第11章](11_http_server/chapter11.md) | HTTP サーバー | 4〜5時間 |
| [第12章](12_json_storage/chapter12.md) | JSON・ファイルストレージ | 3〜4時間 |
| [第13章](13_testing/chapter13.md) | テスト | 3〜4時間 |
| [第14章](14_todo_app/chapter14.md) | TODOアプリ実装 | 6〜10時間 |

**合計目安: 約43〜61時間**

---

## 事前準備

### Go のインストール

https://go.dev/dl/ から最新版をダウンロードしてインストールします。

```bash
# インストール確認
go version
# 例: go version go1.22.0 windows/amd64
```

### エディタ設定

**VS Code** の場合:
1. 拡張機能「Go」（ms-vscode.go）をインストール
2. コマンドパレット → `Go: Install/Update Tools` → 全選択してインストール

### 動作確認

```bash
mkdir hello && cd hello
```

`main.go` を作成:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

```bash
go run main.go
# Hello, Go!
```

---

## A Tour of Go 対応表

| 本カリキュラム | A Tour of Go |
|---------------|-------------|
| 第1章 | Welcome (1-4) |
| 第2章 | Basics: Packages, Variables, Types (1-17) |
| 第3章 | Basics: Flow Control (1-13) |
| 第4章 | Basics: Functions + More Types: Closures |
| 第5章 | More Types: Arrays, Slices, Maps (6-23) |
| 第6章 | Methods and Interfaces (1-7) |
| 第7章 | Methods and Interfaces (8-18) |
| 第10章 | Concurrency (1-11) |

---

## 最終成果物: TODO アプリ

第14章で完成する TODO アプリの機能:

- タスクの追加・一覧表示・更新・削除
- 完了フラグの切り替え
- JSON ファイルへの永続化
- REST API（`net/http` のみ使用、フレームワーク不使用）

```
GET    /todos              # 一覧取得
POST   /todos              # 新規作成
GET    /todos/{id}         # 1件取得
PUT    /todos/{id}         # 更新
DELETE /todos/{id}         # 削除
PATCH  /todos/{id}/complete # 完了にする
```
