package blemodule

import (
	"regexp"
	"strings"
)

// urc.go 负责解析 BLE 模块的 URC 行
// 例如 +QBLESTAT, +QBLEINFO, +QBLEPEERINFO 等

type URCEvent struct {
	Type   string
	Fields []string
	Raw    string
}

// ParseURC 解析一行 URC 字符串
func ParseURC(line string) *URCEvent {
	if !strings.HasPrefix(line, "+QBLE") {
		return nil
	}
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil
	}
	eventType := strings.TrimSpace(parts[0][1:])
	fields := regexp.MustCompile(`\s*,\s*`).Split(strings.TrimSpace(parts[1]), -1)
	return &URCEvent{
		Type:   eventType,
		Fields: fields,
		Raw:    line,
	}
}

// TODO: 可扩展为更详细的结构体和事件分发
