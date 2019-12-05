package auth

//type AdminGroup struct {
//	IsAdmin bool
//}
//
//type BlackGroup struct {
//	IsBlack bool
//}
//
//type AttrGroup struct {
//	IsAttr bool
//}
//
//type AssetGroup struct {
//	Domain     uint
//	Host       uint
//	Ip         uint
//	Pc         uint
//	Product    uint
//	App        uint
//	CheckAsset uint // 资产详情鉴权细化
//}
//
//type EventGroup struct {
//	Auth       uint
//	Vuln       uint
//	CheckEvent uint // 事件详情鉴权细化
//}
//
//type DepartmentGroup struct {
//	SecGroup uint
//	OpsGroup uint
//	Other    uint // 当为其它部门的时候 需要细粒度前面的
//}
//
//// 改条group记录属于用户
//type Group struct {
//	ID uint `gorm:"primary_key"`
//	AdminGroup
//	BlackGroup
//	AttrGroup
//	AssetGroup
//	EventGroup
//	DepartmentGroup
//	User User  `gorm:"ForeignKey:UserID"`
//	UserID uint
//}


type Group struct {
	Id int `gorm:"primary_key"`
	GroupId string `gorm:"not null;unique"`
	GroupName string
}



func (Group) TableName() string {
	return "byte_security_auth_group"
}


