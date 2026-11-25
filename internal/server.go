package internal

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed web
var selfRegFiles embed.FS

func Handler() http.Handler {
	mux := http.NewServeMux()
	webFiles, _ := fs.Sub(selfRegFiles, "web")
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/api/merchant-export", merchantExport)
	mux.HandleFunc("/api/transaction-export", transactionExport)
	mux.HandleFunc("/api/self-reg", selfReg)
	mux.HandleFunc("/self-reg", selfRegPage)
	mux.Handle("/self-reg-assets/", http.StripPrefix("/self-reg-assets/", http.FileServer(http.FS(webFiles))))
	return mux
}

func selfRegPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "请求方式错误")
		return
	}
	page, err := selfRegFiles.ReadFile("web/self-reg.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(page)
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "请求方式错误")
		return
	}
	w.Write([]byte("ok"))
}

func merchantExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "请求方式错误")
		return
	}

	date, ok := exportDate(w, r)
	if !ok {
		return
	}
	accounts, ok := exportAccounts(w, r)
	if !ok {
		return
	}

	file, err := BuildMerchantExport(r.Context(), date, accounts)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	fileName := time.Now().Format("200601021504") + "-商户查询导出.xlsx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(file)
}

func transactionExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "请求方式错误")
		return
	}

	date, ok := exportDate(w, r)
	if !ok {
		return
	}
	accounts, ok := exportAccounts(w, r)
	if !ok {
		return
	}

	file, err := BuildTransactionExport(r.Context(), date, accounts)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	fileName := time.Now().Format("200601021504") + "-交易查询导出.xlsx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(file)
}

func selfReg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "请求方式错误")
		return
	}

	var request struct {
		AccountID string `json:"accountId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.AccountID == "" {
		writeError(w, "accountId 参数不能为空")
		return
	}

	account, err := LoadSelfRegAccount(request.AccountID)
	if err != nil {
		writeError(w, "账号配置不存在")
		return
	}
	result, err := BuildSelfReg(r.Context(), account)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	log.Printf("自主进件码生成成功 accountId=%s selfRegSn=%s", request.AccountID, result.SelfRegSN)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(result)
}

func exportDate(w http.ResponseWriter, r *http.Request) (string, bool) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		writeError(w, "date 格式应为 YYYY-MM-DD")
		return "", false
	}
	return date, true
}

func exportAccounts(w http.ResponseWriter, r *http.Request) (Accounts, bool) {
	client := r.URL.Query().Get("client")
	if client == "" {
		writeError(w, "client 参数不能为空")
		return Accounts{}, false
	}
	accounts, err := LoadClientAccounts(client)
	if err != nil {
		writeError(w, "客户配置不存在")
		return Accounts{}, false
	}
	return accounts, true
}

func writeError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
