// @Contact:    huaxinrui
// @Time:       2019/10/10 下午2:59

package workflow

type NodeUser struct {
	ID       uint   `gorm:"primary_key"`
	NodeId   int    `gorm:"index"`
	UserName string `gorm:"index"`
}

func (NodeUser) TableName() string {
	return "byte_security_workflow_node_user"
}
