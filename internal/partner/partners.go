// Package partner 读取本地合作方编号配置。
package partner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type file struct {
	Partners []string `json:"partners"`
}

// Load 读取合作方编号，过滤空项并拒绝重复编号。
func Load(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg file
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parse partners: %w", err)
	}

	seen := make(map[string]struct{}, len(cfg.Partners))
	partners := make([]string, 0, len(cfg.Partners))
	for i, value := range cfg.Partners {
		id := strings.TrimSpace(value)
		if id == "" {
			return nil, fmt.Errorf("partner #%d is empty", i+1)
		}
		if _, exists := seen[id]; exists {
			return nil, fmt.Errorf("duplicate partner: %s", id)
		}
		seen[id] = struct{}{}
		partners = append(partners, id)
	}
	if len(partners) == 0 {
		return nil, fmt.Errorf("partners is empty")
	}
	return partners, nil
}
