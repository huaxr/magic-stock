// @Contact:    huaxinrui
// @Time:       2019/10/15 下午6:04

package jwt

import (
	"fmt"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	err := Jwt.GenRsaKey(1024)
	if err != nil {
		fmt.Println(err)
	}
}

func TestG(t *testing.T) {
	Jwt.LoadKeyBytes()
	var theMsg = "你好 .世界"
	fmt.Println("Source:", theMsg)
	//私钥签名
	sig, _ := Jwt.RsaSign([]byte(theMsg))
	fmt.Println(string(sig))
	//公钥验证
	fmt.Println(Jwt.RsaSignVer([]byte(theMsg), sig))

	enc, _ := Jwt.RsaEncrypt([]byte(theMsg))
	fmt.Println("Encrypted:", string(enc))
	//私钥解密
	decstr, _ := Jwt.RsaDecrypt(enc)
	fmt.Println("Decrypted:", string(decstr))
}
