package main

import (
	"fmt"
	"net/http"
	"os"

	ztk "ztk-auto/internal"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("服务已启动：http://localhost:" + port)
	_ = http.ListenAndServe(":"+port, ztk.Handler())
}
