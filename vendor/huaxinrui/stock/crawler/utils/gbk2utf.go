// @Contact:    huaxinrui
// @Time:       2019/7/5 下午4:53

package utils

import "github.com/axgle/mahonia"


// gbk 乱码转 utf-8
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}