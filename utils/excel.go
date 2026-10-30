package utils

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
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1F4E78"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "bottom", Color: "17365D", Style: 1}},
	})
	bodyStyle, _ := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "222222", Size: 10},
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "bottom", Color: "D9E2F3", Style: 1}},
	})
	centerStyle, _ := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "222222", Size: 10},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "bottom", Color: "D9E2F3", Style: 1}},
	})
	_ = excel.SetCellStyle(sheet, "A1", "P1", headerStyle)
	_ = excel.SetRowHeight(sheet, 1, 32)
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
			_ = excel.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("P%d", row), bodyStyle)
			_ = excel.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), centerStyle)
			_ = excel.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("D%d", row), centerStyle)
			_ = excel.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), centerStyle)
			_ = excel.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("L%d", row), centerStyle)
			_ = excel.SetCellStyle(sheet, fmt.Sprintf("N%d", row), fmt.Sprintf("N%d", row), centerStyle)
			_ = excel.SetCellStyle(sheet, fmt.Sprintf("P%d", row), fmt.Sprintf("P%d", row), centerStyle)
			_ = excel.SetRowHeight(sheet, row, 22)
			row++
		}
	}
	buffer, err := excel.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
