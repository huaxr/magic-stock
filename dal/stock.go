// @Time:       2019/12/2 下午2:28

package dal

type Code struct {
	ID      uint   `gorm:"primary_key" json:"id"`
	Code    string `gorm:"not null;unique" json:"code"`
	Name    string `json:"name"`
	Belong  string `json:"belong"`  // 所属行业板块
	Concept string `json:"concept"` // 所属概念板块

	CompanyName        string `json:"company_name"`                   // 公司全称
	OrganizationalForm string `json:"organizational_form"`            // 组织形式
	InstitutionalType  string `json:"institutional_type"`             // 机构类型
	Location           string `json:"location"`                       // 地域
	Address            string `json:"address"`                        // 办公地址
	NetAddress         string `json:"net_address"`                    // 公司网址
	MajorBusinesses    string `json:"major_businesses"`               // 主营业务
	BusinessScope      string `sql:"type:text" json:"business_scope"` // 经营范围
	EstablishmentTime  string `json:"establishment_time"`             // 成立日期
	ListingDate        string `json:"listing_date"`                   // 上市日期
	HistoryNames       string `json:"history_names"`                  // 历史用名
	Tape               string `json:"tape"`                           // 盘口 是从 concept 拿到: 小盘,中盘,大盘,超大盘
}

func (Code) TableName() string {
	return "magic_stock_code"
}
