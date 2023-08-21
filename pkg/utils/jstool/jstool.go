package jstool

import (
	"encoding/json"
)

func StructToStr(v any) string {
	t, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(t)
}

func StructToMap(v any) map[string]any {
	t, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	m := make(map[string]any)
	err = json.Unmarshal(t, &m)
	if err != nil {
		return nil
	}
	return m
}
