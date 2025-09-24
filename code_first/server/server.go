package server

import (
	"code_first/bank"
	"encoding/json"
	"net/http"
)

var acc = &bank.Account{
	Id: "001",
	Name: "John",
	Balance: 1000.0,
	AccountType: bank.Giro,
}

type Transaction struct {
	Amount float64 `json:"amount"`
	To string `json:"to"`
}

func showAccountDetails(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	acc.ShowAccountDetails(w, req.URL.Query().Get("name"))
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


	acc.Deposit(transaction.Amount)
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

	acc.Transfer(transaction.Amount, transaction.To)
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
	acc.Withdraw(transaction.Amount)
}

func Router() {

    http.HandleFunc("/show", showAccountDetails)
    http.HandleFunc("/deposit", deposit)
    http.HandleFunc("/transfer", transfer)
    http.HandleFunc("/witdraw", withdraw)

    http.ListenAndServe(":8090", nil)
}