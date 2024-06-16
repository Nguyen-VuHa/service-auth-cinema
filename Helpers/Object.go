package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Trả về keyName với value truyền vào và tìm trong object đó
func GetKeysByValue(objectData interface{}, value interface{}) []string {
	var keys []string

	val := reflect.ValueOf(objectData)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	valueToMatch := reflect.ValueOf(value)
	found := false

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Interface() == valueToMatch.Interface() {
			keys = append(keys, val.Type().Field(i).Name)
			found = true
		}
	}

	if !found {
		return []string{}
	}

	return keys
}

// Trả về value với key truyền vào và tìm trong object đó
func GetValueByKey(objectData interface{}, key string) (interface{}, bool) {
	val := reflect.ValueOf(objectData)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field := val.FieldByName(key)
	if !field.IsValid() {
		return nil, false
	}

	return field.Interface(), true
}

// Hàm tùy chỉnh để chuyển đổi struct sang JSON
func JSON_Stringify(v interface{}) (string, error) {
	// Chuyển đổi struct sang JSON
	jsonData, err := json.Marshal(v)

	if err != nil {
		return "", fmt.Errorf("error converting struct to JSON: %w", err)
	}

	return string(jsonData), nil
}
