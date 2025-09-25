package bank

import (
	"bytes"
	"fmt"
	"strconv"
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
	loadAccFunc = func() ([]Account, error) {
		return []Account{
			{
				Name:        "TestTest",
				Balance:     50,
				AccountType: Giro,
				Transactions: []Transactions{
					{Time: time, Amount: 50, Type: Withdraw},
				},
			},
		}, nil
	}
	defer func() { loadAccFunc = LoadAcc }()

	show_account := map[string]struct {
		person       string
		withdraw     float64
		wantedError  bool
		wantedOutput string
	}{
		"Happy Path": {
			withdraw:     50,
			person:       "",
			wantedError:  false,
			wantedOutput: fmt.Sprintf("Balance: 50.00\nTime: %v, Amount: 50.00, Type: withdraw", time),
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
				Time:   time,
				Amount: tc.withdraw,
				Type:   Withdraw,
			}

			acc.Transactions = append(acc.Transactions, txn)
			err := acc.ShowAccountDetails(&buf, tc.person, "", "")

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

func TestFilterTransactions(t *testing.T) {
	time := time.Now()
	loadAccFunc = func() ([]Account, error) {
		return []Account{
			{
				Name:        "TestTest",
				Balance:     50,
				AccountType: Giro,
				Transactions: []Transactions{
					{Time: time, Amount: 50, Type: Withdraw},
				},
			},
		}, nil
	}
	defer func() { loadAccFunc = LoadAcc }()

	tests := map[string]struct {
		criteria    string
		filter      string
		transaction string
		want        string
		wantError   bool
	}{
		"Happy Path: Show Deposit": {
			criteria:    "type",
			filter:      "deposit",
			transaction: "deposit",
			want:        fmt.Sprintf("Balance: 100.00\nTime: %v, Amount: 50.00, Type: deposit", time),
		},
		"Happy Path: Show Withdraw": {
			criteria:    "type",
			filter:      "withdraw",
			transaction: "withdraw",
			want:        fmt.Sprintf("Balance: 100.00\nTime: %v, Amount: 50.00, Type: withdraw", time),
		},
		"Happy Path: Show Transfer": {
			criteria:    "type",
			filter:      "transfer",
			transaction: "transfer",
			want:        fmt.Sprintf("Balance: 100.00\nTime: %v, Amount: 50.00, Type: transfer", time),
		},
		"Happy Path: Show Amount": {
			criteria:    "amount",
			filter:      "50",
			transaction: "deposit",
			want:        fmt.Sprintf("Balance: 100.00\nTime: %v, Amount: 50.00, Type: deposit", time),
		}, "Happy Path: Show Day": {
			criteria:    "day",
			filter:      strconv.Itoa(time.Day()),
			transaction: "deposit",
			want:        fmt.Sprintf("Balance: 100.00\nTime: %v, Amount: 50.00, Type: deposit", time),
		}, "Happy Path: Show Month": {
			criteria:    "month",
			filter:      fmt.Sprintf("%d", time.Month()),
			transaction: "deposit",
			want:        fmt.Sprintf("Balance: 100.00\nTime: %v, Amount: 50.00, Type: deposit", time),
		}, "Happy Path: Show Year": {
			criteria:    "year",
			filter:      strconv.Itoa(time.Year()),
			transaction: "deposit",
			want:        fmt.Sprintf("Balance: 100.00\nTime: %v, Amount: 50.00, Type: deposit", time),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			acc := Account{
				Balance: 100.00,
				Transactions: []Transactions{
					{
						Time:   time,
						Amount: 50.00,
						Type:   TransactionType(tc.transaction),
					},
				},
			}

			var buf bytes.Buffer
			_ = acc.ShowAccountDetails(&buf, "", tc.criteria, tc.filter)

			got := buf.String()
			if strings.TrimSpace(tc.want) != strings.TrimSpace(got) {
				t.Errorf("got = %v, want = %v", got, tc.want)
			}

		})
	}
}

func TestAccountTransfer(t *testing.T) {
	var testAccounts []Account

	loadAccFunc = func() ([]Account, error) {
		return []Account{
			{Id: "1", Name: "Alice", Balance: 1000, AccountType: Giro},
			{Id: "2", Name: "Bob", Balance: 500, AccountType: Giro},
		}, nil
	}
	defer func() { loadAccFunc = LoadAcc }()

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
			testAccounts = []Account{
				{Id: "1", Name: "Alice", Balance: 1000, AccountType: Giro},
				{Id: "2", Name: "Bob", Balance: 500, AccountType: Giro},
			}

			originalSearchingAcc := searchingAcc
			searchingAcc := func(name string) (*Account, error) {
				for i := range testAccounts {
					if testAccounts[i].Name == name {
						return &testAccounts[i], nil
					}
				}
				return nil, fmt.Errorf("could not find account: %s", name)
			}
			defer func() { searchingAcc = originalSearchingAcc }()

			var fromAcc *Account
			for i := range testAccounts {
				if testAccounts[i].Name == tt.from {
					fromAcc = &testAccounts[i]
					break
				}
			}
			if fromAcc == nil {
				t.Fatalf("From account '%s' not found", tt.from)
			}

			err := fromAcc.Transfer(tt.amount, tt.to)

			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if fromAcc.Balance != tt.wantFromBal {
				t.Errorf("from account (%s) balance = %.2f, want %.2f", tt.from, fromAcc.Balance, tt.wantFromBal)
			}

			if tt.to != "Charlie" {
				toAcc, searchErr := searchingAcc(tt.to)
				if searchErr != nil {
					t.Errorf("Error finding to account after transfer: %v", searchErr)
				} else {
					if toAcc.Balance != tt.wantToBal {
						t.Errorf("to account (%s) balance = %.2f, want %.2f", tt.to, toAcc.Balance, tt.wantToBal)
					}
				}
			}

			if !tt.wantErr {
				t.Logf("Transfer successful: %s (%.2f) -> %s, amount: %.2f",
					tt.from, fromAcc.Balance, tt.to, tt.amount)
			} else {
				t.Logf("Transfer correctly failed: %s -> %s, amount: %.2f, error: %v",
					tt.from, tt.to, tt.amount, err)
			}
		})
	}
}

func TestTransferValidation(t *testing.T) {
	tests := []struct {
		name       string
		balance    float64
		amount     float64
		recipient  string
		shouldFail bool
	}{
		{"Valid transfer", 1000, 200, "Bob", false},
		{"Insufficient funds", 100, 200, "Bob", true},
		{"Negative amount", 1000, -50, "Bob", true},
		{"Zero amount", 1000, 0, "Bob", true},
		{"Non-existent recipient", 1000, 100, "NonExistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			loadAccFunc = func() ([]Account, error) {
				return []Account{
					{Id: "1", Name: "Alice", Balance: tt.balance, AccountType: Giro},
					{Id: "2", Name: "Bob", Balance: 500, AccountType: Giro},
				}, nil
			}
			defer func() { loadAccFunc = LoadAcc }()

			fromAcc := &Account{Id: "1", Name: "Alice", Balance: tt.balance, AccountType: Giro}
			err := fromAcc.Transfer(tt.amount, tt.recipient)

			if tt.shouldFail && err == nil {
				t.Errorf("Expected transfer to fail but it succeeded")
			} else if !tt.shouldFail && err != nil {
				t.Errorf("Expected transfer to succeed but it failed: %v", err)
			}
		})
	}
}
