// @Time:       2020/1/8 下午1:32

package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
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
	// 去掉最后一个周一的值, 保证关系的一对一
	a := max[len(max)-1]
	if _, ok := a["monday"]; ok {
		max = max[0 : len(max)-1]
	}

	must_1 := true
	must_5 := false
	var tmp = []map[string]string{}
	for _, i := range max {
		if must_1 {
			_, ok := i["monday"]
			if ok {
				tmp = append(tmp, i)
				must_1 = false
				must_5 = true
			} else {
				continue
			}
		}

		if must_5 {
			_, ok := i["friday"]
			if ok {
				tmp = append(tmp, i)
				must_1 = true
				must_5 = false
			} else {
				continue
			}
		}
	}

	// 上述保证了是 1515151515 成对
	for i := 0; i <= len(tmp)-1; i += 2 {
		a, b := tmp[i], tmp[i+1]
		a1 := str2Time(a["monday"])
		b1 := str2Time(b["friday"])
		log.Println(b1.Sub(a1).Hours(), a["monday"], b["friday"])
	}
	return tmp
}

func str2Time(str string) time.Time {
	timeLayout := "2006-01-02"                               //转化所需模板
	loc, _ := time.LoadLocation("Local")                     //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, str, loc) //使用模板在对应时区转化为time.time类型
	return theTime
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
