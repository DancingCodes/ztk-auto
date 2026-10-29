// Package auth 封装认证接口及认证相关配置。
package auth

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultLoginURL     = "https://htkapi.lakala.com/auth/oauth/token"
	defaultGetInfoURL   = "https://tkapi.lakala.com/ztkapi/customer?customerNameOrExtCustomerNo=123&status=&pageNum=1&pageSize=10"
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
	u, err := url.Parse(c.getInfoURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("customerNameOrExtCustomerNo", customerNameOrExtCustomerNo)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString("{}"))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+token)
	for key, value := range map[string]string{
		"User-Agent":       "Dart/3.4 (dart:io)",
		"appcode":          "ZHTK",
		"Accept-Encoding":  "gzip",
		"model":            "22041216C",
		"buildnumber":      "10005",
		"ostype":           "Android",
		"apptype":          "ZHTK",
		"Content-Type":     "application/json",
		"platform":         "LKL",
		"version":          "1.0.5",
		"sysversion":       "{baseOS: , securityPatch: 2025-05-05, sdkInt: 35, release: 15, codename: REL, previewSdkInt: 0, incremental: 1171}",
		"request-id":       requestID(),
	} {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("getinfo: HTTP %s: %s", resp.Status, raw)
	}

	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("getinfo: invalid JSON: %w", err)
	}
	return data, nil
}

func requestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	s := hex.EncodeToString(b)
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	if resp.Header.Get("Content-Encoding") != "gzip" {
		return io.ReadAll(resp.Body)
	}
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	return io.ReadAll(gz)
}

// Login 以密码模式登录。agent 为空时不发送 changeAgentNo。
func (c *Client) Login(ctx context.Context, user, pass, agent string) (map[string]any, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range []struct{ key, value string }{
		{"username", user},
		{"password", pass},
		{"grant_type", "password"},
		{"client_id", c.clientID},
		{"client_secret", c.clientSecret},
	} {
		if err := w.WriteField(field.key, field.value); err != nil {
			return nil, err
		}
	}
	if agent != "" {
		if err := w.WriteField("changeAgentNo", agent); err != nil {
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
		return nil, fmt.Errorf("login: HTTP %s: %s", resp.Status, raw)
	}

	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("login: invalid JSON: %w", err)
	}
	if _, ok := data["access_token"]; !ok {
		return nil, fmt.Errorf("login failed: %s", raw)
	}
	return data, nil
}
