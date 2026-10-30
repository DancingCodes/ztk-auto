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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("ok"))
}

func merchantExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	date, ok := exportDate(w, r)
	if !ok {
		return
	}

	file, _ := BuildMerchantExport(r.Context(), date)
	fileName := time.Now().Format("200601021504") + "-商户查询导出.xlsx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Write(file)
}

func transactionExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	date, ok := exportDate(w, r)
	if !ok {
		return
	}

	file, _ := BuildTransactionExport(r.Context(), date)
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
		http.Error(w, "date 格式应为 YYYY-MM-DD", http.StatusBadRequest)
		return "", false
	}
	return date, true
}
