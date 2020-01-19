// @Time:       2020/1/19 下午2:17

package control

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"magic/stock/model"
	mathRand "math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	MemoryCacheVar  *MemoryCache
	AppID                  = "wx921b6afd43dddd8e"
	AppSecret              = "247a017acef5b2a65af1854d2ae4a950"
	AccessTokenHost string = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + AppID + "&secret=" + AppSecret
	JsAPITicketHost string = "https://api.weixin.qq.com/cgi-bin/ticket/getticket"
)

func init() {
	MemoryCacheVar = new(MemoryCache)
	MemoryCacheVar.Items = make(map[string]*Item)
}

func (d *UserControl) GetWxSign(c *gin.Context) {
	var (
		noncestr, jsapi_ticket, timestamp, url, signature, signatureStr, access_token string
		wxAccessToken                                                                 WxAccessToken
		wxJsApiTicket                                                                 WxJsApiTicket
		wxSignature                                                                   WxSignature
	)

	_auth, _ := c.Get("auth")
	authentication := _auth.(*model.AuthResult)

	user, _ := UserControlGlobal.Query("id = ?", []interface{}{authentication.Uid})
	url = c.DefaultQuery("url", "")
	if url == "" {
		c.JSON(200, gin.H{"error_code": 1, "err_msg": "没有指定的url参数", "data": nil})
		return
	}
	if strings.Contains(url, "?") {
		url += "&token=" + user.ShareToken
	} else {
		url += "?token=" + user.ShareToken
	}

	log.Println("分享的url:", url)
	noncestr = RandStringBytes(16)
	timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	//获取access_token，如果缓存中有，则直接取出数据使用；否则重新调用微信端接口获取
	client := &http.Client{}
	if MemoryCacheVar.Get("access_token") == nil {
		request, _ := http.NewRequest("GET", AccessTokenHost, nil)
		response, _ := client.Do(request)
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			c.JSON(200, gin.H{"error_code": 1, "err_msg": err.Error(), "data": nil})
			return
		}
		err = json.Unmarshal(body, &wxAccessToken)
		if err != nil {
			c.JSON(200, gin.H{"error_code": 1, "err_msg": err.Error(), "data": nil})
			return
		}
		if wxAccessToken.Errcode == 0 {
			access_token = wxAccessToken.Access_token
		} else {
			c.JSON(200, gin.H{"error_code": 1, "err_msg": wxAccessToken.Errmsg, "data": nil})
			return
		}
		MemoryCacheVar.Put("access_token", access_token, time.Duration(wxAccessToken.Expires_in)*time.Second)

		//获取 jsapi_ticket
		requestJs, _ := http.NewRequest("GET", JsAPITicketHost+"?access_token="+access_token+"&type=jsapi", nil)
		responseJs, _ := client.Do(requestJs)
		defer responseJs.Body.Close()
		bodyJs, err := ioutil.ReadAll(responseJs.Body)
		if err != nil {
			c.JSON(200, gin.H{"error_code": 1, "err_msg": err.Error(), "data": nil})
			return
		}
		err = json.Unmarshal(bodyJs, &wxJsApiTicket)
		if err != nil {
			c.JSON(200, gin.H{"error_code": 1, "err_msg": err.Error(), "data": nil})
			return
		}
		if wxJsApiTicket.Errcode == 0 {
			jsapi_ticket = wxJsApiTicket.Ticket
		} else {
			c.JSON(200, gin.H{"error_code": 1, "err_msg": wxJsApiTicket.Errmsg, "data": nil})
			return
		}
		MemoryCacheVar.Put("jsapi_ticket", jsapi_ticket, time.Duration(wxJsApiTicket.Expires_in)*time.Second)
	} else {
		//缓存中存在access_token，直接读取
		access_token = MemoryCacheVar.Get("access_token").(*Item).Value
		jsapi_ticket = MemoryCacheVar.Get("jsapi_ticket").(*Item).Value
	}
	log.Println("分享接口access_token:", access_token)
	log.Println("分享接口jsapi_ticket:", jsapi_ticket)

	// 获取 signature
	signatureStr = "jsapi_ticket=" + jsapi_ticket + "&noncestr=" + noncestr + "&timestamp=" + timestamp + "&url=" + url
	signature = GetSha1(signatureStr)
	log.Println("签名:", signatureStr, signature)
	wxSignature.Url = url
	wxSignature.Noncestr = noncestr
	wxSignature.Timestamp = timestamp
	wxSignature.Signature = signature
	wxSignature.AppID = AppID

	c.JSON(200, gin.H{"error_code": 0, "err_msg": nil, "data": gin.H{"url": url, "noncestr": noncestr, "timestamp": timestamp, "signature": signature, "appid": AppID}})

}

//生成指定长度的字符串
func RandStringBytes(n int) string {
	const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mathRand.Intn(len(letterBytes))]
	}
	return string(b)
}

//SHA1加密
func GetSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

type WxAccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}
type WxJsApiTicket struct {
	Ticket     string `json:"ticket"`
	Expires_in int    `json:"expires_in"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}
type WxSignature struct {
	Noncestr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Url       string `json:"url"`
	Signature string `json:"signature"`
	AppID     string `json:"appId"`
}

type WxSignRtn struct {
	ErrCode int         `json:"err_code"`
	ErrMsg  string      `json:"err_msg"`
	Data    WxSignature `json:"data"`
}

// 数据缓存处理
type Item struct {
	Value      string
	CreateTime time.Time
	LifeTime   time.Duration
}

type MemoryCache struct {
	sync.RWMutex
	Items map[string]*Item
}

func (mc *MemoryCache) Put(key string, value string, lifeTime time.Duration) {
	mc.Lock()
	defer mc.Unlock()
	mc.Items[key] = &Item{
		LifeTime:   lifeTime,
		Value:      value,
		CreateTime: time.Now(),
	}
}

func (mc *MemoryCache) Get(key string) interface{} {
	mc.RLock()
	defer mc.RUnlock()
	if e, ok := mc.Items[key]; ok {
		if !e.isExpire() {
			return e
		} else {
			delete(mc.Items, key)
		}
	}
	return nil
}

func (e *Item) isExpire() bool {
	if e.LifeTime == 0 {
		return false
	}
	//根据创建时间和生命周期判断元素是否失效
	return time.Now().Sub(e.CreateTime) > e.LifeTime

}
