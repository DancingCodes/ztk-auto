package main

import (
	"context"
	"encoding/json"
	"fmt"

	ztk "ztk-auto/internal"
)

func main() {
	pcAccount := ztk.LoadPCAccount("config/pcAccount.json")
	pcLogin := ztk.PCLogin(context.Background(), pcAccount.Username, pcAccount.Password)
	pcToken, _ := pcLogin["access_token"].(string)
	_ = pcToken
	result, _ := json.MarshalIndent(pcLogin, "", "  ")
	fmt.Println(string(result))
}
