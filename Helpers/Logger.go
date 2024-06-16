package helpers

import (
	"encoding/json"
	initializers "service-auth/Initializers"

	"go.uber.org/zap"
)

// Custom Write log trong app
// Logtitle ->  tile của đoạn log
// objectLogString -> là key value trong zap.String(key, value)
// logType -> Info, Error, ...
func WriteLogApp(logTitle string, objectLogString map[string]interface{}, logType string) {
	fields := MapToZapFields(objectLogString)

	switch logType {
	case "INFO":
		initializers.Logger.Info(logTitle, fields...)
	case "ERROR":
		initializers.Logger.Error(logTitle, fields...)
	}
}

// create ZapField với object truyền vào
func MapToZapFields(data map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(data))
	for k, v := range data {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			fields = append(fields, zap.String(k, "error marshaling value"))
		} else {
			fields = append(fields, zap.String(k, string(jsonValue)))
		}
	}
	return fields
}
