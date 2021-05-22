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

// ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner of the account
func (a Account) Owner() string {
	return a.owner
}

func (a Account) String() string {
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Blance())
}
