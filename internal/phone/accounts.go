package phone

import (
	"encoding/json"
	"os"
)

type Account struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func LoadAccounts(path string) []Account {
	raw, _ := os.ReadFile(path)
	var data struct {
		PhoneAccounts []Account `json:"phoneAccounts"`
	}
	_ = json.Unmarshal(raw, &data)
	return data.PhoneAccounts
}

func AccountsByPhone(accounts []Account) map[string]Account {
	items := make(map[string]Account, len(accounts))
	for _, item := range accounts {
		items[item.Phone] = item
	}
	return items
}
