package gokms

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"code.byted.org/gopkg/env"
)

// Generate signature
func GetSignature(secretId, secretKey, action string, timestamp int64) string {
	if secretId == "" || secretKey == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(fmt.Sprintf("%s?id=%s&timestamp=%d", action, secretId, timestamp)))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Convert struct to json string
func StructToJson(form ParamForm) (string, error) {
	str, err := json.Marshal(form)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func isBoe() bool {
	//dc := os.Getenv("RUNTIME_IDC_NAME")
	force := os.Getenv("DEBUG_IDC")
	return env.IsBoe() || force == "boe"
}

func isDebug() bool {
	dc := os.Getenv("KMS_ENV")
	force := os.Getenv("DEBUG_IDC")
	return dc == "debug" || force == "local"
}

func isVa() bool {
	//zone := os.Getenv("TCE_ZONE")
	force := os.Getenv("KMS_ZONE")
	return env.Region() == env.R_US || env.Region() == env.R_CA || force == "va"
}

func isSg() bool {
	//zone := os.Getenv("TCE_ZONE")
	force := os.Getenv("KMS_ZONE")
	return env.Region() == env.R_SG || env.Region() == env.R_ALISG || force == "sg"
}

func GetHost() string {
	if isBoe() {
		return BoeHost
	} else if isDebug() {
		return DebugHost
	}
	force := os.Getenv("KMS_ZONE")
	if force != "cn" {
		if isVa() {
			return VaHost
		}else if isSg() {
			return SgHost
		}
	}
	return Host
}
