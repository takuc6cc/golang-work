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
