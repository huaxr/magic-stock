// @Time:       2019/12/4 下午2:26

package crawler

import (
	"fmt"
	"magic/stock/core/store"
	"magic/stock/dal"
	"magic/stock/model"
	"testing"
)

// 得出基金排行并根据这些基金获取持仓股
func TestGetData(t *testing.T) {
	var date = "2019-12-13"
	var code []dal.Code
	store.MysqlClient.GetDB().Model(&dal.Code{}).Where("id>0").Find(&code)
	for _, i := range code {
		x := &model.Params{i.Code, 0, 6, 15, 30, 10, 40}
		y := CrawlerGlobal.CalcResultWithDefined(x, date)
		if y == nil {
			continue
		}
		CrawlerGlobal.Analyze(y, i.Code, i.Name)
	}
}

func TestMultiQuery(t *testing.T) {
	var c []dal.Predict
	err := store.MysqlClient.GetDB().Model(&dal.Predict{}).
		//Where("`condition` regexp ?", "长上影").
		Where("`condition` regexp ?", "总负债不断减小").
		//Where("`condition` regexp ?", "15日均线与30日均线交金叉").
		Where("`condition` regexp ?", "近5日资金净流入总和非负").
		Find(&c).Error
	fmt.Println(err)
	for _, i := range c {
		fmt.Println(i.Code)
	}
}

// 免责声明：本终端所载数据仅供参考，若数据有误，以交易所发布数据为准，不对您构成投资建议。
