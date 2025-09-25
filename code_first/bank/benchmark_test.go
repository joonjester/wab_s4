package bank

import (
	"testing"
)

func BenchmarkDeposit(b *testing.B) {
	account := &Account{
		Id:          "123",
		Name:        "Test",
		Balance:     1000.0,
		AccountType: Giro,
	}

	b.ResetTimer()
	for b.Loop() {
		account.Deposit(100.0)
	}
}

func BenchmarkWithdraw(b *testing.B) {
	account := &Account{
		Id:          "123",
		Name:        "Test",
		Balance:     10000.0,
		AccountType: Giro,
		Overdraw:    500.0,
	}

	b.ResetTimer()
	for b.Loop() {
		account.Withdraw(50.0)
	}
}

func BenchmarkTransfer(b *testing.B) {
	account := &Account{
		Id:          "123",
		Name:        "Test",
		Balance:     10000.0,
		AccountType: Giro,
	}

	b.ResetTimer()
	for b.Loop() {
		account.Transfer(100.0, "recipient")
	}
}

func BenchmarkConversion(b *testing.B) {

	b.ResetTimer()
	for b.Loop() {
		ConvertCurrency(10000.0, EUR, USD)
	}
}
