// @Contact:    huaxinrui
// @Time:       2019/8/29 上午10:48

package utils

func Reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
