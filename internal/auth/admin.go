package auth

import (
	"context"
	"net/http"
	"time"

	"ztk-auto/internal/httpx"
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
	body, contentType := httpx.Multipart([]httpx.Field{
		{"username", username}, {"password", password},
		{"operateId", operateID}, {"grant_type", "password"},
	})
	return httpx.JSON[map[string]any](ctx, c.httpClient, http.MethodPost, c.loginURL, body, map[string]string{
		"Authorization": c.basicAuth, "Content-Type": contentType,
	})
}
