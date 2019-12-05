// @Contact:    huaxinrui
// @Time:       2019/10/14 下午5:36

package soc

type VulnerabilityType struct {
	ID   int `gorm:"primary_key"`
	Name string
	SubName string
}

func (VulnerabilityType) TableName() string {
	return "byte_security_asset_vulntype"
}
