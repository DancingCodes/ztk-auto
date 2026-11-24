package internal

import (
	"net/http"
	"strconv"
	"time"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/api/merchant-export", merchantExport)
	mux.HandleFunc("/api/transaction-export", transactionExport)
	mux.HandleFunc("/api/self-reg-export", selfRegExport)
	return mux
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

func selfRegExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "请求方式错误")
		return
	}

	query := r.URL.Query()
	accounts := query["account"]
	passwords := query["password"]
	serviceProviderNos := query["serviceProviderNo"]
	counts := query["count"]
	if len(accounts) == 0 || len(accounts) != len(passwords) || len(accounts) != len(serviceProviderNos) || len(accounts) != len(counts) {
		writeError(w, "account、password、serviceProviderNo、count 参数数量必须一致")
		return
	}

	requests := make([]SelfRegAccount, 0, len(accounts))
	for index, account := range accounts {
		count, err := strconv.Atoi(counts[index])
		if account == "" || passwords[index] == "" || serviceProviderNos[index] == "" || err != nil || count < 1 {
			writeError(w, "account、password、serviceProviderNo 不能为空且 count 必须大于 0")
			return
		}
		requests = append(requests, SelfRegAccount{
			Account:           account,
			Password:          passwords[index],
			ServiceProviderNo: serviceProviderNos[index],
			Count:             count,
		})
	}

	file, err := BuildSelfRegExport(r.Context(), requests)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	fileName := time.Now().Format("200601021504") + "-自主进件码导出.xlsx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(file)
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
