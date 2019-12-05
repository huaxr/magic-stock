// @Contact:    huaxinrui
// @Time:       2019/9/24 下午5:01

package upload

import (
	"bytes"
	"io/ioutil"
	"path"
	"strconv"
	"time"

	"code.byted.org/byte_security/platform_api/service/check"
	"code.byted.org/byte_security/platform_api/service/encrypt"
	"code.byted.org/gopkg/tos"
)

func (u *Upload) NewTos() (*tos.Tos, error) {
	bucket := tos.WithAuth(u.bucketName, u.accessKey)
	return tos.NewTos(bucket)
}

func (u *Upload) UploadByFilePath(file_path string) error {
	client, err := u.NewTos()
	if err != nil {
		return err
	}
	file_name := u.genHashName(path.Base(file_path)) + path.Ext(file_path)
	byt, err := ioutil.ReadFile(file_path)
	if err != nil {
		return err
	}
	if err := client.PutObject(u.context, file_name, int64(len(byt)), bytes.NewBuffer(byt)); err != nil {
		return err
	}
	return nil
}

func (u *Upload) UploadByBytes(name string, _bytes []byte) (string, error) {
	err := check.Security.PreventFileAnyUpload(path.Ext(name))
	if err != nil {
		return "", err
	}
	client, err := u.NewTos()
	file_name := u.genHashName(name) + path.Ext(name)
	if err != nil {
		return "", err
	}
	if err := client.PutObject(u.context, file_name, int64(len(_bytes)), bytes.NewBuffer(_bytes)); err != nil {
		return "", err
	}
	return u.GetFullPath() + file_name, nil
}

func (u *Upload) GetFullPath() string {
	return u.path
}

func (u *Upload) genHashName(name string) string {
	md5, _ := encrypt.MD5Client.Encrypt(name)
	unix := strconv.FormatInt(time.Now().Unix(), 10)
	return md5 + unix
}
