package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"code.byted.org/inf/infsecc"
)

var (
	// 获取用户组信息
	GET      = "GET"
	POST     = "POST"
	URL      = "https://ee.byted.org/ratak/employees/%s/groups/"
	GROUPURL = "https://open.byted.org/people/employee/?email=%s@bytedance.com"
	key      = "111:59A1D19D5BFF46C7BFED446C0B5E5F66"
	// 内容开放平台
	GET_HOST_TAGS  = "http://console.byted.org/tag/api/v1/host/tags/"
	GET_HOST_OWNER = "http://console.byted.org/tag/api/v1/host/owner/"
	GET_TAG_HOSTS  = "http://console.byted.org/tag/api/v1/tag/hosts/"
	ACCESSKEY      = "IFQP2FAGGH"
	ACCESSSECRET   = "OTZRVIX7DD31KU1PH3ZRZBEPV02CBR324AC7ATEJNJHS398TNJ"
	// 服务树
	ServiceTree = "http://galaxy-api.bytedance.net/service_meta/api/v2/nodes/?name=" // security.soc.online
	// PSM 获取相关信息
	//PSM = "curl http://tce.byted.org/api/v3/3rd/services/get/service_detail/\?psm\=security.soc.online
	//PSM2 = "curl http://tsearch.byted.org/search/\?keyword\=security.soc.online"
)

// Tag 相关接口， 根据指定的 ip
func GetTagHost(ip string, kind string) []byte {
	client := &http.Client{}
	URL := ""
	data := map[string]string{}

	if kind == "tags" {
		URL = GET_HOST_TAGS
		data = map[string]string{
			"hosts": ip,
			//"tags": "system.security",
		}
	} else if kind == "owners" {
		URL = GET_HOST_OWNER
		data = map[string]string{
			"hosts": ip,
			//"tags": "system.security",
		}
	} else if kind == "hosts" {
		URL = GET_TAG_HOSTS
		data = map[string]string{
			"tags": ip,
		}
	}
	dataStr, err := json.Marshal(data)
	req, err := http.NewRequest(POST, URL, bytes.NewReader(dataStr))
	req.SetBasicAuth(ACCESSKEY, ACCESSSECRET)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	bodyByte, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return bodyByte
}

// Get group info by given name
var result = make(chan map[string]string, 1)

func WatchGroup(name string) {
	client := &http.Client{}
	req, _ := http.NewRequest(GET, fmt.Sprintf(URL, name), nil)
	encodeString := base64.StdEncoding.EncodeToString([]byte(key))
	req.Header.Add("Authorization", "Basic "+encodeString)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]string)
	err := json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	result <- m
}

func GetGroupByName(ctx context.Context, name string) map[string]string {
	go WatchGroup(name)

	select {
	case x := <-result:
		return x

	case <-ctx.Done():
		//fmt.Println(ctx.Err()) // prints "context deadline exceeded"
		return nil
	}
}

// 以下接口只返回指定group
var result2 = make(chan map[string]interface{}, 1)

func WatchGroup2(name string) {
	type group_query struct {
		Employees []map[string]interface{}
		Success   bool
	}
	var grouper group_query

	client := &http.Client{}
	req, _ := http.NewRequest(GET, fmt.Sprintf(GROUPURL, name), nil)
	req.Header.Add("Authorization", "Basic MTExOjU5QTFEMTlENUJGRjQ2QzdCRkVENDQ2QzBCNUU1RjY2")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(body, &grouper)
	if err != nil {
		panic(err)
	}
	if grouper.Success {
		result2 <- grouper.Employees[0]
	}
}

func GetGroupByName2(ctx context.Context, name string) map[string]interface{} {
	go WatchGroup2(name)

	select {
	case x := <-result2:
		return x

	case <-ctx.Done():
		//fmt.Println(ctx.Err()) // prints "context deadline exceeded"
		return nil
	}
}

// 根据sdk获取group， 在boe,tce下使用
var result3 = make(chan map[string]interface{}, 1)

type group_query struct {
	Employees []map[string]interface{}
	Success   bool
}

func WatchGroup3(name string) {
	var grouper group_query
	client := &http.Client{}
	req, _ := http.NewRequest(GET, fmt.Sprintf(GROUPURL, name), nil)
	req.Header.Add("Authorization", "Basic MTExOjU5QTFEMTlENUJGRjQ2QzdCRkVENDQ2QzBCNUU1RjY2")
	token, err := infsecc.GetToken(false)
	if err == nil {
		req.Header.Add("X-Dps-Token", token)
	}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &grouper)
	if err != nil {
		panic(err)
	}
	if grouper.Success {
		result3 <- grouper.Employees[0]
	}
}

func GetGroupByName3(ctx context.Context, name string) map[string]interface{} {
	go WatchGroup3(name)

	select {
	case x := <-result3:
		return x

	case <-ctx.Done():
		//fmt.Println(ctx.Err()) // prints "context deadline exceeded"
		return nil
	}
}

// 根据域名获取psm
func GetPSMByDomain(domain string) string {
	command := "curl -I -H'Get-Svc: 1' " + domain
	s := System(command)
	reg := regexp.MustCompile(`\nX-Svc: (.*?)\n`)
	params := reg.FindAllString(s, 1)
	if len(params) > 0 {
		return strings.Split(params[0], " ")[1]
	}
	return ""
}

// 根据psm获取负责人 （根据服务树）
type RS struct {
	ErrorCode int `json:"error_code"`
	Data      []struct {
		Owners    string `json:"owners"`
		CreatedBy string `json:"created_by"`
		Path      string `json:"path"`
	}
}

//func Cache(f func (psm string) RS, psm string) RS{
//	r, _ := common.Backend.Redis.Get(psm).Result()
//	if r != "" {
//		log.Println("cache redis found")
//		var res RS
//		json.Unmarshal([]byte(r), &res)
//		return res
//	} else {
//		log.Println("no cache redis found")
//	}
//
//	rs := f(psm)
//
//	defer func(rs RS) {
//		t, _ := json.Marshal(rs)
//		err := common.Backend.Redis.Set(psm, string(t), 0).Err()
//		if err != nil {
//			panic(err)
//		}
//	}(rs)
//
//	return rs
//}
//
//func GetPsmUser(psm string) RS {
//	pc, _, _, _ := runtime.Caller(1)
//	caller := runtime.FuncForPC(pc).Name()
//	if !strings.HasSuffix(caller, "api.Cache") {
//		log.Println("Suggest using cache to call this function. e.g. Cache(GetPsmUser, psm)")
//	}
//
//	_, body, _ := gorequest.New().Get(ServiceTree + psm).End()
//
//	var res RS
//	json.Unmarshal([]byte(body), &res)
//
//	if res.ErrorCode != 0 {
//		return RS{}
//	} else {
//		return res
//	}
//
//}
