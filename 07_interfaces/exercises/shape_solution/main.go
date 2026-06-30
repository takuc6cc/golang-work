package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func printShape(s Shape) {
	fmt.Printf("面積: %.2f, 周囲: %.2f\n", s.Area(), s.Perimeter())
}

func describe(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Printf("整数: %d\n", val)
	case string:
		fmt.Printf("文字列: %s\n", val)
	case float64:
		fmt.Printf("浮動小数点数: %.2f\n", val)
	case bool:
		fmt.Printf("真偽値: %t\n", val)
	default:
		fmt.Printf("不明な型: %T\n", val)
	}
}

type Todo struct {
	ID    int
	Title string
	Done  bool
}

func (t Todo) String() string {
	mark := " "
	if t.Done {
		mark = "x"
	}
	return fmt.Sprintf("[%s] %d: %s", mark, t.ID, t.Title)
}

func main() {
	c := Circle{Radius: 5}
	r := Rectangle{Width: 4, Height: 6}
	printShape(c)
	printShape(r)

	describe(42)
	describe("Hello")
	describe(3.14)
	describe(true)

	todo := Todo{ID: 1, Title: "Goを学ぶ", Done: true}
	fmt.Println(todo)
}
