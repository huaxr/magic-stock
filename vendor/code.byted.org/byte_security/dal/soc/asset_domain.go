// 资产域名

package soc

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

// 域名->产品->责任人
// 域名->责任人
type Domain struct {
	Id        int         `gorm:"primary_key" json:"id"`
	Created   time.Time   `json:"created"`
	Updated   time.Time   `json:"updated"`
	Name      string      `gorm:"not null;unique" json:"name"`
	Visible   int         `json:"visible"`
	Product   Product     `gorm:"ForeignKey:ProductID" json:"product"`
	ProductID int         `json:"product_id"`
	Owner     string      `json:"owner"`
	PSM       string      `json:"psm"`
	Icp       common.JSON `sql:"type:json" json:"icp"`
	Whois     common.JSON `sql:"type:json" json:"whois"`
	Cert      common.JSON `sql:"type:json" json:"cert"`
	Dns       common.JSON `sql:"type:json" json:"dns"`
	Cluster   common.JSON `sql:"type:json" json:"cluster"` // 集群信息
	Status    string      `json:"status"`
	Scope     string      `json:"scope"`                 //可见范围
	AccType   string      `json:"acc_type"`              // 业务类型
	Operator  string      `json:"operator"`              //运营商
	Vendor    string      `json:"vendor"`                // cdn
	Comment   string      `json:"comment"`               // 域名用途
	GroupId   string      `gorm:"index" json:"group_id"` // 所属组
	Debug     string      `json:"debug"`                 // 对外开放端口开启调试模式的端口
	Extra     common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (Domain) New() interface{} {
	return &Domain{}
}

func (Domain) TableName() string {
	return "byte_security_asset_domain"
}

func (domain Domain) GetAssetKey() string {
	return "name"
}

func (domain Domain) GetAssetValue() string {
	return domain.Name
}

//func (d *Domain) GetOwnerNames() (name []string) {
//	name = append(name, d.Owner)
//	common.Backend.DB.Model(&d).Related(&d.Product).Find(&d.Product)
//	name = append(name, d.Product.Owner)
//	return
//}
//
//func (d *Domain) GetPSM() {
//	c := fmt.Sprintf("curl -I -H'Get-Svc: 1' %s", d.Name)
//	psm := api.GetPSMByDomain(c)
//	fmt.Println(psm)
//}
