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
)

type PCAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoadPCAccount(path string) PCAccount {
	raw, _ := os.ReadFile(path)
	var account PCAccount
	_ = json.Unmarshal(raw, &account)
	return account
}

func LoadPCPartners(path string) []string {
	raw, _ := os.ReadFile(path)
	var data struct {
		PCPartners []string `json:"pcPartners"`
	}
	_ = json.Unmarshal(raw, &data)
	partners := make([]string, 0, len(data.PCPartners))
	for _, item := range data.PCPartners {
		if value := strings.TrimSpace(item); value != "" {
			partners = append(partners, value)
		}
	}
	return partners
}

func PCLogin(ctx context.Context, username, password string) map[string]any {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range [][2]string{{"username", username}, {"password", password}, {"grant_type", "password"}} {
		_ = w.WriteField(field[0], field[1])
	}
	_ = w.Close()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, pcLoginURL, &body)
	req.Header.Set("Authorization", pcBasicAuth)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := pcHTTPClient().Do(req)
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
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var data any
	_ = json.Unmarshal(raw, &data)
	return data
}

func pcHTTPClient() *http.Client { return &http.Client{Timeout: 25 * time.Second} }
