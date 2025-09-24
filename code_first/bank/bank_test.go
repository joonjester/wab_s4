package bank

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDeposit(t *testing.T) {
	deposit_test := map[string]struct {
		amount      float64
		wantBalance float64
		wantError   bool
	}{
		"Happy Path": {
			amount:      1000,
			wantBalance: 1000,
			wantError:   false,
		},
		"Unhappy Path: deposit 0": {
			amount:      0,
			wantBalance: 0,
			wantError:   true,
		},
		"Unhappy Path: deposit < 0": {
			amount:      -100,
			wantBalance: 0,
			wantError:   true,
		},
	}

	for name, tc := range deposit_test {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			acc := &Account{}
			err := acc.Deposit(tc.amount)

			if tc.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v\n", err)
				}
				if acc.Balance != tc.wantBalance {
					t.Errorf("balance = %v, want %v", acc.Balance, tc.wantBalance)
				}
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	withdraw_test := map[string]struct {
		amount      float64
		overdraw    float64
		accountType AccountType
		wantBalance float64
		wantError   bool
	}{
		"Happy Path": {
			amount:      50,
			overdraw:    100,
			accountType: Giro,
			wantBalance: 500,
			wantError:   false,
		},
		"Happy Path: Overdrawed": {
			amount:      650,
			overdraw:    100,
			accountType: Giro,
			wantBalance: -100,
			wantError:   false,
		},
		"Unhappy Path: withdraw 0": {
			amount:      0,
			overdraw:    100,
			accountType: Giro,
			wantBalance: 500,
			wantError:   true,
		},
		"Unhappy Path: withdraw negative amount": {
			amount:      -50,
			overdraw:    100,
			accountType: Giro,
			wantBalance: 500,
			wantError:   true,
		},
		"Unhappy Path: insufficiend fund for giro": {
			amount:      800,
			overdraw:    100,
			accountType: Giro,
			wantBalance: 0,
			wantError:   true,
		},
		"Unhappy Path: insufficiend fund for saving": {
			amount:      560,
			overdraw:    0,
			accountType: Savings,
			wantBalance: 0,
			wantError:   true,
		},
	}

	for name, tc := range withdraw_test {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			acc := &Account{
				Id:          "TestTest",
				Name:        "TestTest",
				Balance:     550,
				Overdraw:    tc.overdraw,
				AccountType: tc.accountType,
			}
			err := acc.Withdraw(tc.amount)

			if tc.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v\n", err)
				}
				if acc.Balance != tc.wantBalance {
					t.Errorf("balance = %v, want = %v\n", acc.Balance, tc.wantBalance)
				}
			}
		})
	}
}

func TestShowAccount(t *testing.T) {
	time := time.Now()
	show_account := map[string]struct {
		person       string
		withdraw     float64
		wantedError  bool
		wantedOutput string
	}{
		"Happy Path": {
			withdraw:     50,
			person:       "Alice",
			wantedError:  false,
			wantedOutput: "Balance: 1000.00\n",
		},
		"Happy Path: Own Account": {
			withdraw:     50,
			person:       "",
			wantedError:  false,
			wantedOutput: fmt.Sprintf("Balance: 50.00\nTime: %v, Purpose: Get Money, Amount: 50.00, Type: withdraw", time),
		},
		"Unhappy Path: No account": {
			withdraw:     50,
			person:       "THISDOESNOTEXITS",
			wantedError:  true,
			wantedOutput: "",
		},
	}

	for name, tc := range show_account {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer

			acc := Account{
				Id:          "TestTest",
				Name:        "TestTest",
				Balance:     50,
				AccountType: Giro,
			}

			txn := Transactions{
				Time:    time,
				Purpose: "Get Money",
				Amount:  tc.withdraw,
				Type:    Withdraw,
			}

			acc.Transactions = append(acc.Transactions, txn)
			err := acc.ShowAccountDetails(&buf, tc.person)

			got := buf.String()
			if tc.wantedError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if strings.TrimSpace(got) != strings.TrimSpace(tc.wantedOutput) {
					t.Errorf("got = %s, wanted = %s\n", got, tc.wantedOutput)
				}
			}

		})
	}
}

func TestAccount_Transfer(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		amount      float64
		wantErr     bool
		wantFromBal float64
		wantToBal   float64
	}{
		{
			name:        "Happy Path",
			from:        "Alice",
			to:          "Bob",
			amount:      200,
			wantErr:     false,
			wantFromBal: 800,
			wantToBal:   700,
		},
		{
			name:        "Unhappy Path: Insufficient funds",
			from:        "Alice",
			to:          "Bob",
			amount:      2000,
			wantErr:     true,
			wantFromBal: 1000,
			wantToBal:   500,
		},
		{
			name:        "Unhappy Path: Negative amount",
			from:        "Alice",
			to:          "Bob",
			amount:      -50,
			wantErr:     true,
			wantFromBal: 1000,
			wantToBal:   500,
		},
		{
			name:        "Unhappy Path: Recipient does not exist",
			from:        "Alice",
			to:          "Charlie",
			amount:      50,
			wantErr:     true,
			wantFromBal: 1000,
			wantToBal:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fromAcc := searchingAcc(tt.from)
			toAccBefore := searchingAcc(tt.to)
			fromBalBefore := fromAcc.Balance
			toBalBefore := 0.0
			if toAccBefore != nil {
				toBalBefore = toAccBefore.Balance
			}

			err := fromAcc.Transfer(tt.amount, tt.to)

			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}

			if fromAcc.Balance != tt.wantFromBal {
				t.Errorf("from account balance = %v, want %v", fromAcc.Balance, tt.wantFromBal)
			}

			if toAcc := searchingAcc(tt.to); toAcc != nil {
				if toAcc.Balance != tt.wantToBal {
					t.Errorf("to account balance = %v, want %v", toAcc.Balance, tt.wantToBal)
				}
			} else if !tt.wantErr {
				t.Errorf("to account not found")
			}

			fromAcc.Balance = fromBalBefore
			if toAccBefore != nil {
				toAccBefore.Balance = toBalBefore
			}
		})
	}
}
