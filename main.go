package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ztk "ztk-auto/internal"
)

func main() {
	// PCToken
	pc := ztk.LoadPC("config/pc.json")
	pcLogin := ztk.PCLogin(context.Background(), pc.Username, pc.Password, pc.OperateID)
	pcToken, _ := pcLogin["access_token"].(string)

	// PCToken -> PC商户列表
	date := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var results []map[string]any
	for _, partner := range pc.PCPartners {
		results = append(results, map[string]any{
			"partner": partner,
			"data":    ztk.PCFindChange(context.Background(), pcToken, date, partner),
		})
	}

	// PC商户列表 -> App登录 -> 商户
	appAccounts := ztk.AppAccountsByServiceProvider(ztk.LoadAppAccounts("config/app.json"))
	for _, result := range results {
		partner := result["partner"].(string)
		account, ok := appAccounts[partner]
		if !ok || account.Phone == "" || account.Password == "" {
			continue
		}
		appLogin := ztk.AppLogin(context.Background(), account.Phone, account.Password, account.ServiceProviderNo)
		appToken, _ := appLogin["access_token"].(string)
		if appToken == "" {
			continue
		}

		pcData, ok := result["data"].(map[string]any)
		if !ok {
			continue
		}
		records, ok := pcData["records"].([]any)
		if !ok {
			continue
		}
		for _, value := range records {
			record, ok := value.(map[string]any)
			if !ok {
				continue
			}
			externalCustomerNo, _ := record["externalCustomerNo"].(string)
			if externalCustomerNo != "" {
				record["appData"] = ztk.AppFindCustomer(context.Background(), appToken, externalCustomerNo)
			}
		}
	}

	result, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(result))
}
