package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	defaultAdminLoginURL  = "https://tkapi.lakala.com/auth/oauth/token"
	defaultAdminBasicAuth = "Basic bHN5OlNSNjJRc1lLaTFsMDFJZUI="
)

// AdminClient 负责 tkapi 管理端认证请求。
type AdminClient struct {
	loginURL   string
	basicAuth  string
	httpClient *http.Client
}

// NewAdminClient 返回使用默认网页端客户端凭据的管理端认证客户端。
func NewAdminClient() *AdminClient {
	return &AdminClient{
		loginURL:   defaultAdminLoginURL,
		basicAuth:  defaultAdminBasicAuth,
		httpClient: &http.Client{Timeout: 25 * time.Second},
	}
}

// Login 以管理端账号、密码和 operateId 获取 access token。
func (c *AdminClient) Login(ctx context.Context, username, password, operateID string) (map[string]any, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range []struct{ key, value string }{
		{"username", username},
		{"password", password},
		{"operateId", operateID},
		{"grant_type", "password"},
	} {
		if err := w.WriteField(field.key, field.value); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.loginURL, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.basicAuth)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("admin login: HTTP %s: %s", resp.Status, raw)
	}

	var result map[string]any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("admin login: invalid JSON: %w", err)
	}
	if _, ok := result["access_token"]; !ok {
		return nil, fmt.Errorf("admin login: 未返回 access_token: %s", raw)
	}
	return result, nil
}
