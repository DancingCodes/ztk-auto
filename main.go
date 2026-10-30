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
	startedAt := time.Now()
	fmt.Println("开始时间:", startedAt.Format("2006-01-02 15:04:05"))
	defer func() {
		endedAt := time.Now()
		fmt.Println("结束时间:", endedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("用时:", endedAt.Sub(startedAt))
	}()

	// PCToken
	pc := ztk.LoadPC("config/pc.json")
	pcLogin := ztk.PCLogin(context.Background(), pc.Username, pc.Password, pc.OperateID)
	pcToken, _ := pcLogin["access_token"].(string)

	if pcToken == "" {
		result, _ := json.MarshalIndent(pcLogin, "", "  ")
		fmt.Println(string(result))
		return
	}

	// PCToken -> PC商户列表
	date := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var results []map[string]any
	for _, partner := range pc.PCPartners {
		results = append(results, map[string]any{
			"partner": partner,
			"data":    ztk.PCFindChange(context.Background(), pcToken, date, partner),
		})
	}

	for _, result := range results {
		if _, ok := result["data"].(map[string]any)["records"].([]any); !ok {
			response, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(response))
			return
		}
	}

	// PC商户列表 -> App登录 -> 商户
	appAccounts := ztk.AppAccountsByServiceProvider(ztk.LoadAppAccounts("config/app.json"))
	for _, result := range results {
		partner := result["partner"].(string)
		records := result["data"].(map[string]any)["records"].([]any)
		account, ok := appAccounts[partner]
		if !ok {
			for _, value := range records {
				value.(map[string]any)["isFirst"] = "无账号"
				fmt.Println(value)
			}
			continue
		}
		appLogin := ztk.AppLogin(context.Background(), account.Phone, account.Password, account.ServiceProviderNo)
		appToken, _ := appLogin["access_token"].(string)
		if appToken == "" {
			for _, value := range records {
				value.(map[string]any)["isFirst"] = "账密不正确"
			}
			continue
		}
		for _, value := range records {
			record := value.(map[string]any)
			externalCustomerNo := record["externalCustomerNo"].(string)
			appData := ztk.AppFindCustomer(context.Background(), appToken, externalCustomerNo)
			appResult, ok := appData.(map[string]any)
			if !ok {
				record["isFirst"] = "接口返回异常"
				continue
			}
			if _, ok := appResult["records"].([]any); !ok {
				record["isFirst"] = "无商户数据"
				continue
			}
			record["appData"] = appResult["records"].([]any)
		}
	}

	if err := utils.ExportExcel(results); err != nil {
		fmt.Println("导出 Excel 失败:", err)
		return
	}
}
