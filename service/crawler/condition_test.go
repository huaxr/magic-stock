// @Time:       2019/12/4 下午2:26

package crawler

import (
	"fmt"
	"magic/stock/core/store"
	"magic/stock/dal"
	"testing"
)

func TestMultiQuery(t *testing.T) {
	var c []dal.Predict
	err := store.MysqlClient.GetDB().Model(&dal.Predict{}).
		Where("date = ?", "2020-01-03").
		//Where("`condition` regexp ?", "高位回调").
		//Where("`condition` regexp ?", "金叉").
		Where("`condition` regexp ?", "近期60日均线与收盘价黏合").
		Find(&c).Error
	fmt.Println(err)
	for _, i := range c {
		fmt.Println(i.Code)
	}
}
