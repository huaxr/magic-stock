// @Contact:    huaxinrui
// @Time:       2019/9/8 下午3:31

package soc

import (
	"code.byted.org/gopkg/gorm"
	"encoding/json"
	"log"
)

type Structer interface {

}

func struct2Json(form Structer) (string, error) {
	str, err := json.Marshal(form)
	if err != nil {
		return "", err
	}
	return string(str), nil
}


func GetAssetDetail(db *gorm.DB, from, asset string) (string, error){
	var result string
	var err error
	switch from {
	case "WAF":
		var domain Domain
		db.Model(&Domain{}).Where("name = ?", asset).Find(&domain)
		result, err =  struct2Json(domain)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "HIDS":
		var host Host
		db.Model(&Host{}).Where("ip = ?", asset).Find(&host)
		result, err =  struct2Json(host)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	return result, nil
}