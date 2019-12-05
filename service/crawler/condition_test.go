// @Time:       2019/12/4 下午2:26

package crawler

import (
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"testing"
)

// 得出基金排行并根据这些基金获取持仓股
func TestGetData(t *testing.T) {
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id>0").Find(&code)
	for _, i := range code {
		x := &model.Params{i.Code, 0, 6, 15, 30, 10, 40}
		y := CrawlerGlobal.CalcResultWithDefined(x)
		if y == nil {
			continue
		}
		CrawlerGlobal.Analyze(y, i.Code, i.Name)
	}
}
