package internal

import (
	"context"
)

func BuildMerchantExport(ctx context.Context, date string) ([]byte, error) {
	// Pc登录
	pc := LoadPC("config/pc.json")
	pcLogin := PCLogin(ctx, pc.Username, pc.Password, pc.OperateID)
	pcToken, _ := pcLogin["access_token"].(string)
	if pcToken == "" {
		return BuildMessageExcel("PC登录失败")
	}

	// Pc商户信息
	var results []map[string]any
	for _, partner := range pc.PCPartners {
		results = append(results, map[string]any{
			"partner": partner,
			"data":    PCFindChange(ctx, pcToken, date, partner),
		})
	}

	for _, result := range results {
		data, ok := result["data"].(map[string]any)
		if !ok {
			return BuildMessageExcel("PC商户列表获取失败：" + result["partner"].(string))
		}
		if _, ok := data["records"].([]any); !ok {
			return BuildMessageExcel("PC商户列表获取失败：" + result["partner"].(string))
		}
	}

	// App登录 + App商户信息
	appAccounts := AppAccountsByServiceProvider(LoadAppAccounts("config/app.json"))
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
		appLogin := AppLogin(ctx, account.Phone, account.Password, account.ServiceProviderNo)
		appToken, _ := appLogin["access_token"].(string)
		if appToken == "" {
			for _, value := range records {
				value.(map[string]any)["isFirst"] = "App登录失败"
			}
			continue
		}
		for _, value := range records {
			record := value.(map[string]any)
			externalCustomerNo := record["externalCustomerNo"].(string)
			appData := AppFindCustomer(ctx, appToken, externalCustomerNo)
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

	// 导出
	file, err := BuildExcel(results)
	if err != nil {
		return BuildMessageExcel("商户列表导出失败")
	}
	return file, nil
}
