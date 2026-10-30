package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	pcLoginURL  = "https://tkapi.lakala.com/auth/oauth/token"
	pcBasicAuth = "Basic bHN5OlNSNjJRc1lLaTFsMDFJZUI="
	pcChangeURL = "https://tkapi.lakala.com/ztkapi/customer/findChange"
	pcStatusURL = "https://tkapi.lakala.com/ztkapi/dictionary?dictTypeId=customerStatus"
)

type PCAccount struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	OperateID string `json:"operateId"`
}

type PCConfig struct {
	PCAccount
	PCPartners []string `json:"pcPartners"`
}

func LoadPC(path string) PCConfig {
	raw, _ := os.ReadFile(path)
	var config PCConfig
	_ = json.Unmarshal(raw, &config)
	partners := make([]string, 0, len(config.PCPartners))
	for _, item := range config.PCPartners {
		if value := strings.TrimSpace(item); value != "" {
			partners = append(partners, value)
		}
	}
	config.PCPartners = partners
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

func PCGetCustomerStatus(ctx context.Context, token string) map[string]string {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, pcStatusURL, nil)
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

func pcHTTPClient() *http.Client { return &http.Client{Timeout: 25 * time.Second} }
