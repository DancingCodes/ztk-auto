// Package admin 封装 tkapi 管理端业务接口。
package admin

import (
	"context"
	"net/http"
	"time"

	"ztk-auto/internal/httpx"
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
	url := findChangeURL + "?startTime=" + date + "&endTime=" + date + "&agentNo=" + agentNo + "&pageNum=1&pageSize=999"
	return httpx.JSON[any](ctx, c.httpClient, http.MethodGet, url, nil, map[string]string{"Authorization": "bearer " + token})
}
