package helpers

import "reflect"

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
