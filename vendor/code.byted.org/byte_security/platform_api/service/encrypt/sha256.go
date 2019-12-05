// @Contact:    huaxinrui
// @Time:       2019/11/5 下午7:26

package encrypt

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

type Sha256Client struct {
}

func (k *Sha256Client) Decrypt(secret string) (string, error) {
	return "", errors.New("md5 hash can not decrypt ! please user rainbow table")
}

func (k *Sha256Client) Encrypt(plain string) (string, error) {
	h := sha256.New()
	h.Write([]byte(plain))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}

func InitSha256Encrypt() *Sha256Client {
	k := new(Sha256Client)
	return k
}
