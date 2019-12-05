package gokms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Response struct
type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Extra   map[string]interface{} `json:"extra"`
}

// Method that sends request to kms platform
func request(path, form string) (Response, error) {
	var r Response
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/%s", GetHost(), path), strings.NewReader(form))
	if err != nil {
		return r, err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return r, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return r, err
	}
	defer rsp.Body.Close()
	//fmt.Println(string(body))
	err = json.Unmarshal(body, &r)
	return r, err
}
