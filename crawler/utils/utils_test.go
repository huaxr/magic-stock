// @Contact:    huaxinrui
// @Time:       2019/8/29 上午10:50

package utils

import (
	"fmt"
	"testing"
)

func TestReverse(t *testing.T) {
	x := []string{"a", "b", "c"}
	y := Reverse(x)
	fmt.Println(x, y)
}
