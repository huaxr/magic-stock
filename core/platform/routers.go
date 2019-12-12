// @Contact:    huaxinrui
// @Time:       2019/10/20 下午5:54

package platform

import "magic/stock/control"

func (r *Router) bindRouters() {
	r.addRouters()
}

func (r *Router) addRouters() {
	router := r.Router.Group("/api")
	router.GET("/user", control.UserControlGlobal.GetUserInfo)
	router.GET("/wx_login", control.UserControlGlobal.LoginByWeChat)
	router.POST("/jsapi_pay", control.UserControlGlobal.PayByWeChat)
	router.GET("/callback/:order_id", control.UserControlGlobal.TradeCallBack)

	router.POST("/predicts_list", control.PredictControlGlobal.GetPredict)
	router.GET("/stock/detail", control.PredictControlGlobal.GetDetail)
	router.GET("/stock/fund", control.PredictControlGlobal.GetFunds)
	// 机构持仓
	router.GET("/fund_hold", control.PredictControlGlobal.FundHold)
	// 流通股东持仓
	router.GET("/top_holder_hold", control.PredictControlGlobal.TopHolderHold)
}
