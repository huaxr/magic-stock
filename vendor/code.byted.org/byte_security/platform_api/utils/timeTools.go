package utils

import (
	"errors"
	"math"
	"strings"
	"time"
)

func StringToLocTime(timeStr string) time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return tt
}

func StringToTime(timeStr string) time.Time {
	tt, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	return tt
}

func UnixToString(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func TimeToString(t time.Time) string {
	return time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05")
}

func CalculateInterval(startTime, endTime string, count int) int64 {
	sTime := StringToTime(startTime)
	eTime := StringToTime(endTime)
	interval := int64(math.Abs(eTime.Sub(sTime).Seconds())) / int64(count)
	if interval > 0 {
		return interval
	} else {
		return int64(3600 * 24)
	}
}

func CreateTimeBucketByInterval(startTime, endTime string, interval int64) ([]int, map[int]string, error) {
	if len(startTime) == 0 || len(endTime) == 0 || strings.HasPrefix(endTime, "000") {
		return nil, nil, errors.New("wrong startTime or endTime")
	}
	if interval == 0 {
		return nil, nil, errors.New("wrong interval")
	}
	sTime := StringToTime(startTime)
	eTime := StringToTime(endTime)
	var buckets []int
	bucketTimes := make(map[int]string, 0)
	for b := sTime.Unix() / interval; b <= eTime.Unix()/interval; b++ {
		buckets = append(buckets, int(b))
		bucketTimes[int(b)] = UnixToString(b * interval)
	}
	return buckets, bucketTimes, nil
}

func CreateTimeBucketByCount(startTime, endTime string, count int) ([]int, map[int]string, error) {
	interval := CalculateInterval(startTime, endTime, count)
	return CreateTimeBucketByInterval(startTime, endTime, interval)
}
