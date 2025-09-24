package bank

import (
	"errors"
	"fmt"
	"io"
	"time"
)

type TransactionType int16
type AccountType int16

const (
	Deposit TransactionType = iota
	Withdraw
	Fee
	Transfer
)

func (tt TransactionType) String() string{
	switch tt {
	case Deposit:
		return "deposit"
	case Withdraw:
		return "withdraw"
	case Fee:
		return "fee"
	case Transfer:
		return "transfer"
	}
	return "unknown"
}

const (
	Giro AccountType = iota
	Savings
)

type Transactions struct {
	Time time.Time
	Purpose string
	Amount float64
	Type TransactionType
}

type Account struct {
	Id string
	Name string
	Balance float64
	AccountType AccountType
	Transactions []Transactions
}

var accounts = []Account{
    {
        Id: "002",
        Name: "Alice",
        Balance: 1000.0,
        AccountType: Giro,
    },
    {
        Id: "003",
        Name: "Bob",
        Balance: 500.0,
        AccountType: Savings,
    },
  }


func (account *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("Amount should be larger then 0")
	}

	account.Balance += amount

	account.addTransaction(amount, "Deposit", Deposit)
	return nil
}

func (account *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("Amount should be larger then 0")
	}

	if account.Balance-amount < 0 {
		return errors.New("Insufficient funds")
	}

	account.Balance -= amount

	account.addTransaction(amount, "Get Money", Withdraw)
	return nil
}

func (account *Account) Transfer(amount float64, to string) error {
	if amount <= 0 {
		return errors.New("Amount should be larger then 0")
	}

	if account.Balance-amount < 0{
		return errors.New("Insufficient funds")
	}

	acc := searchingAcc(to)

	if acc == nil {
		return errors.New("No account found")
	}

	account.Balance -= amount
	acc.Balance += amount

	account.addTransaction(amount, "send it to " + to, Transfer )
	acc.addTransaction(amount, `received from {account.Name}`, Transfer )
	return nil
}

func (account *Account) ShowAccountDetails(w io.Writer, name string) {

	if name == "" {
		fmt.Fprintf(w, "Balance: %.2f\n", account.Balance)
		for _, txn := range account.Transactions {
			fmt.Fprintf(w, "Time: %v, Purpose: %s, Amount: %.2f, Type: %v\n",
				txn.Time, txn.Purpose, txn.Amount, txn.Type)
		}
		return
	}
	 
	acc := searchingAcc(name)
	if acc == nil {
		fmt.Errorf("No account found with the name of: %s", name)
		return
	}

	fmt.Fprintf(w, "Balance: %.2f\n", acc.Balance)
	for _, txn := range acc.Transactions {
		fmt.Fprintf(w, "Time: %v, Purpose: %s, Amount: %.2f, Type: %v\n",
			txn.Time, txn.Purpose, txn.Amount, txn.Type.String())
	}

}

func searchingAcc(searchingAccount string) *Account {
	
	for _, account := range accounts {
		if account.Name == searchingAccount {
			return &account
		}
		continue
	}
	
	return nil 
}

func (account *Account)addTransaction(amount float64, purpose string, tt TransactionType) {
	account.Transactions = append(account.Transactions, Transactions {
		Time: time.Now(),
		Purpose: purpose,
		Amount: amount,
		Type: tt,
	})
}
