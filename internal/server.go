package internal

import (
	"net/http"
	"time"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/api/merchant-export", merchantExport)
	mux.HandleFunc("/api/transaction-export", transactionExport)
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

	file, err := BuildMerchantExport(r.Context(), date)
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

	file, err := BuildTransactionExport(r.Context(), date)
	if err != nil {
		writeError(w, err.Error())
		return
	}
	fileName := time.Now().Format("200601021504") + "-交易查询导出.xlsx"
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

func writeError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
