package phone

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	loginURL     = "https://htkapi.lakala.com/auth/oauth/token"
	userInfoURL  = "https://zhtk.lklshuzikeji.com/auth/user/operator/userInfo/zhtk/"
	getInfoURL   = "https://tkapi.lakala.com/ztkapi/customer"
	clientID     = "zhtk"
	clientSecret = "prKTcg0iN46omMAF"
)

func Login(ctx context.Context, phone, password string) map[string]any {
	agentNo := UserNo(ctx, phone)
	fields := [][2]string{{"username", phone}, {"password", password}, {"grant_type", "password"}, {"client_id", clientID}, {"client_secret", clientSecret}, {"changeAgentNo", agentNo}}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range fields {
		_ = w.WriteField(field[0], field[1])
	}
	_ = w.Close()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := httpClient().Do(req)
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	data := map[string]any{}
	_ = json.Unmarshal(raw, &data)
	return data
}

func UserNo(ctx context.Context, phone string) string {
	url := userInfoURL + phone + "?appCode=ZHTK&loginAccount=" + phone
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, _ := httpClient().Do(req)
	defer resp.Body.Close()
	reader := io.Reader(resp.Body)
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, _ := gzip.NewReader(resp.Body)
		defer gz.Close()
		reader = gz
	}
	raw, _ := io.ReadAll(reader)
	var data []map[string]any
	_ = json.Unmarshal(raw, &data)
	if len(data) == 0 {
		return ""
	}
	userNo, _ := data[0]["userNo"].(string)
	return userNo
}

func GetInfo(ctx context.Context, token, customer string) any {
	url := getInfoURL + "?customerNameOrExtCustomerNo=" + customer + "&status=&pageNum=1&pageSize=10"
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString("{}"))
	headers := phoneHeaders()
	headers["Authorization"] = "bearer " + token
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, _ := httpClient().Do(req)
	defer resp.Body.Close()
	reader := io.Reader(resp.Body)
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, _ := gzip.NewReader(resp.Body)
		defer gz.Close()
		reader = gz
	}
	raw, _ := io.ReadAll(reader)
	var data any
	_ = json.Unmarshal(raw, &data)
	return data
}

func phoneHeaders() map[string]string {
	return map[string]string{
		"User-Agent": "Dart/3.4 (dart:io)", "appcode": "ZHTK", "Accept-Encoding": "gzip",
		"model": "22041216C", "buildnumber": "10005", "ostype": "Android",
		"apptype": "ZHTK", "Content-Type": "application/json", "platform": "LKL", "version": "1.0.5",
		"sysversion": "{baseOS: , securityPatch: 2025-05-05, sdkInt: 35, release: 15, codename: REL, previewSdkInt: 0, incremental: 1171}",
		"request-id": requestID(),
	}
}

func httpClient() *http.Client {
	return &http.Client{
		Timeout: 25 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
}

func requestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6], b[8] = (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80
	s := hex.EncodeToString(b)
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}
