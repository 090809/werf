package secretvalues

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func MaskSecretValuesInString(secretValues []string, targetStr string) string {
	for _, secretValue := range secretValues {
		targetStr = strings.ReplaceAll(targetStr, secretValue, "***")
	}
	return targetStr
}

func ExtractSecretValuesFromMap(data map[string]interface{}) []string {
	queue := []interface{}{data}
	maskedValues := []string{}

	for len(queue) > 0 {
		var elemI interface{}
		elemI, queue = queue[0], queue[1:]

		switch reflect.TypeOf(elemI).Kind() {
		case reflect.Slice, reflect.Array:
			elem := reflect.ValueOf(elemI)
			for i := 0; i < elem.Len(); i++ {
				value := elem.Index(i)
				queue = append(queue, value.Interface())
			}
		case reflect.Map:
			elem := reflect.ValueOf(elemI)
			for _, key := range elem.MapKeys() {
				value := elem.MapIndex(key)
				queue = append(queue, value.Interface())
			}
		default:
			elemStr := fmt.Sprintf("%v", elemI)
			if len(elemStr) >= 4 {
				maskedValues = append(maskedValues, elemStr)
			}
			for _, line := range strings.Split(elemStr, "\n") {
				trimmedLine := strings.TrimSpace(line)
				if len(trimmedLine) >= 4 {
					maskedValues = append(maskedValues, trimmedLine)
				}
			}

			dataMap := map[string]interface{}{}
			if err := json.Unmarshal([]byte(elemStr), &dataMap); err == nil {
				for _, v := range dataMap {
					queue = append(queue, v)
				}
			}

			dataArr := []interface{}{}
			if err := json.Unmarshal([]byte(elemStr), &dataArr); err == nil {
				for _, v := range dataArr {
					queue = append(queue, v)
				}
			}
		}
	}

	return maskedValues
}
