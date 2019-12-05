package models

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
	Id    int    `json:"id"`
	Name  string `json:"name"`
	UName string `json:"u_name"`
	Email string `json:"email"`
}

type Emp struct {
	Id         int        `json:"id"`
	Name       string     `json:"name"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	AvatarUrl  string     `json:"avatar_url"`
	Department department `json:"department"`
	Leader     leader     `json:"leader"`
}

type Info struct {
	Employees []Emp `json:"employees"`
	Success   bool  `json:"success"`
}
