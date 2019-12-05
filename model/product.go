// @Time:       2019/12/1 下午4:08

package model

// 产品详情(新增DiscountPrice， ActivityLabel)
type ProductInfo struct {
	MerchantId    int     `json:"merchant_id"`    // 商户id
	ProductName   string  `json:"product_name"`   // 产品名称
	Unit          string  `json:"unit"`           // 计算单位 5g, 1kg
	Price         float64 `json:"price"`          // 产品价格  为浮点型变量
	PostType      string  `json:"post_type"`      // 物流类型
	PostMoney     float64 `json:"post_money"`     // 物流总价
	DiscountPrice float64 `json:"discount_price"` // 优惠价
	ActivityLabel string  `json:"activity_label"` // 活动标签
	VipPrice      float64 `json:"vip_price"`      // 会员价
	Desc          string  `json:"desc"`           // 产品描述
	Image1        string  `json:"image_1"`        // 产品图片1
	Image2        string  `json:"image_2"`        // 产品图片2
	Image3        string  `json:"image_3"`        // 产品图片3
	Video1        string  `json:"video_1"`        // 产品视频1
	Video2        string  `json:"video_2"`        // 产品视频2
	SaleCount     int     `json:"sale_count"`     // 已购数量
	LikeCount     int     `json:"like_count"`     // 喜欢/好评数量
}
