package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"code.byted.org/gopkg/env"
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
	if env.IsBoe() {
		return "boe"
	} else if env.IsProduct() {
		return "tce"
	} else {
		return "loc"
	}
}

func GetRecentDayTimeStr(recent int) string {
	t := time.Now()
	nTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	yesTime := nTime.AddDate(0, 0, -recent)
	lastWeekDay := yesTime.Format("2006-01-02")
	return lastWeekDay
}

func GetTableName(field string) (string, string) {
	switch field {
	case "host":
		return "byte_security_asset_host", "ip"
	case "domain":
		return "byte_security_asset_domain", "name"
	case "psm":
		return "byte_security_asset_psm", "psm"
	default:
		return "", ""
	}
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符m +
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}
