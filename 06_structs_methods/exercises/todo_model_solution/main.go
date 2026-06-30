package main

import (
	"errors"
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	return fmt.Sprintf("こんにちは！私は%sです。%d歳です。", p.Name, p.Age)
}

type BankAccount struct {
	Balance int
}

func (b *BankAccount) Deposit(amount int) {
	b.Balance += amount
}

func (b *BankAccount) Withdraw(amount int) error {
	if amount > b.Balance {
		return errors.New("残高不足です")
	}
	b.Balance -= amount
	return nil
}

func (b *BankAccount) GetBalance() int {
	return b.Balance
}

type Todo struct {
	ID    int
	Title string
	Done  bool
}

func (t *Todo) Complete() {
	t.Done = true
}

func (t Todo) String() string {
	mark := " "
	if t.Done {
		mark = "x"
	}
	return fmt.Sprintf("[%s] %d: %s", mark, t.ID, t.Title)
}

func main() {
	p := Person{Name: "Alice", Age: 30}
	fmt.Println(p.Greet())

	account := &BankAccount{Balance: 1000}
	account.Deposit(500)
	fmt.Printf("残高: %d円\n", account.GetBalance())
	if err := account.Withdraw(2000); err != nil {
		fmt.Println("エラー:", err)
	}
	account.Withdraw(300)
	fmt.Printf("残高: %d円\n", account.GetBalance())

	todo := Todo{ID: 1, Title: "Goを学ぶ"}
	fmt.Println(todo)
	todo.Complete()
	fmt.Println(todo)
}
