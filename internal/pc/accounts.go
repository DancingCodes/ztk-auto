package pc

import (
	"encoding/json"
	"os"
	"strings"
)

func LoadPartners(path string) []string {
	raw, _ := os.ReadFile(path)
	var data struct {
		PCPartners []string `json:"pcPartners"`
	}
	_ = json.Unmarshal(raw, &data)
	partners := make([]string, 0, len(data.PCPartners))
	for _, item := range data.PCPartners {
		if value := strings.TrimSpace(item); value != "" {
			partners = append(partners, value)
		}
	}
	return partners
}
