package main

import "fmt"

type Weekday int

const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	// 問題1: 摂氏→華氏変換
	celsius := 100
	fahrenheit := float64(celsius)*9/5 + 32
	fmt.Printf("摂氏 %d 度は華氏 %.1f 度です。\n", celsius, fahrenheit)

	// 問題2: 曜日
	weekdays := []string{"月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日", "日曜日"}
	for i, day := range weekdays {
		fmt.Printf("%s: %d\n", day, i)
	}

	// 問題3: fmt.Printf
	name := "山田花子"
	age := 25
	height := 162.5
	likesGo := true
	fmt.Printf("名前: %s\n", name)
	fmt.Printf("年齢: %d歳\n", age)
	fmt.Printf("身長: %.1fcm\n", height)
	fmt.Printf("Go好き: %t\n", likesGo)
}
