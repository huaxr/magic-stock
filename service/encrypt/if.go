// @Contact:    huaxinrui
// @Time:       2019/9/26 上午10:15

package encrypt

type EncryptIF interface {
	Decrypt(secret string) (string, error)
	Encrypt(plain string) (string, error)
}

var MD5Client EncryptIF = InitMd5Encrypt()
var SHA256Client EncryptIF = InitSha256Encrypt()
