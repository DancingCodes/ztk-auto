// Package admin 封装 tkapi 管理端业务接口。
package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const findChangeURL = "https://tkapi.lakala.com/ztkapi/customer/findChange"

// Client 负责使用 admin token 调用管理端业务接口。
type Client struct {
	httpClient *http.Client
}

// NewClient 返回默认的管理端业务客户端。
func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: 25 * time.Second}}
}

// FindChange 查询某合作方在 date（YYYY-MM-DD）当天的全部变更客户记录。
// 请求固定 pageNum=1、pageSize=999。
func (c *Client) FindChange(ctx context.Context, token, date, agentNo string) (any, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("admin token is empty")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, fmt.Errorf("invalid date %q: expected YYYY-MM-DD", date)
	}
	if strings.TrimSpace(agentNo) == "" {
		return nil, fmt.Errorf("agentNo is empty")
	}

	u, err := url.Parse(findChangeURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("startTime", date)
	q.Set("endTime", date)
	q.Set("agentNo", agentNo)
	q.Set("pageNum", "1")
	q.Set("pageSize", strconv.Itoa(999))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+token)

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
		return nil, fmt.Errorf("findChange: HTTP %s: %s", resp.Status, raw)
	}

	var result any
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("findChange: invalid JSON: %w", err)
	}
	return result, nil
}
