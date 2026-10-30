package internal

import (
	"context"
	"errors"
)

func GetTransactions(ctx context.Context, token, date string, partners []Partner) []map[string]any {
	var results []map[string]any
	for _, partner := range partners {
		results = append(results, map[string]any{
			"partner": partner.ServiceProviderNo,
			"data":    PCFindTransactions(ctx, token, date, partner.ServiceProviderNo),
		})
	}
	return results
}

func BuildTransactionExport(ctx context.Context, date string) ([]byte, error) {
	accounts := LoadAccounts("config/accounts.json")
	login := PCLogin(ctx, accounts.Admin.Username, accounts.Admin.Password, accounts.Admin.OperateID)
	token, _ := login["access_token"].(string)
	if token == "" {
		return nil, errors.New("PC登录失败")
	}
	results := GetTransactions(ctx, token, date, accounts.Partners)

	for _, result := range results {
		response, ok := result["data"].(map[string]any)
		if !ok {
			return nil, errors.New("交易列表获取失败：" + result["partner"].(string))
		}
		data, ok := response["data"].(map[string]any)
		if !ok {
			return nil, errors.New("交易列表获取失败：" + result["partner"].(string))
		}
		if _, ok := data["records"].([]any); !ok {
			return nil, errors.New("交易列表获取失败：" + result["partner"].(string))
		}
	}

	tranType := PCGetDictionary(ctx, token, "tranType")
	extTransCode := PCGetDictionary(ctx, token, "extTransCode")
	transCardType := PCGetDictionary(ctx, token, "transCardType")
	for _, result := range results {
		records := result["data"].(map[string]any)["data"].(map[string]any)["records"].([]any)
		for _, value := range records {
			record := value.(map[string]any)
			if name := tranType[record["type"].(string)]; name != "" {
				record["type"] = name
			}
			if name := extTransCode[record["extTransCode"].(string)]; name != "" {
				record["extTransCode"] = name
			}
			if name := transCardType[record["transCardType"].(string)]; name != "" {
				record["transCardType"] = name
			}
		}
	}

	file, err := BuildTransactionExcel(results)
	if err != nil {
		return nil, errors.New("交易列表导出失败")
	}
	return file, nil
}

func BuildMerchantExport(ctx context.Context, date string) ([]byte, error) {
	// Admin登录
	accounts := LoadAccounts("config/accounts.json")
	pcLogin := PCLogin(ctx, accounts.Admin.Username, accounts.Admin.Password, accounts.Admin.OperateID)
	pcToken, _ := pcLogin["access_token"].(string)
	if pcToken == "" {
		return nil, errors.New("PC登录失败")
	}
	customerStatus := PCGetCustomerStatus(ctx, pcToken)

	// Admin商户信息
	var results []map[string]any
	for _, partner := range accounts.Partners {
		results = append(results, map[string]any{
			"partner": partner.ServiceProviderNo,
			"account": partner,
			"data":    PCFindChange(ctx, pcToken, date, partner.ServiceProviderNo),
		})
	}

	for _, result := range results {
		data, ok := result["data"].(map[string]any)
		if !ok {
			return nil, errors.New("PC商户列表获取失败：" + result["partner"].(string))
		}
		if _, ok := data["records"].([]any); !ok {
			return nil, errors.New("PC商户列表获取失败：" + result["partner"].(string))
		}
	}

	for _, result := range results {
		records := result["data"].(map[string]any)["records"].([]any)
		for _, value := range records {
			record := value.(map[string]any)
			if name := accountType[record["agentNature"].(string)]; name != "" {
				record["agentNature"] = name
			}
			if name := merchantType[record["channelType"].(string)]; name != "" {
				record["channelType"] = name
			}
			if name := merchantSource[record["merchantSource"].(string)]; name != "" {
				record["merchantSource"] = name
			}
			if name := customerStatus[record["customerStatus"].(string)]; name != "" {
				record["customerStatus"] = name
			}
		}
	}

	// App登录 + App商户信息
	for _, result := range results {
		records := result["data"].(map[string]any)["records"].([]any)
		account := result["account"].(Partner)
		if account.Phone == "" || account.Password == "" {
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
		return nil, errors.New("商户列表导出失败")
	}
	return file, nil
}

var accountType = map[string]string{
	"PERSONAL":           "个人",
	"ZHTK_OPERATE_CENTER": "运营中心",
}

var merchantType = map[string]string{
	"TP_MERCHANT": "企业商户",
	"TP_PERSONAL": "个人商户",
}

var merchantSource = map[string]string{
	"AGENT": "服务商平台",
	"APP":   "APP平台",
	"H5":    "H5平台",
}
