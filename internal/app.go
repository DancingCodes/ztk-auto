package internal

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
	appLoginURL      = "https://htkapi.lakala.com/auth/oauth/token"
	appGetInfoURL    = "https://tkapi.lakala.com/ztkapi/customer"
	selfRegCreateURL = "https://tkapi.lakala.com/htkmerchants/selfRegMp/create"
	selfRegDetailURL = "https://tkapi.lakala.com/htkmerchants/selfRegMp/detail/"
	appClientID      = "zhtk"
	appClientSecret  = "prKTcg0iN46omMAF"
)

type SelfRegCode struct {
	ID         json.Number `json:"id"`
	SelfRegSN  string      `json:"selfRegSn"`
	SelfRegURL string      `json:"selfRegUrl"`
}

func AppCreateSelfReg(ctx context.Context, token string) SelfRegCode {
	body := bytes.NewBufferString(`{"activityType":"ZTKZMTSCZC"}`)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, selfRegCreateURL, body)
	setAppHeaders(req, token)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, _ := appHTTPClient().Do(req)
	if resp == nil {
		return SelfRegCode{}
	}
	defer resp.Body.Close()

	var code SelfRegCode
	_ = json.Unmarshal(appBody(resp), &code)
	return code
}

func AppGetSelfRegDetail(ctx context.Context, token string, id json.Number) SelfRegCode {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, selfRegDetailURL+id.String(), nil)
	setAppHeaders(req, token)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := appHTTPClient().Do(req)
	if resp == nil {
		return SelfRegCode{}
	}
	defer resp.Body.Close()

	var detail SelfRegCode
	_ = json.Unmarshal(appBody(resp), &detail)
	return detail
}

func setAppHeaders(req *http.Request, token string) {
	for key, value := range map[string]string{
		"Authorization": "bearer " + token, "User-Agent": "Dart/3.4 (dart:io)", "appcode": "ZHTK",
		"Accept-Encoding": "gzip", "model": "23113RKC6C", "buildnumber": "10005",
		"ostype": "Android", "apptype": "ZHTK", "platform": "LKL", "version": "1.0.5",
		"sysversion": "{baseOS: , securityPatch: 2022-04-05, sdkInt: 32, release: 12, codename: REL, previewSdkInt: 0, incremental: 2721}",
		"request-id": appRequestID(),
	} {
		req.Header.Set(key, value)
	}
}

func AppLogin(ctx context.Context, phone, password, serviceProviderNo string) map[string]any {
	fields := [][2]string{{"username", phone}, {"password", password}, {"grant_type", "password"}, {"client_id", appClientID}, {"client_secret", appClientSecret}, {"changeAgentNo", serviceProviderNo}}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range fields {
		_ = w.WriteField(field[0], field[1])
	}
	_ = w.Close()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, appLoginURL, &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := appHTTPClient().Do(req)
	if resp == nil {
		return map[string]any{}
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	data := map[string]any{}
	_ = json.Unmarshal(raw, &data)
	return data
}

func AppFindCustomer(ctx context.Context, token, customer string) any {
	url := appGetInfoURL + "?customerNameOrExtCustomerNo=" + customer + "&status=&pageNum=1&pageSize=10"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	for key, value := range map[string]string{
		"Authorization": "bearer " + token, "User-Agent": "Dart/3.4 (dart:io)", "appcode": "ZHTK",
		"Accept-Encoding": "gzip", "model": "22041216C", "buildnumber": "10005",
		"ostype": "Android", "apptype": "ZHTK", "Content-Type": "application/json",
		"platform": "LKL", "version": "1.0.5",
		"sysversion": "{baseOS: , securityPatch: 2025-05-05, sdkInt: 35, release: 15, codename: REL, previewSdkInt: 0, incremental: 1171}",
		"request-id": appRequestID(),
	} {
		req.Header.Set(key, value)
	}
	resp, _ := appHTTPClient().Do(req)
	if resp == nil {
		return nil
	}
	defer resp.Body.Close()
	raw := appBody(resp)
	var data any
	_ = json.Unmarshal(raw, &data)
	return data
}

func appHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   25 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
}

func appBody(resp *http.Response) []byte {
	reader := io.Reader(resp.Body)
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, _ := gzip.NewReader(resp.Body)
		defer gz.Close()
		reader = gz
	}
	raw, _ := io.ReadAll(reader)
	return raw
}

func appRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6], b[8] = (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80
	s := hex.EncodeToString(b)
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}
