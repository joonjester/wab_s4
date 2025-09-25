package bank

import (
	"encoding/json"
	"fmt"
	"os"
)

const dbFile = "acc_db.json"

func AddOrUpdateAcc(newAcc *Account) {
	allAcc, err := LoadAcc()
	if err != nil {
		fmt.Println("could not load accounts:", err)
	}

	updated := false
	for index, account := range allAcc {
		if account.Id == newAcc.Id {
			allAcc[index] = *newAcc
			updated = true
			break
		}
	}

	if !updated {
		allAcc = append(allAcc, *newAcc)
	}

	data, err := json.MarshalIndent(allAcc, "", "  ")
	if err != nil {
		fmt.Println("could not save account:", err)
		return
	}

	if err := os.WriteFile(dbFile, data, 0644); err != nil {
		fmt.Println("could not write file:", err)
	}
}

func LoadAcc() ([]Account, error) {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Account{}, nil
		}
		return nil, err
	}
	var users []Account
	err = json.Unmarshal(data, &users)
	return users, err
}
