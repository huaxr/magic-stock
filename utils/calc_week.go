// @Time:       2020/1/8 下午1:32

package utils

import (
	"strconv"
	"time"
)

func GetWeekPair(date []string) (string, [][]interface{}) {
	//// 去除首部直到遇到1
	//for i := 0; i <= 5; i++ {
	//	if str2Week(date[0]) != 1 {
	//		date = date[1:len(date)]
	//	} else {
	//		break
	//	}
	//}
	//// 去除尾部直到遇到5
	//for i := 5; i >= 0; i-- {
	//	if str2Week(date[len(date)-1]) != 5 {
	//		date = date[0 : len(date)-1]
	//	} else {
	//		break
	//	}
	//}

	var x string
	var xx = [][]interface{}{}
	for _, i := range date {
		week := str2Week(i)
		x += strconv.Itoa(week)
		xx = append(xx, []interface{}{week, i})
	}
	y := calc(x)
	return y, xx

}

func str2Week(str string) int {
	timeLayout := "2006-01-02"                               //转化所需模板
	loc, _ := time.LoadLocation("Local")                     //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, str, loc) //使用模板在对应时区转化为time.time类型
	return int(theTime.Weekday())
}

func calc(x string) string {
	for j := 0; j <= len(x)-1; j++ {
		i := string(x[j])
		//log.Println(j, string(i), j%5)
		if j%5 == 0 {
			if string(i) != "1" && string(i) != "0" {
				x = x[0:j] + "0" + x[j:len(x)]
			}
		}
		if j%5 == 1 {
			if string(i) != "2" && string(i) != "0" {
				x = x[0:j] + "0" + x[j:len(x)]
			}
		}
		if j%5 == 2 {
			if string(i) != "3" && string(i) != "0" {
				x = x[0:j] + "0" + x[j:len(x)]
			}
		}
		if j%5 == 3 {
			if string(i) != "4" && string(i) != "0" {
				x = x[0:j] + "0" + x[j:len(x)]
			}
		}
		if j%5 == 4 {
			if string(i) != "5" && string(i) != "0" {
				x = x[0:j] + "0" + x[j:len(x)]
			}
		}
	}

	return x
}
