// @Time:       2019/12/4 下午2:26

package crawler

import (
	"fmt"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"testing"
)

// 获取具体日期的分析结果
func TestGetData(t *testing.T) {
	var date = "2019-12-25"
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Find(&code)
	for _, i := range code {
		x := &model.Params{i.Code, date, 0, 6, 15, 30, 10, 40}
		y := CrawlerGlobal.CalcResultWithDefined(x)
		if y == nil {
			continue
		}
		CrawlerGlobal.Analyze(y, i.Code, i.Name)
	}
}

func TestMultiQuery(t *testing.T) {
	var c []dal.Predict
	err := store.MysqlClient.GetDB().Model(&dal.Predict{}).
		Where("date = ?", "2019-12-20").
		//Where("`condition` regexp ?", "高位回调").
		//Where("`condition` regexp ?", "金叉").
		Where("`condition` regexp ?", "一字板").
		Find(&c).Error
	fmt.Println(err)
	for _, i := range c {
		fmt.Println(i.Code)
	}
}
