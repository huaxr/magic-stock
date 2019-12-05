// @Contact:    huaxinrui
// @Time:       2019/10/15 下午5:53

package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

func (a *jwt) GetPublicKeyStr() string {
	if a.publicKey == nil {
		a.LoadKeyBytes()
	}
	return string(a.publicKey)
}

func (a *jwt) GetEnv() bool {
	return false
}

func (a *jwt) LoadKeyBytesCycle() {
J:
	if a.privateKey == nil || a.publicKey == nil {
		var err error
		a.publicKey, err = ioutil.ReadFile(PUBLIC)
		a.privateKey, err = ioutil.ReadFile(PRIVATE)
		if err != nil {
			fmt.Println(err)
			a.GenRsaKey(1024)
			goto J
		}
	}
}

func (a *jwt) LoadKeyBytes() {
	if a.privateKey == nil || a.publicKey == nil {
		var err error
		a.publicKey, err = ioutil.ReadFile(PUBLIC)
		a.privateKey, err = ioutil.ReadFile(PRIVATE)
		if err != nil {
			panic(err)
		}
	}
}

//私钥签名
func (a *jwt) RsaSign(data BYTES) (BYTES, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	//获取私钥
	block, _ := pem.Decode(a.privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed)
}

//公钥验证
func (a *jwt) RsaSignVer(data BYTES, signature BYTES) error {
	hashed := sha256.Sum256(data)
	block, _ := pem.Decode(a.publicKey)
	if block == nil {
		return errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//验证签名
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}

// 公钥加密
func (a *jwt) RsaEncrypt(data BYTES) (BYTES, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(a.publicKey)
	if block == nil {
		return nil, errors.New("public key error.")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// 私钥解密
func (a *jwt) RsaDecrypt(ciphertext BYTES) (BYTES, error) {
	//获取私钥
	block, _ := pem.Decode(a.privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func (a *jwt) PrivateEncrypt(data BYTES) (string, error) {
	block, _ := pem.Decode(a.privateKey)
	if block == nil {
		panic("decode error")
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	signData, err := rsa.SignPKCS1v15(nil, private, crypto.Hash(0), data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signData), nil
}

func (a *jwt) PublicDecrypt(datas string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(datas)
	if err != nil {
		log.Println(datas, err)
		return "", errors.New("BASE64 解密错误")
	}
	block, _ := pem.Decode(a.publicKey)
	if block == nil {
		return "", errors.New("公钥存在错误") // panic
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub := pubInterface.(*rsa.PublicKey)
	decData, err := a.publicDecrypt(pub, crypto.Hash(0), nil, data)
	if err != nil {
		return "", err
	}
	return string(decData), nil
}
