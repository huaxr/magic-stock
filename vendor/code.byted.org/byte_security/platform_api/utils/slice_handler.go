package utils

import (
	"math/rand"
	"strings"
	"time"
)

func StringsDistinct(src []string) (dst []string) {
	dst = make([]string, 0)
	for i := 0; i < len(src); i++ {
		repeat := false
		for j := i + 1; j < len(src); j++ {
			if src[i] == src[j] {
				repeat = true
				break
			}
		}
		if !repeat && len(strings.TrimSpace(src[i])) > 0 {
			dst = append(dst, src[i])
		}
	}
	return
}

func RandomString(src []string) string {
	if len(src) == 0 {
		return ""
	} else if len(src) == 1 {
		return src[0]
	}
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
	return src[rand.Intn(len(src))]
}
