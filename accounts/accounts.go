package accounts

import (
	"errors"
	"fmt"
)

var errNoMoney = errors.New("Can't without you are pool")

// Account struct
type Account struct {
	owner  string
	blance int
}

// NewAccount creates Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, blance: 0}
	return &account //memory address
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	fmt.Println("Gonna deposit", amount, "hahaha")
	a.blance += amount
}

func (a Account) Blance() int {
	return a.blance
}

// Withdraw x amount from your account
func (a *Account) Withdraw(amount int) error {
	if a.blance < amount {
		return errNoMoney
	}
	a.blance -= amount
	return nil
}
