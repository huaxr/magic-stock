package gokms

import (
	"code.byted.org/inf/infsecc"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"code.byted.org/gopkg/logs"
)

const (
	Host      = "https://kms.bytedance.net"
	VaHost    = "https://kms-va.bytedance.net"
	SgHost    = "https://kms-sg.bytedance.net"
	DebugHost = "http://127.0.0.1:13888"
	BoeHost   = "https://kms-boe.bytedance.net"
)

// Kms Client
type KMSClient struct {
	secretId  string
	secretKey string
	PSM       string
	Retries   int
}

// Method to create kms client
func NewKMSClient(PSM string) (*KMSClient, error) {
	logs.Info(fmt.Sprintf("Choose kms host: %s", GetHost()))
	return &KMSClient{
		secretId:  os.Getenv("kms_id"),
		secretKey: os.Getenv("kms_key"),
		PSM:       PSM,
		Retries:   3,
	}, nil
}

// Method to get tce token
func getToken() (token string) {
	token, err := infsecc.GetToken(false)
	if err != nil {
		//log.Printf("Get tce token err: %v\n", err)
		logs.Warnf("Get tce token err: %v\n", err)
		return ""
	} else {
		//fmt.Println(infsecc.ParseToken(token))
		logs.Info("Get tce token success\n")
		return token
	}
}

// Method that check struct
func (c KMSClient) checkStruct() bool {
	if c.PSM == "" {
		return false
	}
	return true
}

// Method that send request
func (c KMSClient) request(path, form string) (Response, error) {
	var err error
	var rsp Response
	if c.checkStruct() == false {
		return Response{}, errors.New("client value incorrect, please 'use NewKMSClient(PSM string)' to create client")
	}
	retries := c.Retries
	for retries > 0 {
		rsp, err = request(path, form)
		if err != nil {
			retries -= 1
			logs.Warn(fmt.Sprintf("kms request server error: %s, request %d times", err.Error(), c.Retries - retries))
		}else{
			break
		}
	}
	if rsp.Code == 0 {
		return rsp, err
	}
	return rsp, errors.New(fmt.Sprintf("Error: code: %d, message: %s", rsp.Code, rsp.Message))
}

// Method That generate public form
func (c KMSClient) getPublicParams(action string) PublicForm {
	timestamp := time.Now()
	return PublicForm{
		SecretId:  c.secretId,
		Signature: GetSignature(c.secretId, c.secretKey, action, timestamp.Unix()),
		Timestamp: timestamp.Unix(),
		PSM:       c.PSM,
		Token:     getToken(),
	}
}

