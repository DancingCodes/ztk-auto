package main

import (
	"fmt"
	"net/http"
	ztk "ztk-auto/internal"
)

func main() {
	port := "8080"

	fmt.Println("服务已启动：http://localhost:" + port)
	_ = http.ListenAndServe(":"+port, ztk.Handler())
}
