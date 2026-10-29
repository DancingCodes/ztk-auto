package main

import (
	"context"
	"encoding/json"
	"os"

	ztk "ztk-auto/internal"
)

func main() {
	items := ztk.LoadPhoneAccounts("config/phoneAccounts.json")
	phoneAccounts := ztk.PhoneAccountsByPhone(items)
	results := make(map[string]any, len(phoneAccounts))
	for phoneNo, item := range phoneAccounts {
		results[phoneNo] = ztk.AppLogin(context.Background(), item.Phone, item.Password)
	}
	_ = json.NewEncoder(os.Stdout).Encode(results)
}
