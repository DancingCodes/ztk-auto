package main

import (
	"context"
	"time"

	ztk "ztk-auto/internal"
	"ztk-auto/utils"
)

func main() {
	// PCToken
	pc := ztk.LoadPC("config/pc.json")
	pcLogin := ztk.PCLogin(context.Background(), pc.Username, pc.Password, pc.OperateID)
	pcToken, _ := pcLogin["access_token"].(string)

	if pcToken == "" {
		_ = utils.ExportMessageExcel("PC登录失败")
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
			_ = utils.ExportMessageExcel("PC商户列表获取失败：" + result["partner"].(string))
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
			appRecords, ok := appResult["records"].([]any)
			if !ok || len(appRecords) == 0 {
				record["isFirst"] = "无商户数据"
				continue
			}
			record["appData"] = appRecords
		}
	}

	if utils.ExportExcel(results) != nil {
		err := utils.ExportMessageExcel("商户列表导出失败")
		if err != nil {
			return
		}
	}
}
