// @Time:       2019/12/15 下午7:35

package control

import (
	"fmt"
	"magic/stock/core/store"
	"magic/stock/dal"
	"testing"
)

type xx struct {
	Date string `json:"date"`
}

func TestQuery(t *testing.T) {
	var x []xx
	err := store.MysqlClient.GetDB().Model(&dal.Predict{}).Select("distinct(date) as date").Order("date desc").Scan(&x).Error
	fmt.Println(err, x)
}
