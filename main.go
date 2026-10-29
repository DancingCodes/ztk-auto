package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ztk "ztk-auto/internal"
)

func main() {
	// 拿PCToken
	pc := ztk.LoadPC("config/pc.json")
	pcLogin := ztk.PCLogin(context.Background(), pc.Username, pc.Password, pc.OperateID)
	pcToken, _ := pcLogin["access_token"].(string)

	// 拿列表数据
	date := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	results := map[string]any{}
	for _, partner := range pc.PCPartners {
		results[partner] = ztk.PCFindChange(context.Background(), pcToken, date, partner)
	}
	result, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(result))


}
