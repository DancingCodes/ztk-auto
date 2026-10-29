package utils

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

func ExportExcel(results []map[string]any) error {
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
	_ = excel.SetColWidth(sheet, "A", "A", 14)
	_ = excel.SetColWidth(sheet, "B", "B", 14)
	_ = excel.SetColWidth(sheet, "C", "C", 14)
	_ = excel.SetColWidth(sheet, "D", "D", 21)
	_ = excel.SetColWidth(sheet, "E", "E", 28)
	_ = excel.SetColWidth(sheet, "F", "F", 18)
	_ = excel.SetColWidth(sheet, "G", "G", 24)
	_ = excel.SetColWidth(sheet, "H", "H", 14)
	_ = excel.SetColWidth(sheet, "I", "J", 20)
	_ = excel.SetColWidth(sheet, "K", "K", 12)
	_ = excel.SetColWidth(sheet, "L", "L", 14)
	_ = excel.SetColWidth(sheet, "M", "M", 28)
	_ = excel.SetColWidth(sheet, "N", "N", 14)
	_ = excel.SetColWidth(sheet, "O", "O", 22)
	_ = excel.SetColWidth(sheet, "P", "P", 12)

	row := 2
	for _, result := range results {
		records := result["data"].(map[string]any)["records"].([]any)
		for _, value := range records {
			record := value.(map[string]any)
			appData := record["appData"].([]any)
			isFirst := "否"
			if appData[0].(map[string]any)["posAndIdentityFirst"] == true {
				isFirst = "是"
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
	fileName := time.Now().Format("200601021504") + "-商户查询导出.xlsx"
	if err := excel.SaveAs(fileName); err != nil {
		return err
	}
	fmt.Println("已生成 " + fileName)
	return nil
}
