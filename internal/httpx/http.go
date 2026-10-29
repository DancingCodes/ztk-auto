// Package httpx 收口网络请求的兜底逻辑。
package httpx

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type Field struct{ Key, Value string }

func Multipart(fields []Field) (*bytes.Buffer, string) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for _, field := range fields {
		_ = w.WriteField(field.Key, field.Value)
	}
	_ = w.Close()
	return &body, w.FormDataContentType()
}

// JSON 负责创建请求、执行、解压 gzip、检查状态码并解析 JSON。
func JSON[T any](ctx context.Context, client *http.Client, method, url string, body io.Reader, headers map[string]string) (T, error) {
	var data T
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return data, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	raw, err := readBody(resp)
	if err != nil {
		return data, err
	}
	if resp.StatusCode/100 != 2 {
		return data, fmt.Errorf("HTTP %s: %s", resp.Status, raw)
	}
	return data, json.Unmarshal(raw, &data)
}

func readBody(resp *http.Response) ([]byte, error) {
	if !strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		return io.ReadAll(resp.Body)
	}
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	return io.ReadAll(gz)
}
