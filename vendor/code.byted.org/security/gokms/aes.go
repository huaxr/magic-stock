package gokms

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"crypto/sha1"
	"code.byted.org/gopkg/logs"
)

type AesCipher struct {
	key []byte
}

func NewAesCipher(dataKey string) (*AesCipher, error) {
	if len(dataKey) == 0 {
		logs.Error("[-] kms get nil data key")
		return nil, errors.New("kms get nil data key")
	}
	return &AesCipher{
		key: []byte(dataKey),
	}, nil
}

func (c AesCipher) Encrypt(origData string) (string, error) {
	orig := []byte(origData)
	result, err := c.aesEncrypt(orig, true)
	if err != nil {
		return "", err
	}else{
		return base64.StdEncoding.EncodeToString(result), nil
	}
}

func (c AesCipher) EncryptForQuery(origData string) (string, error) {
	orig := []byte(origData)
	s, e := c.aesEncrypt(orig, false)
	if e != nil {
		return fmt.Sprintf("@%s", s), e
	}else{
		return fmt.Sprintf("@%s", base64.StdEncoding.EncodeToString(s)), nil
	}
}

func (c AesCipher) EncryptGCM(origData string) (string, error) {
	return c.Encrypt(origData)
}

/*
func (c AesCipher) Encrypt(origData string) {
	return c.encrypt(origData, true, "GCM")
}
*/

// 加密
func (c AesCipher) aesEncrypt(orig []byte, ivFlag bool) ([]byte, error) {
	var iv []byte
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return []byte{}, err
	}
	var blockSize int
	blockSize = 12
	orig = c.pKCS5Padding(orig, block.BlockSize())
	//orig = c.zeroPadding(orig, block.BlockSize())
	if ivFlag == true {
		iv = make([]byte, blockSize)
		_, err := rand.Read(iv)
		if err != nil {
			return []byte{}, err
		}
	} else {
		if len(c.key) < blockSize {
			return []byte{}, errors.New("[kms] length of key smaller than block size, please check key")
		}
		iv = c.key[:blockSize]
	}
	crypted := make([]byte, len(orig))
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	crypted = aesgcm.Seal(nil, iv, orig, nil)
	if ivFlag == true {
		crypted = append(iv, crypted...)
	}
	return crypted, nil
}

func (c AesCipher) Decrypt(crypted string) (string, error) {
	var cryptedData string
	if strings.HasPrefix(crypted, "@") {
		cryptedData = strings.TrimPrefix(crypted, "@")
	}else{
		cryptedData = crypted
	}
	cryText, err := base64.StdEncoding.DecodeString(cryptedData)
	if err != nil {
		return "", err
	}
	var result []byte
	if strings.HasPrefix(crypted, "@") {
		result, err = c.aesDecrypt(cryText, false)
	}else {
		result, err = c.aesDecrypt(cryText, true)
	}
	return string(result), err
}

func (c AesCipher) DecryptGCM(crypted string) (string, error) {
	return c.Decrypt(crypted)
}

// 解密
func (c AesCipher) aesDecrypt(crypted []byte, ivFlag bool) ([]byte, error) {
	var iv []byte
	var cryBlob []byte
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return []byte{}, err
	}
	var blockSize int
	blockSize = 12
	if ivFlag == true {
		if len(crypted) < blockSize {
			return []byte{}, errors.New("[kms] length of decrypted data smaller than block size, please check data")
		}
		iv = crypted[0:blockSize]
		cryBlob = crypted[blockSize:]
	} else {
		if len(c.key) < blockSize {
			return []byte{}, errors.New("[kms] length of key smaller than block size, please check key")
		}
		iv = c.key[:blockSize]
		cryBlob = crypted
	}
	origData := make([]byte, len(cryBlob))
	//nonce := make([]byte, 12)
	if len(cryBlob) < 12 {
		return []byte{}, errors.New("crypto/cipher: input error")
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	origData, err = aesgcm.Open(nil, iv, cryBlob, nil)
	if err != nil {
		return []byte{}, err
	}
	origData, err = c.pKCS5UnPadding(origData)
	//origData, err = c.zeroUnPadding(origData)
	return origData, err
}

func (c AesCipher) zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func (c AesCipher) zeroUnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return nil, errors.New("error ciphertext found")
	}
	return origData[:(length - unpadding)], nil
}

func (c AesCipher) pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (c AesCipher) pKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	if length < unpadding {
		return nil, errors.New("error ciphertext found")
	}
	return origData[:(length - unpadding)], nil
}

func (c AesCipher) ID() string {
	if len(c.key) == 0 {
		return "<nil>"
	}
	h := sha1.New()
	h.Write(c.key)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (c AesCipher) String() string {
	return fmt.Sprintf("kms.aes, key: %s", c.ID())
}




func (c AesCipher) EncryptBin(orig []byte) ([]byte, error){
	return c.aesEncrypt(orig, true)
}

func (c AesCipher) DecryptBin(crypted []byte) ([]byte, error){
	return c.aesDecrypt(crypted, true)
}