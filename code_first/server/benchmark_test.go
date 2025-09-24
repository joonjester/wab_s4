package server

import (
	"bytes"
	"code_first/bank"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestAccountForBenchmark() {
	acc = &bank.Account{
		Id:          "test123",
		Name:        "Test Account",
		Balance:     1000.0,
		AccountType: bank.Giro,
		Overdraw:    500.0,
	}
}

func BenchmarkShowAccountDetails(b *testing.B) {
	setupTestAccountForBenchmark()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/show?name=Test", nil)
		w := httptest.NewRecorder()
		showAccountDetails(w, req)
	}
}

func BenchmarkDeposit(b *testing.B) {
	setupTestAccountForBenchmark()

	transaction := Transaction{Amount: 100.0}
	jsonData, _ := json.Marshal(transaction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setupTestAccountForBenchmark()

		req := httptest.NewRequest("POST", "/deposit", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		deposit(w, req)
	}
}

func BenchmarkWithdraw(b *testing.B) {
	transaction := Transaction{Amount: 50.0}
	jsonData, _ := json.Marshal(transaction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setupTestAccountForBenchmark()
		req := httptest.NewRequest("POST", "/withdraw", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		withdraw(w, req)
	}
}

func BenchmarkTransfer(b *testing.B) {
	transaction := Transaction{Amount: 100.0, To: "recipient"}
	jsonData, _ := json.Marshal(transaction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setupTestAccountForBenchmark()
		req := httptest.NewRequest("POST", "/transfer", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		transfer(w, req)
	}
}

func BenchmarkDepositParallel(b *testing.B) {
	setupTestAccountForBenchmark()

	transaction := Transaction{Amount: 10.0}
	jsonData, _ := json.Marshal(transaction)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("POST", "/deposit", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			deposit(w, req)
		}
	})
}

func BenchmarkFullTransactionFlow(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/deposit", deposit)
	mux.HandleFunc("/withdraw", withdraw)
	mux.HandleFunc("/show", showAccountDetails)

	server := httptest.NewServer(mux)
	defer server.Close()

	transaction := Transaction{Amount: 100.0}
	jsonData, _ := json.Marshal(transaction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := http.Post(server.URL+"/deposit", "application/json", bytes.NewBuffer(jsonData))
		resp.Body.Close()

		resp, _ = http.Get(server.URL + "/show?name=Test Account")
		resp.Body.Close()

		resp, _ = http.Post(server.URL+"/withdraw", "application/json", bytes.NewBuffer(jsonData))
		resp.Body.Close()
	}
}
