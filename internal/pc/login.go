package pc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	loginURL  = "https://tkapi.lakala.com/auth/oauth/token"
	basicAuth = "Basic bHN5OlNSNjJRc1lLaTFsMDFJZUI="
	changeURL = "https://tkapi.lakala.com/ztkapi/customer/findChange"
)

func Login(ctx context.Context, username, password, operateID string) map[string]any {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range [][2]string{{"username", username}, {"password", password}, {"operateId", operateID}, {"grant_type", "password"}} {
		_ = w.WriteField(field[0], field[1])
	}
	_ = w.Close()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, &body)
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := httpClient().Do(req)
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	data := map[string]any{}
	_ = json.Unmarshal(raw, &data)
	return data
}

func FindChange(ctx context.Context, token, date, agentNo string) any {
	url := changeURL + "?startTime=" + date + "&endTime=" + date + "&agentNo=" + agentNo + "&pageNum=1&pageSize=999"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "bearer "+token)
	resp, _ := httpClient().Do(req)
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var data any
	_ = json.Unmarshal(raw, &data)
	return data
}

func httpClient() *http.Client { return &http.Client{Timeout: 25 * time.Second} }
