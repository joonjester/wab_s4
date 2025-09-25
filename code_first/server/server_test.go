package server

import (
	"bytes"
	"code_first/bank"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestAccount() {
	acc = &bank.Account{
		Id:          "123",
		Name:        "Alice",
		Balance:     100.0,
		AccountType: bank.Giro,
		Overdraw:    50.0,
	}
}

func TestShowAccountDetails(t *testing.T) {
	setupTestAccount()

	tests := []struct {
		name       string
		method     string
		queryParam string
		wantCode   int
	}{
		{"valid GET", http.MethodGet, "Alice", http.StatusOK},
		{"invalid method", http.MethodPost, "Alice", http.StatusMethodNotAllowed},
		{"unknown account", http.MethodGet, "DoesNotExist", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/show?name="+tt.queryParam, nil)
			rr := httptest.NewRecorder()

			showAccountDetails(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("got %d, want %d", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestDeposit(t *testing.T) {
	setupTestAccount()

	tests := []struct {
		name     string
		method   string
		body     interface{}
		wantCode int
	}{
		{"valid deposit", http.MethodPost, Transaction{Amount: 50}, http.StatusOK},
		{"invalid method", http.MethodGet, Transaction{Amount: 50}, http.StatusMethodNotAllowed},
		{"invalid json", http.MethodPost, "{bad json}", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/deposit", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			deposit(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("got %d, want %d", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	setupTestAccount()

	tests := []struct {
		name     string
		method   string
		body     interface{}
		wantCode int
	}{
		{"valid withdraw", http.MethodPost, Transaction{Amount: 30}, http.StatusOK},
		{"overdraw attempt", http.MethodPost, Transaction{Amount: 1000}, http.StatusBadRequest},
		{"invalid method", http.MethodGet, Transaction{Amount: 20}, http.StatusMethodNotAllowed},
		{"invalid json", http.MethodPost, "{bad json}", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/withdraw", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			withdraw(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("got %d, want %d", rr.Code, tt.wantCode)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	setupTestAccount()

	tests := []struct {
		name     string
		method   string
		body     any
		wantCode int
	}{
		{"invalid method", http.MethodGet, Transaction{Amount: 20, To: "Bob"}, http.StatusMethodNotAllowed},
		{"invalid json", http.MethodPost, "{bad json}", http.StatusBadRequest},
		{"insufficient funds", http.MethodPost, Transaction{Amount: 2000, To: "Bob"}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/transfer", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			transfer(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("got %d, want %d", rr.Code, tt.wantCode)
			}
		})
	}
}
