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
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < b.N; i++ {
		account.Transfer(100.0, "recipient")
	}
}
