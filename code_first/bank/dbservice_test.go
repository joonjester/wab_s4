package bank

import (
	"os"
	"path/filepath"
	"testing"
)

func withTempDB(t *testing.T, testFunc func()) {
	tmp := filepath.Join(t.TempDir(), "accounts.json")
	oldDB := dbFile
	dbFile = tmp
	defer func() { dbFile = oldDB }()
	testFunc()
}

func TestAddOrUpdateAcc(t *testing.T) {
	withTempDB(t, func() {
		// 1. Add a new account
		acc1 := &Account{Id: "1", Name: "Alice"}
		AddOrUpdateAcc(acc1)

		accounts, err := LoadAcc()
		if err != nil {
			t.Fatalf("unexpected error loading accounts: %v", err)
		}
		if len(accounts) != 1 {
			t.Fatalf("expected 1 account, got %d", len(accounts))
		}
		if accounts[0].Name != "Alice" {
			t.Errorf("expected name Alice, got %s", accounts[0].Name)
		}

		// 2. Add another account
		acc2 := &Account{Id: "2", Name: "Bob"}
		AddOrUpdateAcc(acc2)

		accounts, _ = LoadAcc()
		if len(accounts) != 2 {
			t.Fatalf("expected 2 accounts, got %d", len(accounts))
		}

		// 3. Update existing account (Alice → Alicia)
		acc1Updated := &Account{Id: "1", Name: "Alicia"}
		AddOrUpdateAcc(acc1Updated)

		accounts, _ = LoadAcc()
		if len(accounts) != 2 {
			t.Fatalf("expected 2 accounts after update, got %d", len(accounts))
		}

		found := false
		for _, a := range accounts {
			if a.Id == "1" {
				found = true
				if a.Name != "Alicia" {
					t.Errorf("expected updated name Alicia, got %s", a.Name)
				}
			}
		}
		if !found {
			t.Errorf("updated account not found")
		}
	})
}

func TestLoadAcc_FileNotExist(t *testing.T) {
	withTempDB(t, func() {
		// No file written yet
		accounts, err := LoadAcc()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(accounts) != 0 {
			t.Errorf("expected empty accounts, got %d", len(accounts))
		}
	})
}

func TestLoadAcc_InvalidJSON(t *testing.T) {
	withTempDB(t, func() {
		// Write invalid JSON
		if err := os.WriteFile(dbFile, []byte("{not-json}"), 0644); err != nil {
			t.Fatalf("could not write temp file: %v", err)
		}

		_, err := LoadAcc()
		if err == nil {
			t.Errorf("expected error on invalid JSON, got nil")
		}
	})
}
