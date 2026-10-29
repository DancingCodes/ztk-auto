package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ztk "ztk-auto/internal"
	"ztk-auto/utils"
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
		account := appAccounts[partner]
		appLogin := ztk.AppLogin(context.Background(), account.Phone, account.Password, account.ServiceProviderNo)
		appToken, _ := appLogin["access_token"].(string)
		records := result["data"].(map[string]any)["records"].([]any)
		for _, value := range records {
			record := value.(map[string]any)
			externalCustomerNo, _ := record["externalCustomerNo"].(string)
			record["appData"] = ztk.AppFindCustomer(context.Background(), appToken, externalCustomerNo)
		}
	}

	// 打印数据 / 到处execl
	result, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(result))
	utils.ExportExcel(results)
}
