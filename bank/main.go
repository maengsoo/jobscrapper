package main

import (
	"fmt"

	"github.com/maengsoo/learngo/bank/accounts"
)

func main() {
	account := accounts.NewAccount("maengsoo")
	account.Deposit(10)
	err := account.Withdraw(20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(account)
}
