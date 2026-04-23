package kongc

import (
	"fmt"
	"strings"
)

type ErrCode int

type ApiErr struct {
	Code    *ErrCode `json:"code"`
	Fields  *any     `json:"fields"`
	Message *string  `json:"message"`
	Name    *string  `json:"name"`
}

const (
	UniqueErr = 5
)

func (e ApiErr) Error() string {
	parts := make([]string, 0, 4)

	if e.Code != nil {
		parts = append(parts, fmt.Sprintf("\"code\":%d", *e.Code))
	}
	if e.Message != nil {
		parts = append(parts, fmt.Sprintf("\"message\":\"%s\"", strings.ReplaceAll(*e.Message, "\"", "\\\"")))
	}
	if e.Name != nil {
		parts = append(parts, fmt.Sprintf("\"name\":\"%s\"", strings.ReplaceAll(*e.Name, "\"", "\\\"")))
	}
	if e.Fields != nil {
		parts = append(parts, fmt.Sprintf("\"fields\":\"%+v\"", *e.Fields))
	}

	return "{" + strings.Join(parts, ",") + "}"
}
