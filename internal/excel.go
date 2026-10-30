package internal

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func BuildMessageExcel(message string) ([]byte, error) {
	excel := excelize.NewFile()
	sheet := "提示"
	_ = excel.SetSheetName(excel.GetSheetName(0), sheet)
	_ = excel.SetCellValue(sheet, "A1", message)
	_ = excel.SetColWidth(sheet, "A", "A", 36)
	buffer, err := excel.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func BuildExcel(results []map[string]any) ([]byte, error) {
	excel := excelize.NewFile()
	sheet := "商户列表"
	_ = excel.SetSheetName(excel.GetSheetName(0), sheet)
	headers := []any{
		"合作方编号", "合作方名称", "账户类型",
		"商户号", "商户名称", "商户类型",
		"机具序列号", "联系人", "开通时间",
		"激活时间", "商户状态", "产品类型",
		"活动类型", "进件来源", "失效原因", "是否首台",
	}
	_ = excel.SetSheetRow(sheet, "A1", &headers)
	headerStyle, _ := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"D9E2F3"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	bodyStyle, _ := excel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	_ = excel.SetCellStyle(sheet, "A1", "P1", headerStyle)
	_ = excel.SetRowHeight(sheet, 1, 22)
	_ = excel.SetColWidth(sheet, "A", "P", 24)

	row := 2
	for _, result := range results {
		records := result["data"].(map[string]any)["records"].([]any)
		for _, value := range records {
			record := value.(map[string]any)
			isFirst, _ := record["isFirst"].(string)
			if isFirst == "" {
				appData := record["appData"].([]any)
				isFirst = "否"
				if appData[0].(map[string]any)["posAndIdentityFirst"] == true {
					isFirst = "是"
				}
			}
			values := []any{
				record["agentNo"], record["fullName"], record["agentNature"],
				record["externalCustomerNo"], record["customerName"], record["channelType"],
				record["posSn"], record["contactManName"], record["openTime"],
				record["activityTime"], record["customerStatus"], record["posType"],
				record["activityTypeName"], record["merchantSource"], record["failReason"], isFirst,
			}
			_ = excel.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &values)
			row++
		}
	}
	if row > 2 {
		_ = excel.SetCellStyle(sheet, "A2", fmt.Sprintf("P%d", row-1), bodyStyle)
	}
	buffer, err := excel.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func BuildTransactionExcel(results []map[string]any) ([]byte, error) {
	excel := excelize.NewFile()
	sheet := "交易列表"
	_ = excel.SetSheetName(excel.GetSheetName(0), sheet)
	headers := []any{
		"渠道名称", "商户编号", "商户名称", "机具序列号", "活动类型",
		"交易金额", "资金类型", "交易类型", "支付渠道", "流水号",
		"交易日期", "合作方编号", "合作方名称", "合作方角色",
	}
	_ = excel.SetSheetRow(sheet, "A1", &headers)
	headerStyle, _ := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"D9E2F3"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	bodyStyle, _ := excel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	_ = excel.SetCellStyle(sheet, "A1", "N1", headerStyle)
	_ = excel.SetRowHeight(sheet, 1, 22)
	_ = excel.SetColWidth(sheet, "A", "N", 24)

	row := 2
	for _, result := range results {
		records := result["data"].(map[string]any)["data"].(map[string]any)["records"].([]any)
		for _, value := range records {
			record := value.(map[string]any)
			values := []any{
				record["appCode"], record["externalCustomerNo"], record["customerName"], record["posSn"], record["activityName"],
				record["amount"], record["type"], record["extTransCode"], record["transCardType"], record["externalNo"],
				record["transTime"], record["agentNo"], record["agentName"], record["roleName"],
			}
			_ = excel.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &values)
			row++
		}
	}
	if row > 2 {
		_ = excel.SetCellStyle(sheet, "A2", fmt.Sprintf("N%d", row-1), bodyStyle)
	}
	buffer, err := excel.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
