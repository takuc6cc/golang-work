# 第1章: Go言語入門・環境構築

## Go言語とは

Go（Golang）は Google が2009年に公開したプログラミング言語です。

**特徴:**
- **コンパイル型**: Python のように逐次実行ではなく、実行前にバイナリに変換
- **静的型付け**: 変数の型はコンパイル時に決まる（Python は動的型付け）
- **ガベージコレクション**: メモリ管理は自動（C/C++と違い手動解放不要）
- **並行処理に強い**: goroutine による軽量な並行処理
- **シンプルな文法**: キーワードが25個しかない

### Python/PHP との思想比較

| 項目 | Python | PHP | Go |
|------|--------|-----|-----|
| 型付け | 動的 | 動的 | 静的 |
| 実行方式 | インタープリタ | インタープリタ | コンパイル |
| 速度 | 遅い | 中程度 | 速い |
| 並行処理 | 難しい（GIL） | 難しい | 得意 |
| 文法の複雑さ | 中 | 中 | シンプル |

---

## 環境構築

### Go のインストール

1. https://go.dev/dl/ から OS に合わせたインストーラをダウンロード
2. インストール後、ターミナルで確認:

```bash
go version
# go version go1.22.0 windows/amd64
```

### 重要な環境変数

```bash
go env GOPATH   # Goのワークスペース（通常 ~/go）
go env GOROOT   # Goのインストール先
```

Python との対比:

| Python | Go |
|--------|-----|
| `python --version` | `go version` |
| `PYTHONPATH` | `GOPATH` |
| インタープリタ本体 | `GOROOT` |

---

## 最初の Go プログラム

### ファイルを作成して実行

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

**Python との比較:**

```python
# Python: ファイルを作ってそのまま実行
print("Hello, Python!")
```

```go
// Go: package宣言とmain関数が必須
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

### 実行方法の違い

```bash
# Python
python hello.py

# Go: その場で実行（コンパイル+実行を自動でやってくれる）
go run main.go

# Go: バイナリにコンパイルして実行
go build -o hello main.go
./hello          # Linux/Mac
hello.exe        # Windows
```

---

## Go のコード規則

### package main と func main()

- すべての Go ファイルは `package 名前` で始まる
- 実行可能なプログラムは必ず `package main` にする
- プログラムのエントリポイントは `func main()` （引数・戻り値なし）

### import

```go
import "fmt"          // 1つだけ
import (              // 複数
    "fmt"
    "strings"
)
```

Python との比較:
```python
import os
from os import path
```

```go
import "os"
// Goはパッケージ全体をインポートし、os.xxx のように使う
// from os import path のような部分インポートはない
```

### コメント

```go
// 1行コメント（Python の # に相当）

/*
  複数行コメント
  （Python の """ とは違い、文字列ではない）
*/
```

### セミコロン

Go にはセミコロンが**ありません**（PHPと違い不要）。
ただし Go コンパイラが自動的に行末にセミコロンを挿入するルールがあるため、**開き中括弧 `{` は行末に書く**必要があります。

```go
// OK
func main() {
    fmt.Println("OK")
}

// NG（コンパイルエラー）
func main()
{
    fmt.Println("NG")
}
```

### go fmt: 自動フォーマット

Python の `black`、PHP の PHP-CS-Fixer に相当するツールが標準搭載されています。

```bash
go fmt main.go       # ファイルを整形
go fmt ./...         # プロジェクト全体を整形
```

VS Code に Go 拡張を入れていれば保存時に自動実行されます。

---

## A Tour of Go

A Tour of Go（https://go.dev/tour/welcome/1）は Go 公式のインタラクティブチュートリアルです。
ブラウザ上で Go コードを書いて実行できるため、環境構築前でも学習を始められます。

本章に対応するセクション: **Welcome (1-4)**

---

## 練習問題

### 問題1: Hello, World!

`exercises/hello/main.go` に以下を実装してください。

- `"Hello, Go!"` を出力するプログラムを書け

```bash
go run exercises/hello/main.go
# Hello, Go!
```

### 問題2: 自己紹介

自分の名前と年齢を以下の形式で出力するプログラムを書け。
`fmt.Printf` を使うこと。

```
私の名前は田中太郎です。年齢は30歳です。
```

ヒント: `fmt.Printf("私の名前は%sです。年齢は%d歳です。\n", name, age)`

### 問題3: go build でバイナリ生成

`go build -o hello main.go` でバイナリを生成し、直接実行せよ。
（go run を使わずに実行できることを確認する）

---

## 解答

`exercises/hello/main.go` の解答例:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")

    // 問題2
    name := "田中太郎"
    age := 30
    fmt.Printf("私の名前は%sです。年齢は%d歳です。\n", name, age)
}
```

---

## まとめ

| 概念 | Python | Go |
|------|--------|----|
| 実行 | `python main.py` | `go run main.go` |
| ビルド | — | `go build` |
| フォーマット | `black` | `go fmt` |
| エントリポイント | ファイルの先頭から実行 | `func main()` |
| コメント | `#` | `//` |
| セミコロン | 不要 | 不要（自動挿入） |

次の章では変数と型を学びます。→ [第2章](../02_basics/chapter02.md)
