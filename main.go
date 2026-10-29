package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"ztk-auto/internal/account"
	"ztk-auto/internal/auth"
)

func main() {
	items, err := account.Load("config/accounts.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "读取 config/accounts.json 失败：", err)
		fmt.Fprintln(os.Stderr, "请复制 config/accounts.example.json 为 config/accounts.json 后填写账号。")
		os.Exit(2)
	}

	accounts, err := account.ByUsername(items)
	if err != nil {
		fmt.Fprintln(os.Stderr, "账号配置错误：", err)
		os.Exit(1)
	}
	client := auth.NewClient()
	results := make(map[string]any, len(accounts))
	failed := false
	for username, item := range accounts {
		data, err := client.Login(context.Background(), item.Username, item.Password, item.AgentNo)
		if err != nil {
			results[username] = map[string]string{"error": err.Error()}
			failed = true
			continue
		}
		results[username] = data
	}
	if err := json.NewEncoder(os.Stdout).Encode(results); err != nil {
		fmt.Fprintln(os.Stderr, "输出结果失败：", err)
		os.Exit(1)
	}
	if failed {
		os.Exit(1)
	}
}
