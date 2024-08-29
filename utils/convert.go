package utils

import (
	"fmt"
	"reflect"
	"strconv"
)

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

func MapStringToStructure(m map[string]string, s interface{}) error {
	// Lấy giá trị của struct
	v := reflect.ValueOf(s).Elem()

	// Lặp qua các phần tử của map
	for key, value := range m {
		// Lấy field tương ứng trong struct
		f := v.FieldByName(key)

		// Kiểm tra xem field có tồn tại và có thể được gán giá trị không
		if f.IsValid() && f.CanSet() {
			switch f.Kind() {
			case reflect.String:
				f.SetString(value)
			case reflect.Int:
				// Convert string to int
				intVal, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("failed to convert %s to int for field %s: %w", value, key, err)
				}
				f.SetInt(int64(intVal))
			// Thêm các kiểu dữ liệu khác nếu cần
			default:
				return fmt.Errorf("unsupported kind %s for field %s", f.Kind(), key)
			}
		}
	}

	return nil
}
