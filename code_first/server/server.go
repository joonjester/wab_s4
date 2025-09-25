package server

import (
	"code_first/bank"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var acc *bank.Account

type Transaction struct {
	Amount         float64       `json:"amount"`
	To             string        `json:"to"`
	BaseCurrency   bank.Currency `json:"base"`
	TargetCurrency bank.Currency `json:"target"`
}

func showAccountDetails(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	err := acc.ShowAccountDetails(w, req.URL.Query().Get("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func deposit(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var transaction Transaction

	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	err = acc.Deposit(transaction.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func transfer(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var transaction Transaction

	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	err = acc.Transfer(transaction.Amount, transaction.To)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func withdraw(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var transaction Transaction

	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	err = acc.Withdraw(transaction.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func convert(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var transaction Transaction

	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	converted, err := bank.ConvertCurrency(transaction.Amount, transaction.BaseCurrency, transaction.TargetCurrency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = acc.Withdraw(*converted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func InitializeAcc(args []string) error {
	var accType bank.AccountType
	var argsLenght int
	switch strings.ToLower(args[1]) {
	case "giro":
		accType = bank.Giro
		argsLenght = 6
	case "saving":
		accType = bank.Savings
		argsLenght = 5
	default:
		return errors.New("give a valid account type: (Giro | Saving)")
	}

	if len(args) < argsLenght {
		return errors.New("Please passe Id, Name, Balance and Account Type")
	}

	balance, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return errors.New("Please give valid balance number")
	}

	overdraw, err := strconv.ParseFloat(args[5], 64)
	if err != nil && accType == bank.Giro {
		return errors.New("Please give valid overdraw value")
	}

	acc = &bank.Account{
		Id:          args[2],
		Name:        args[3],
		Balance:     balance,
		AccountType: accType,
		Overdraw:    overdraw,
	}

	return nil
}

func Router() {

	http.HandleFunc("/show", showAccountDetails)
	http.HandleFunc("/deposit", deposit)
	http.HandleFunc("/transfer", transfer)
	http.HandleFunc("/withdraw", withdraw)
	http.HandleFunc("/convert", convert)

	http.ListenAndServe(":8090", nil)
}
