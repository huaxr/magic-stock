// @Contact:    huaxinrui
// @Time:       2019/7/5 下午5:00

package utils

import (
	"fmt"
	"strings"
	"time"
)

func GetPriceLink(code string) string{
	if strings.HasPrefix(code, "6"){
		code = "sh" + code //上海
	}  else {
		code = "sz" + code //深圳
	}

	now := time.Now()
	end := now.AddDate(0, 0 , 0).Format("2006-01-02")
	start := now.AddDate(0, 0 , -40).Format("2006-01-02")

	s := fmt.Sprintf("http://market.finance.sina.com.cn/pricehis.php?symbol=%s&startdate=%s&enddate=%s", code, start, end)
	return s
}
