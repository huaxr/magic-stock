package model

type AuthResult struct {
	Err       error  // 鉴权是否失败
	User      string // 用户和用户组
	Uid       int    // 用户表中id
	Member    bool   // 会员
	QueryLeft int    // 查询的剩余次数
}
