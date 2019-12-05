// @Contact:    huaxinrui
// @Time:       2019/9/24 下午5:00

package upload

import (
	"context"

	"code.byted.org/gopkg/tos"

	"magic/stock/utils"
)

type UploadIF interface {
	NewTos() (*tos.Tos, error)
	UploadByFilePath(file_ath string) error
	UploadByBytes(name string, _bytes []byte) (string, error)
	GetFullPath() string
	genHashName(name string) string
}

const (
	AccessKey    = ""
	BOEAccessKey = ""
	BucketName   = ""

	TPath = ""
	BPath = ""
)

var GlobalUpload UploadIF

func init() {
	tmp := new(Upload)
	tmp.context = context.TODO()
	switch utils.TellEnv() {
	case "loc":
		tmp.accessKey = BOEAccessKey
		tmp.path = BPath
	case "boe":
		tmp.accessKey = BOEAccessKey
		tmp.path = BPath
	case "tce":
		tmp.path = TPath
		tmp.accessKey = AccessKey
	}
	tmp.bucketName = BucketName
	GlobalUpload = tmp
}

type Upload struct {
	context    context.Context
	accessKey  string
	bucketName string
	path       string
}
