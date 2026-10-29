package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ExportExcel(results []map[string]any) {
	excel := excelize.NewFile()
	sheet := "商户列表"
	_ = excel.SetSheetName(excel.GetSheetName(0), sheet)
	headers := []any{
		"合作方编号", "合作方名称", "账户类型",
		"商户号", "商户名称", "商户类型",
		"机具序列号", "联系人", "开通时间",
		"激活时间", "商户状态", "产品类型",
		"活动类型", "进件来源", "失效原因",
	}
	_ = excel.SetSheetRow(sheet, "A1", &headers)
	headerStyle, _ := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1D2330"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	_ = excel.SetCellStyle(sheet, "A1", "O1", headerStyle)
	_ = excel.SetColWidth(sheet, "A", "O", 18)
	_ = excel.SetColWidth(sheet, "E", "E", 28)
	_ = excel.SetColWidth(sheet, "J", "J", 22)

	row := 2
	for _, result := range results {
		records := result["data"].(map[string]any)["records"].([]any)
		for _, value := range records {
			record := value.(map[string]any)
			values := []any{
				record["agentNo"], record["fullName"], record["agentNature"],
				record["externalCustomerNo"], record["customerName"], record["channelType"],
				record["posSn"], record["contactManName"], record["openTime"],
				record["activityTime"], record["customerStatus"], record["posType"],
				record["activityTypeName"], record["merchantSource"], record["failReason"],
			}
			_ = excel.SetSheetRow(sheet, fmt.Sprintf("A%d", row), &values)
			row++
		}
	}
	_ = excel.SaveAs("result.xlsx")
	fmt.Println("已生成 result.xlsx")
}