// Encrypt method, return text that after encrypting
func (c KMSClient) Encrypt(keyid, plaintext string) (string, error) {
	form := EncryptForm{
		PublicForm: c.getPublicParams("Encrypt"),
		Keyid:      keyid,
		Plaintext:  plaintext,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("encrypt", params)
	if err != nil {
		return "", err
	}
	return rsp.Extra["text"].(string), err
}

// Decrypt method, return text that after decrypting
func (c KMSClient) Decrypt(text string) (string, error) {
	form := DecryptForm{
		PublicForm:     c.getPublicParams("Decrypt"),
		CiphertextBlob: text,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("decrypt", params)
	if err != nil {
		logs.Error(fmt.Sprintf("[-] kms cloud decrypt error: %s", err.Error()))
		return "", err
	}
	return rsp.Extra["text"].(string), err
}

// Batch Decrypt method, return texts that after decrypting
func (c KMSClient) DecryptBatch(texts ...string) ([]string, error) {
	if texts == nil {
		return []string{}, nil
	}
	form := BatchDecryptForm{
		PublicForm:  c.getPublicParams("Decrypt"),
		Ciphertexts: texts,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("decrypt/batch", params)
	if err != nil {
		return []string{}, err
	}
	interfaceTexts := rsp.Extra["texts"].([]interface{})
	var stringTexts []string
	for _, text := range interfaceTexts {
		stringTexts = append(stringTexts, text.(string))
	}
	return stringTexts, err
}

// Create master key method, return new key's keyid
func (c KMSClient) CreateMasterKey(alias, description string) (string, error) {
	form := NewKeyForm{
		PublicForm:  c.getPublicParams("CreateKey"),
		Alias:       alias,
		Description: description,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("key/create", params)
	if err != nil {
		return "", err
	}
	return rsp.Extra["keyid"].(string), err
}

// Method that list all master key, return the list of keyinfo struct
func (c KMSClient) ListMasterKeys() ([]KeyInfo, error) {
	var keys []KeyInfo
	form := ListKeyForm{
		PublicForm: c.getPublicParams("ListKey"),
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("key/list", params)
	if err != nil {
		return keys, err
	}
	data := rsp.Extra["data"].([]interface{})
	for _, value := range data {
		v := value.(map[string]interface{})
		keys = append(keys, KeyInfo{
			Alias:       v["alias"].(string),
			Createdtime: v["createdtime"].(string),
			Description: v["description"].(string),
			Keyid:       v["keyid"].(string),
			Status:      v["status"].(string),
		})
	}
	return keys, err
}

// Method that get info of a master key, return the list of keyinfo struct
func (c KMSClient) GetMasterKey(keyid string) (KeyInfo, error) {
	form := InfoKeyForm{
		PublicForm: c.getPublicParams("GetKey"),
		Keyid:      keyid,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("key/info", params)
	if err != nil {
		return KeyInfo{}, err
	}
	return KeyInfo{
		Alias:       rsp.Extra["alias"].(string),
		Createdtime: rsp.Extra["createdtime"].(string),
		Description: rsp.Extra["description"].(string),
		Keyid:       rsp.Extra["keyid"].(string),
		Status:      rsp.Extra["status"].(string),
	}, err
}

// Method that create data key, return the content of data key
func (c KMSClient) GenerateDataKey(keyid, name, psm string) (string, error) {
	form := GenDataKeyForm{
		PublicForm: c.getPublicParams("GenerateDataKey"),
		Keyid:      keyid,
		Name:       name,
		P:          psm,
		Type:       "AES_256",
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("datakey/create", params)
	if err != nil {
		return "", err
	}
	return rsp.Extra["id"].(string), err
}

// Method that used to decrypt data key
func (c KMSClient) DecryptDataKey(text string) (*AesCipher, error) {
	if text == "" {
		return c.GetDataKey()
	} else if len(text) == 6 && !strings.Contains(text, ".") {
		return c.DecryptDataKeyById(text)
	} else if len(text) < 200 {
		return c.DecryptDataKeyByPSM(text)
	}
	form := DecryptDataKeyForm{
		PublicForm:     c.getPublicParams("Decrypt"),
		CiphertextBlob: text,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("decrypt/datakey", params)
	if err != nil {
		logs.Error(fmt.Sprintf("[-] kms get data key error: %s", err.Error()))
		return nil, err
	}
	dk := rsp.Extra["text"].(string)
	return NewAesCipher(dk)
}

// Method that decrypt data key by id
func (c KMSClient) DecryptDataKeyById(keyid string) (*AesCipher, error) {
	form := DecryptDataKeyByIdForm{
		PublicForm: c.getPublicParams("Decrypt"),
		Keyid:      keyid,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("decrypt/datakey/id", params)
	if err != nil {
		logs.Error(fmt.Sprintf("[-] kms get data key by id error: %s", err.Error()))
		return nil, err
	}
	dk := rsp.Extra["text"].(string)
	return NewAesCipher(dk)
}

// Method that decrypt data key by psm
func (c KMSClient) DecryptDataKeyByPSM(psm string) (*AesCipher, error) {
	if psm == "" {
		return c.GetDataKey()
	}
	form := DecryptDataKeyByPSMForm{
		PublicForm: c.getPublicParams("Decrypt"),
		P:          psm,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("decrypt/datakey/psm", params)
	if err != nil {
		logs.Error(fmt.Sprintf("[-] kms get data key by psm error: %s", err.Error()))
		return nil, err
	}
	dk := rsp.Extra["text"].(string)
	return NewAesCipher(dk)
}

// Method that decrypt data key by psm
func (c KMSClient) GetDataKey() (*AesCipher, error) {
	form := DecryptDataKeyByPSMForm{
		PublicForm: c.getPublicParams("Decrypt"),
		P:          "",
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("decrypt/datakey/psm", params)
	if err != nil {
		logs.Error(fmt.Sprintf("[-] kms get data key error: %s", err.Error()))
		return nil, err
	}
	dk := rsp.Extra["text"].(string)
	return NewAesCipher(dk)
}

// Method that share data key
func (c KMSClient) ShareDataKey(psm, description string) (string, error) {
	form := ShareDataKeyForm{
		PublicForm:  c.getPublicParams("ShareDataKey"),
		P:           psm,
		Description: description,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("datakey/share", params)
	if err != nil {
		return "", err
	}
	return rsp.Extra["text"].(string), err
}

// Method that share master key
func (c KMSClient) ShareMasterKey(text, username string) (string, error) {
	form := ShareMasterKeyForm{
		PublicForm: c.getPublicParams("ShareMasterKey"),
		Text:       text,
		Username:   username,
	}
	params, _ := StructToJson(form)
	rsp, err := c.request("key/share", params)
	if err != nil {
		return "", err
	}
	return rsp.Extra["text"].(string), err
}
