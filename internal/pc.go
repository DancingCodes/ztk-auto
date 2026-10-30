package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const (
	pcLoginURL  = "https://tkapi.lakala.com/auth/oauth/token"
	pcBasicAuth = "Basic bHN5OlNSNjJRc1lLaTFsMDFJZUI="
	pcChangeURL = "https://tkapi.lakala.com/ztkapi/customer/findChange"
	pcDictionaryURL = "https://tkapi.lakala.com/ztkapi/dictionary?dictTypeId="
	pcTransURL  = "https://tkapi.lakala.com/ztkapi/trans/list"
)

type Admin struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	OperateID string `json:"operateId"`
}

type Accounts struct {
	Admin    Admin     `json:"admin"`
	Partners []Partner `json:"partners"`
}

type Partner struct {
	ServiceProviderNo string `json:"serviceProviderNo"`
	Phone             string `json:"phone"`
	Password          string `json:"password"`
}

func LoadAccounts(path string) Accounts {
	raw, _ := os.ReadFile(path)
	var config Accounts
	_ = json.Unmarshal(raw, &config)
	return config
}

func PCLogin(ctx context.Context, username, password, operateID string) map[string]any {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range [][2]string{{"username", username}, {"password", password}, {"operateId", operateID}, {"grant_type", "password"}} {
		_ = w.WriteField(field[0], field[1])
	}
	_ = w.Close()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, pcLoginURL, &body)
	req.Header.Set("Authorization", pcBasicAuth)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := pcHTTPClient().Do(req)
	if resp == nil {
		return map[string]any{}
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	data := map[string]any{}
	_ = json.Unmarshal(raw, &data)
	return data
}

func PCFindChange(ctx context.Context, token, date, agentNo string) any {
	url := pcChangeURL + "?startTime=" + date + "&endTime=" + date + "&agentNo=" + agentNo + "&pageNum=1&pageSize=999"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "bearer "+token)
	resp, _ := pcHTTPClient().Do(req)
	if resp == nil {
		return map[string]any{}
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var data any
	_ = json.Unmarshal(raw, &data)
	return data
}

func PCFindTransactions(ctx context.Context, token, date, agentNo string) any {
	url := pcTransURL + "?transStartTime=" + date + "+00:00:00&transEndTime=" + date + "+23:59:59&agentNo=" + agentNo + "&pageNum=1&pageSize=999"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "bearer "+token)
	resp, _ := pcHTTPClient().Do(req)
	if resp == nil {
		return map[string]any{}
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var data any
	_ = json.Unmarshal(raw, &data)
	return data
}

func PCGetDictionary(ctx context.Context, token, dictTypeID string) map[string]string {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, pcDictionaryURL+dictTypeID, nil)
	req.Header.Set("Authorization", "bearer "+token)
	resp, _ := pcHTTPClient().Do(req)
	if resp == nil {
		return map[string]string{}
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var items []map[string]any
	_ = json.Unmarshal(raw, &items)
	status := map[string]string{}
	for _, item := range items {
		status[item["dictId"].(string)] = item["dictName"].(string)
	}
	return status
}

func PCGetCustomerStatus(ctx context.Context, token string) map[string]string {
	return PCGetDictionary(ctx, token, "customerStatus")
}

func pcHTTPClient() *http.Client { return &http.Client{Timeout: 25 * time.Second} }
