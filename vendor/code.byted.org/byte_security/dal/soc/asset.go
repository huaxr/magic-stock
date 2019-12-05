package soc

import "code.byted.org/byte_security/dal/seal"

type Asset interface {
	GetAssetKey() string
	GetAssetValue() string
	New() interface{}
}

func GetAssetStructByName(tableName string) interface{} {
	switch tableName {
	case "app", "byte_security_asset_app":
		return App{}
	case "domain", "byte_security_asset_domain":
		return Domain{}
	case "device", "byte_security_seal_device":
		return seal.Device{}
	case "host", "byte_security_asset_host":
		return Host{}
	case "idc", "byte_security_asset_idc":
		return IDC{}
	case "network", "byte_security_asset_network":
		return NetWork{}
	case "product", "byte_security_asset_product":
		return Product{}
	case "psm", "byte_security_asset_psm":
		return PSM{}
	case "repo", "byte_security_asset_repo":
		return Repo{}
	case "website", "byte_security_asset_website":
		return Website{}
	default:
		return Host{}
	}
}
