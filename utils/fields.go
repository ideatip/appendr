package utils

import (
	"fmt"
	"go.ideatip.dev.appendr/models"
)

func FieldsToString(fields []models.Field) string {
	result := ""
	for _, field := range fields {
		result += fmt.Sprintf("%s=%v ", field.Key, field.Value)
	}
	return result
}

func FieldsToMap(fields []models.Field) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		result[field.Key] = field.Value
	}
	return result
}
