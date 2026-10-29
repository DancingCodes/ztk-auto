package main

import (
	"context"
	"encoding/json"
	"os"

	"ztk-auto/internal/phone"
)

func main() {
	items := phone.LoadAccounts("config/phoneAccounts.json")
	phoneAccounts := phone.AccountsByPhone(items)
	results := make(map[string]any, len(phoneAccounts))
	for phoneNo, item := range phoneAccounts {
		results[phoneNo] = phone.Login(context.Background(), item.Phone, item.Password)
	}
	_ = json.NewEncoder(os.Stdout).Encode(results)
}
