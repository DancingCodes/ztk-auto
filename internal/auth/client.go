// Package auth 封装认证接口及认证相关配置。
package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"net/http"
	"time"

	"ztk-auto/internal/httpx"
)

const (
	defaultLoginURL     = "https://htkapi.lakala.com/auth/oauth/token"
	defaultGetInfoURL   = "https://tkapi.lakala.com/ztkapi/customer"
	defaultClientID     = "zhtk"
	defaultClientSecret = "prKTcg0iN46omMAF"
)

// Client 负责登录、刷新 token 等认证请求。
type Client struct {
	loginURL     string
	getInfoURL   string
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

// NewClient 返回使用默认移动端 OAuth 凭据的认证客户端。
func NewClient() *Client {
	return &Client{
		loginURL:     defaultLoginURL,
		getInfoURL:   defaultGetInfoURL,
		clientID:     defaultClientID,
		clientSecret: defaultClientSecret,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
			// 服务端证书域名不匹配，保持与原客户端兼容。
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		},
	}
}

// GetInfo 使用移动端 access token 按手机号或客户编号查询客户信息。
func (c *Client) GetInfo(ctx context.Context, token, customerNameOrExtCustomerNo string) (any, error) {
	url := c.getInfoURL + "?customerNameOrExtCustomerNo=" + customerNameOrExtCustomerNo + "&status=&pageNum=1&pageSize=10"
	return httpx.JSON[any](ctx, c.httpClient, http.MethodPost, url, bytes.NewBufferString("{}"), map[string]string{
		"Authorization":   "bearer " + token,
		"User-Agent":      "Dart/3.4 (dart:io)",
		"appcode":         "ZHTK",
		"Accept-Encoding": "gzip",
		"model":           "22041216C",
		"buildnumber":     "10005",
		"ostype":          "Android",
		"apptype":         "ZHTK",
		"Content-Type":    "application/json",
		"platform":        "LKL",
		"version":         "1.0.5",
		"sysversion":      "{baseOS: , securityPatch: 2025-05-05, sdkInt: 35, release: 15, codename: REL, previewSdkInt: 0, incremental: 1171}",
		"request-id":      requestID(),
	})
}

func requestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	s := hex.EncodeToString(b)
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

// Login 以密码模式登录。agent 为空时不发送 changeAgentNo。
func (c *Client) Login(ctx context.Context, user, pass, agent string) (map[string]any, error) {
	fields := []httpx.Field{{"username", user}, {"password", pass}, {"grant_type", "password"}, {"client_id", c.clientID}, {"client_secret", c.clientSecret}}
	if agent != "" {
		fields = append(fields, httpx.Field{"changeAgentNo", agent})
	}
	body, contentType := httpx.Multipart(fields)
	return httpx.JSON[map[string]any](ctx, c.httpClient, http.MethodPost, c.loginURL, body, map[string]string{"Content-Type": contentType})
}
