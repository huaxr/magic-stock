// @Contact:    huaxinrui
// @Time:       2019/10/20 下午5:54

package platform

import (
	"magic/stock/control"
	"magic/stock/service/middleware/normal"
)

func (r *Router) bindRouters() {
	r.addRouters()
}

func (r *Router) addRouters() {
	router := r.Router.Group("/api")
	router.POST("/callback/:order_id", control.UserControlGlobal.TradeCallBack)
	router.POST("/h5_pay", control.UserControlGlobal.PayByWeChatH5)
	router.GET("/wx_login", control.UserControlGlobal.LoginByWeChat)
	router.GET("/logout", control.UserControlGlobal.LogOut)
	router.Use(normal.LoginRequired())
	{
		router.GET("/user", control.UserControlGlobal.GetUserInfo)
		router.POST("/jsapi_pay", control.UserControlGlobal.PayByWeChatJsApi)
		router.GET("/my_conditions", control.UserControlGlobal.GetConditions)
		router.POST("/edit_condition", control.UserControlGlobal.EditUserConditions)
		router.POST("/delete_condition", control.UserControlGlobal.DeleteUserConditions)
		router.GET("/predicts_dates", control.PredictControlGlobal.GetPredictDates)
		router.GET("/conditions", control.PredictControlGlobal.GetConditions)
		router.GET("/concepts", control.PredictControlGlobal.GetConcepts)
		router.GET("/labels", control.PredictControlGlobal.GetLabels)
		router.GET("/belongs", control.PredictControlGlobal.GetBelongs)
		router.GET("/locations", control.PredictControlGlobal.GetLocations)
		router.POST("/predicts_list", control.PredictControlGlobal.PredictList)
		router.GET("/stock/detail", control.PredictControlGlobal.GetDetail)
		router.GET("/stock/fund", control.PredictControlGlobal.GetFunds)
		// 机构持仓
		router.GET("/fund_hold", control.PredictControlGlobal.FundHold)
		// 流通股东持仓
		router.GET("/top_holder_hold", control.PredictControlGlobal.TopHolderHold)

		router.GET("/stock/fhsgzz", control.PredictControlGlobal.GetFenHong)
		router.GET("/stock/pgzf", control.PredictControlGlobal.GetPeiGuZhuangZeng)
	}
}
