package bank

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type TransactionType string
type AccountType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Fee      TransactionType = "fee"
	Transfer TransactionType = "transfer"
	Giro     AccountType     = "giro"
	Savings  AccountType     = "savings"
)

type Transactions struct {
	Time    time.Time
	Purpose string
	Amount  float64
	Type    TransactionType
}

type Account struct {
	Id           string
	Name         string
	Balance      float64
	Overdraw     float64
	AccountType  AccountType
	Transactions []Transactions
}

var accounts = []Account{
	{
		Id:          "002",
		Name:        "Alice",
		Overdraw:    100.0,
		Balance:     1000.0,
		AccountType: Giro,
	},
	{
		Id:          "003",
		Name:        "Bob",
		Overdraw:    100.0,
		Balance:     500.0,
		AccountType: Savings,
	},
}

func (account *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount should be larger then 0")
	}

	account.Balance += amount

	account.addTransaction(amount, "Deposit", Deposit)
	return nil
}

func (account *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("mmount should be larger then 0")
	}

	overdraw := 0.0
	if account.AccountType == Giro {
		overdraw = -account.Overdraw
	}

	if account.Balance-amount < overdraw {
		return fmt.Errorf("Insufficient funds")
	}

	account.Balance -= amount

	account.addTransaction(amount, "Get Money", Withdraw)
	return nil
}

func (account *Account) Transfer(amount float64, to string) error {
	if amount <= 0 {
		return fmt.Errorf("Amount should be larger then 0")
	}

	if account.Balance-amount < 0 {
		return fmt.Errorf("Insufficient funds")
	}

	acc := searchingAcc(to)

	if acc == nil {
		return fmt.Errorf("No account found")
	}

	account.Balance -= amount
	acc.Balance += amount

	account.addTransaction(amount, "send it to "+to, Transfer)
	acc.addTransaction(amount, `received from {account.Name}`, Transfer)
	return nil
}

func (account *Account) ShowAccountDetails(w io.Writer, name string) error {

	if name == "" {
		fmt.Fprintf(w, "Balance: %.2f\n", account.Balance)
		for _, txn := range account.Transactions {
			fmt.Fprintf(w, "Time: %v, Purpose: %s, Amount: %.2f, Type: %v\n",
				txn.Time, txn.Purpose, txn.Amount, txn.Type)
		}
		return nil
	}

	acc := searchingAcc(name)
	if acc == nil {
		return fmt.Errorf("No account found with the name of: %s\n", name)
	}

	fmt.Fprintf(w, "Balance: %.2f\n", acc.Balance)
	for _, txn := range acc.Transactions {
		fmt.Fprintf(w, "Time: %v, Purpose: %s, Amount: %.2f, Type: %v\n",
			txn.Time, txn.Purpose, txn.Amount, txn.Type)
	}

	return nil
}

func searchingAcc(name string) *Account {

	for i := range accounts {
		if strings.EqualFold(accounts[i].Name, name) {
			return &accounts[i]
		}
		continue
	}

	return nil
}

func (account *Account) addTransaction(amount float64, purpose string, tt TransactionType) {
	account.Transactions = append(account.Transactions, Transactions{
		Time:    time.Now(),
		Purpose: purpose,
		Amount:  amount,
		Type:    tt,
	})
}
