package seal

import (
	"time"

	"code.byted.org/byte_security/dal/common"
)

type User struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	UserName  string `gorm:"not null;unique"`
	FullName  string
	Email     string
	RoleList  string `sql:"type:text"`
	Wifi      string `sql:"type:text"`
	IsTest    bool
	IsAdmin   bool
	IsGrey    int
	Extra     common.JSON `sql:"type:json" json:"extra,omitempty"`
}

func (User) TableName() string {
	return "byte_security_seal_user"
}
