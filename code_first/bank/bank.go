package bank

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type TransactionType string
type AccountType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
	Transfer TransactionType = "transfer"
	Giro     AccountType     = "giro"
	Savings  AccountType     = "savings"
)

type Transactions struct {
	Time   time.Time
	Amount float64
	Type   TransactionType
}

type Account struct {
	Id           string
	Name         string
	Balance      float64
	Overdraw     float64
	AccountType  AccountType
	Transactions []Transactions
}

var initialAccounts = []Account{
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

func InitialAccounts() {
	for _, acc := range initialAccounts {
		AddOrUpdateAcc(&acc)
	}
}

func (account *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("Amount should be larger then 0")
	}

	account.Balance += amount

	account.addTransaction(amount, Deposit)
	AddOrUpdateAcc(account)
	return nil
}

func (account *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount should be larger then 0")
	}

	overdraw := 0.0
	if account.AccountType == Giro {
		overdraw = -account.Overdraw
	}

	if account.Balance-amount < overdraw {
		return fmt.Errorf("Insufficient funds")
	}

	account.Balance -= amount

	account.addTransaction(amount, Withdraw)
	AddOrUpdateAcc(account)
	return nil
}

func (account *Account) Transfer(amount float64, to string) error {
	if amount <= 0 {
		return fmt.Errorf("Amount should be larger then 0")
	}

	if account.Balance-amount < 0 {
		return fmt.Errorf("Insufficient funds")
	}

	recipientAcc, err := searchingAcc(to)
	if err != nil {
		return fmt.Errorf("unexcepteced error: %v\n", err)
	}

	account.Balance -= amount
	recipientAcc.Balance += amount

	account.addTransaction(amount, Transfer)
	recipientAcc.addTransaction(amount, Transfer)
	AddOrUpdateAcc(account)
	AddOrUpdateAcc(recipientAcc)
	return nil
}

func (account *Account) ShowAccountDetails(w io.Writer, name string, criteria, filter string) error {
	acc := account

	if name != "" {
		var err error
		acc, err = searchingAcc(name)
		if err != nil {
			return fmt.Errorf("unexcepteced error: %v\n", err)
		}
	}

	fmt.Fprintf(w, "Balance: %.2f\n", acc.Balance)
	for _, txn := range acc.Transactions {
		if filterTo(txn, criteria, filter) {
			fmt.Fprintf(w, "Time: %v, Amount: %.2f, Type: %v\n",
				txn.Time, txn.Amount, txn.Type)
		}
	}

	return nil
}

var loadAccFunc func() ([]Account, error) = LoadAcc

func searchingAcc(name string) (*Account, error) {
	accounts, err := loadAccFunc()
	if err != nil {
		return nil, err
	}

	for i, account := range accounts {
		if strings.EqualFold(account.Name, name) {
			return &accounts[i], nil
		}
		continue
	}

	return nil, errors.New("could not find account")
}

func (account *Account) addTransaction(amount float64, tt TransactionType) {
	account.Transactions = append(account.Transactions, Transactions{
		Time:   time.Now(),
		Amount: amount,
		Type:   tt,
	})
}

func filterTo(txn Transactions, criteria, filter string) bool {
	switch strings.ToLower(criteria) {
	case "type":
		return txn.Type == TransactionType(filter)

	case "amount":
		amount, err := strconv.ParseFloat(filter, 64)
		if err != nil {
			return false
		}
		return txn.Amount == amount

	case "day":
		day, err := strconv.Atoi(filter)
		if err != nil {
			return false
		}
		return txn.Time.Day() == day

	case "month":
		month, err := strconv.Atoi(filter)
		if err != nil {
			return false
		}
		return txn.Time.Month() == time.Month(month)

	case "year":
		year, err := strconv.Atoi(filter)
		if err != nil {
			return false
		}
		return txn.Time.Year() == year
	case "":
		return true
	default:
		return false
	}
}
