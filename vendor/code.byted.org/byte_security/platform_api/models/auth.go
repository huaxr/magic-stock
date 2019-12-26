package models

import "time"

type AuthRequestBody struct {
	Url     string                      `json:"url"`
	Body    map[interface{}]interface{} `json:"body"`
	Message string                      `json:"message"`
}

type AdminAddRule struct {
	UserName string `json:"user_name"`
	Role     string `json:"role"`
}

type AuthResult struct {
	Err         error    // 鉴权是否失败
	User, Group string   // 用户和用户组
	Email       string   // 用户邮箱
	Uid         int      // 用户表中id
	Admin       bool     // 是否管理员
	Perms       []string // kani 权限
}

type SsoRes struct {
	User       string
	EmployeeId string
	Name       string
}

type GroupQuery struct {
	Employees []map[string]interface{}
	Success   bool
}

func (a *AuthResult) HasPerm() string {
	return ""
}

type department struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type leader struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	UName       string `json:"u_name"`
	Email       string `json:"email"`
	EmployeeNum int    `json:"employee_number"`
}

type PeopleEmp struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Terminated  bool       `json:"terminated"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Sequence    string     `json:"sequence"`
	AvatarUrl   string     `json:"avatar_url"`
	Department  department `json:"department"`
	Leader      leader     `json:"leader"`
	EmployeeNum int        `json:"employee_number"`
}

type PeopleRspInfo struct {
	Employees []PeopleEmp `json:"employees"`
	Success   bool        `json:"success"`
}

type UserInfoForWeb struct {
	Name      string `json:"name"`
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// 用户完整信息，仅在平台内部调用使用，不可返回到前端
type UserInfoInner struct {
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserName  string    `json:"user_name"`
	RealName  string    `json:"real_name"`
	UserNum   string    `json:"user_num"`
	GroupID   string    `json:"group_id"`
	GroupName string    `json:"group_name"`
	AvatarURL string    `json:"avatar_url"`
	Leader    string    `json:"leader"`
	Email     string    `json:"email"`
}
