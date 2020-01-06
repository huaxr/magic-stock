// @Contact:    huaxinrui
// @Time:       2019/10/15 下午5:57

package jwt

import (
	"crypto"
)

type JwtIF interface {
	GetPublicKeyStr() string
	LoadKeyBytes()
	GenRsaKey(bits int) error
	RsaSign(data BYTES) (BYTES, error)
	RsaSignVer(data BYTES, signature BYTES) error
	// 公钥加密, 私钥解密
	RsaEncrypt(data BYTES) (BYTES, error)
	RsaDecrypt(cipher BYTES) (BYTES, error)
	GetEnv() bool
	GenHashPrefix()
	// 私钥加密, 公钥解密
	PrivateEncrypt(data BYTES) (string, error)
	PublicDecrypt(data string) (string, error)
}

type BYTES []byte

const (
	JWToken = "BS-TOKEN"
	PRIVATE = "conf/x.pem"
	PUBLIC  = "conf/t.pem"
)

var Jwt JwtIF

func init() {
	//tmp := new(jwt)
	//if utils.TellEnv() == "loc" {
	//	tmp.loc = true
	//}
	//Jwt = tmp
	//tmp.GenHashPrefix()
	//tmp.GetPublicKeyStr()
}

type jwt struct {
	loc                   bool
	privateKey, publicKey BYTES
	hashPrefixes          map[crypto.Hash]BYTES
}
