// @Contact:    huaxinrui
// @Time:       2019/10/25 下午5:49

package encrypt

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type Md5Client struct {
}

func (k *Md5Client) Decrypt(secret string) (string, error) {
	return "", errors.New("md5 hash can not decrypt ! please user rainbow table")
}

func (k *Md5Client) Encrypt(plain string) (string, error) {
	data := []byte(plain)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str, nil
}

func InitMd5Encrypt() *Md5Client {
	k := new(Md5Client)
	return k
}
