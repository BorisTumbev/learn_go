package pointers

import (
	"errors"
	"fmt"
)

type Bitcoin int

// type Stringer interface {
// 	String() string
// }

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

type Wallet struct {
	balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
	fmt.Printf("address of balance in Deposit is %p \n", &w.balance)
	w.balance += amount
}

// Technically you do not need to change Balance to use a pointer receiver as taking a copy of the balance is fine.
//
//	However, by convention you should keep your method receiver types the same for consistency.
func (w *Wallet) Balance() Bitcoin {
	// (*w).balance is the same
	return w.balance
}

var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return ErrInsufficientFunds
	}
	w.balance -= amount
	return nil
}
