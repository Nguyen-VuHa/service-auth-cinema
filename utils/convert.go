package utils

import "reflect"

// Hàm chuyển struct thành map[string]interface{}
func StructureToMapString(obj interface{}) map[string]interface{} {

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		result[typ.Field(i).Name] = field.Interface()
	}

	return result
}
