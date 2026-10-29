// Package account 读取本地多账号配置。
package account

import (
	"encoding/json"
	"fmt"
	"os"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	AgentNo  string `json:"agent_no"`
}

type file struct {
	Accounts []Account `json:"accounts"`
}

// Load 从 JSON 文件读取账号，并拒绝缺少必要凭据的记录。
func Load(path string) ([]Account, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg file
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parse accounts: %w", err)
	}
	if len(cfg.Accounts) == 0 {
		return nil, fmt.Errorf("accounts is empty")
	}
	for i, a := range cfg.Accounts {
		if a.Username == "" || a.Password == "" {
			return nil, fmt.Errorf("account #%d requires username and password", i+1)
		}
	}
	return cfg.Accounts, nil
}

// ByUsername 为批处理按手机号查找登录凭据。
func ByUsername(accounts []Account) (map[string]Account, error) {
	index := make(map[string]Account, len(accounts))
	for _, a := range accounts {
		if _, exists := index[a.Username]; exists {
			return nil, fmt.Errorf("duplicate username: %s", a.Username)
		}
		index[a.Username] = a
	}
	return index, nil
}
