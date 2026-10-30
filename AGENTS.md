# ztk-auto 开发说明

## 项目用途

这是一个 Go HTTP 服务，用于：

1. 使用 PC 账号登录并查询指定日期的商户变更列表。
2. 从合作方配置中读取对应的 App 账号。
3. 使用 App 查询商户是否为首台。
4. 将结果生成 Excel，并通过 HTTP 接口直接下载。

## 目录结构

```text
cmd/api/main.go       服务启动入口
internal/app.go       App 登录与商户查询接口
internal/pc.go        PC 登录与商户列表接口
internal/excel.go     Excel 内存生成逻辑
internal/server.go    HTTP 路由与请求参数处理
internal/service.go   PC/App 查询业务流程编排
config/               运行配置（真实配置不提交）
  accounts.json       Admin 登录账号、合作方编号及对应 App 账号
Dockerfile            Docker 镜像构建文件
```

## 启动与测试

本地启动：

```bash
go run ./cmd/api
```

健康检查：

```text
GET /health
```

导出接口：

```text
GET /api/merchant-export?date=YYYY-MM-DD
GET /api/transaction-export?date=YYYY-MM-DD
```

`date` 可省略；省略时默认查询昨天。接口会直接返回 `.xlsx` 文件。

## 配置约定

程序每次调用导出接口时，都会重新读取：

```text
config/accounts.json
```

因此修改账号、密码或合作方编号后，不需要重启服务；下一次接口请求会使用最新配置。

真实 `config/accounts.json` 包含敏感账号信息，必须保持在 `.gitignore` 中，禁止提交到仓库、镜像或日志。只提交对应的示例配置文件。

## 业务规则

- PC 登录失败：返回内容为“PC登录失败”的 Excel 提示文件，后续步骤不执行。
- PC 商户列表接口异常：返回内容为“PC商户列表获取失败：合作方编号”的 Excel 提示文件，后续步骤不执行。
- 某合作方没有 App 账号：该合作方的商户仍导出，“是否首台”显示“无账号”。
- App 登录失败：对应商户“是否首台”显示“App登录失败”。
- App 查询异常：对应商户“是否首台”显示“接口返回异常”。
- App 查询无商户数据：对应商户“是否首台”显示“无商户数据”。
- App 查询正常时：`posAndIdentityFirst == true` 显示“是”，其他情况显示“否”。

## 开发约定

- `main.go` 只负责启动服务；路由放在 `internal/server.go`；流程编排放在 `internal/service.go`。
- 新增外部接口时，按来源放入 `internal/pc.go` 或 `internal/app.go`；不要将请求代码写进 HTTP handler。
- Excel 只在内存中生成并直接返回 HTTP 响应，不要在服务端落地导出文件。
- 当前代码偏向简洁，避免无必要的层级、抽象和重复兼容判断。
- 修改 Go 代码后执行：

```bash
gofmt -w .
go test ./...
```

## Docker 部署

构建：

```bash
docker build -t ztk-auto:latest .
```

运行时通过目录挂载提供配置，不要把真实配置复制进镜像：

```yaml
volumes:
  - ./config:/app/config:ro
```

容器内部服务端口为 `8080`；对外端口由 `docker-compose.yml` 决定，例如：

```yaml
ports:
  - "9001:8080"
```

## 注意事项

- PC 合作方列表应避免重复编号，否则会重复请求和重复导出。
- 当前 App 查询为串行请求；商户量显著增长时，导出耗时会增加。并发优化前需先确认接口限流规则。
- 不要在代码、提交记录、截图、日志或接口响应中暴露账号、密码、token、客户端密钥等敏感信息。
