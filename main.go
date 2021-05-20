package main

import (
	"fmt"

	"github.com/maengsoo/learngo/accounts"
)

func main() {
	account := accounts.NewAccount("maengsoo")
	account.Deposit(10)
	fmt.Println(account.Blance())
	err := account.Withdraw(20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(account.Blance())
}
