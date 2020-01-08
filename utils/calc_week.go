// @Time:       2020/1/8 下午1:32

package utils

import (
	"fmt"
	"strconv"
	"strings"
)

var weekday = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

type Monday struct {
	Monday string
}

type Friday struct {
	Friday string
}

// 例如  1 5 1 5 5 5 1 5 1 5 1 1 1 5 1 5

func GetWeekPair(date []string) {
	//pairs := []Pair{}
	getRes(date)

}

// 获取  [map[friday:1991-01-04] map[monday:1991-01-07] map[friday:1991-01-11] map[monday:1991-01-14] 。。。 一个序列
func getRes(dates []string) []map[string]string {
	max := []map[string]string{}
	for _, i := range dates {
		ma := map[string]string{}
		x := strings.Split(i, "-")
		u, m, d := parseUint(x[0]), parseUint(x[1]), parseUint(x[2])
		//fmt.Printf("%d年%d月%d日是:%s\n", u, m, d, ZellerFunction2Week(u, m, d))
		res := ZellerFunction2Week(u, m, d)
		if res == 5 {
			ma["friday"] = i
			max = append(max, ma)
		}
		if res == 1 {
			ma["monday"] = i
			max = append(max, ma)
		}
	}
	return max
}

func parseUint(str string) uint16 {
	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	wd := uint16(u64)
	return wd
}

func ZellerFunction2Week(year, month, day uint16) int {
	var y, m, c uint16
	if month >= 3 {
		m = month
		y = year % 100
		c = year / 100
	} else {
		m = month + 12
		y = (year - 1) % 100
		c = (year - 1) / 100
	}

	week := y + (y / 4) + (c / 4) - 2*c + ((26 * (m + 1)) / 10) + day - 1
	if week < 0 {
		week = 7 - (-week)%7
	} else {
		week = week % 7
	}
	return int(week)
}
