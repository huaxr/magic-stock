package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/axgle/mahonia"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// refer to the test, this is the map w struct method definition
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func MapToStr(m map[interface{}]interface{}) string {
	s := ""
	for k, v := range m {
		if (k != "" || k != nil) && (v != "" || v != nil) {
			s += k.(string) + "=" + v.(string) + "&"
		}
	}
	if len(s) == 0 {
		return ""
	}
	return s
}

func Contains(list []interface{}, value interface{}) bool {
	for _, i := range list {
		if value == i {
			return true
		}
	}
	return false
}

func ContainsString(list []string, value string) bool {
	for _, i := range list {
		if value == i {
			return true
		}
	}
	return false
}

func TellEnv() string {
	if os.Getenv("USER") == "tiger" {
		return "online"
	}
	return "loc"
}

// gbk 乱码转 utf-8
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func TellWeChat(agent string) bool {
	if strings.Contains(strings.ToLower(agent), "micromessenger") {
		return true
	}
	return false
}
