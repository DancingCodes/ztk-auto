package main

import (
	"fmt"
	"net/http"
	ztk "ztk-auto/internal"
)

func main() {
	port := "9101"

	fmt.Println("服务已启动：port:" + port)
	_ = http.ListenAndServe(":"+port, ztk.Handler())
}
