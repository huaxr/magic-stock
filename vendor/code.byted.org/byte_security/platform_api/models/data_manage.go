package models

import (
	"code.byted.org/byte_security/platform_api/utils/base/enum"
	"code.byted.org/gopkg/pkg/errors"
	"strings"
)

//请求单条参数结构
type GetPar struct {
	SourceID int `form:"source_id" binding:"required"`
}

//拉去全部参数结构
type ListAll struct {
	ProductID      int `form:"product_id" binding:"required"`
	ServiceType    int `form:"service_type"`
	DataSourceType int `form:"datasource_type"`
	Page           int `form:"page"`
	PageSize       int `form:"page_size"`
}

type DataSource struct {
	ID             uint                `json:"id"`
	SourceIds      []int               `json:"source_ids"`
	ProductID      int                 `json:"product_id"`
	ServiceType    enum.ServiceType    `json:"service_type"`
	SourceName     string              `json:"source_name"`
	DataSourceType enum.DataSourceType `json:"data_source_type"`

	//数据源链接信息与资源定位信息,tos,kafka,mysql
	ConnectToken   map[string]string `json:"connect_token"`
	LocationDetail map[string]string `json:"location_detail"`
	Page           int               `json:"page"`
	PageSize       int               `json:"page_size"`
}

//新增任务表单验证
func (d *DataSource) CheckAll() error {
	d.purify()
	errLegal := d.isLegal()
	if errLegal != nil {
		return errLegal
	}
	errIds := d.CheckIds()
	if errIds != nil {
		return errIds
	}
	return nil
}

//新增任务表单验证
func (d *DataSource) CheckIds() error {
	for i := range d.SourceIds {
		if d.SourceIds[i] <= 0 {
			return errors.New("ID列表不合法")
		}
	}
	return nil
}

//去除字符串空格
func (d *DataSource) purify() {
	strings.TrimSpace(d.SourceName)
}

//判断表单合法性
func (d *DataSource) isLegal() error {
	if !d.DataSourceType.Check() {
		return errors.New("数据用途参数值输入有误")
	} else if d.SourceName == "" {
		return errors.New("数据名称不能为空")
	} else if d.DataSourceType.Check() {
		return errors.New("数据类型值输入有误")
	} else if len(d.LocationDetail) == 0 {
		return errors.New("数据源属性定义为空")
	} else if d.ID == 0 {
		return errors.New("数据源ID不能为空")
	}
	return nil
}
