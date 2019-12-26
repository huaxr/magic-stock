package d_system

import (
	"code.byted.org/byte_security/dal/system"
	"code.byted.org/byte_security/platform_api/core/store"
	"code.byted.org/gopkg/logs"
)

type RobotIF struct {
	Dal store.ModelInterface
}

var RobotDao RobotIF

func init() {
	tmp := new(RobotIF)
	tmp.Dal = store.DB
	RobotDao = *tmp
}

func (r *RobotIF) GetRobotInfoByID(id int) (robot *system.FAASRobot, err error) {
	param := r.Dal.NewDefaultQueryParam()
	param.Table = system.FAASRobot{}
	param.Where = "id = ?"
	param.Args = []interface{}{id}
	result, err := r.Dal.QueryOne(param)
	if err != nil {
		logs.Errorf("【RobotIF】GetRobotInfoByID error:%v", err)
		return
	}
	robot = result.(*system.FAASRobot)
	return
}

func (r *RobotIF) GetRobotInfoByURL(url string) (robot *system.FAASRobot, err error) {
	param := r.Dal.NewDefaultQueryParam()
	param.Table = system.FAASRobot{}
	param.Where = "url = ?"
	param.Args = []interface{}{url}
	result, err := r.Dal.QueryOne(param)
	if err != nil {
		logs.Errorf("【RobotIF】GetRobotInfoByID error:%v", err)
		return
	}
	robot = result.(*system.FAASRobot)
	return
}
