package lark_bot

import (
	"fmt"
	"strings"
	"time"

	"code.byted.org/byte_security/platform_api/utils"
	"code.byted.org/gopkg/logs"
)

const (
	OnCallBotAPIID      = "cli_9dc3ff012afa910b"
	OnCallBotAPISecret  = "6Md9dpoKGwQqt3n98g77bcxrehQB3XNm"
	BOEBotAPIID         = "cli_9de22e46e33ed10b"
	BOEBotAPISecret     = "SCJq9hygjM3EzOcIrgMJfAo2LWUSzPnq"
	ByteSecBotAPIID     = "cli_9dc3ff012afa910b"
	ByteSecBotAPISecret = "6Md9dpoKGwQqt3n98g77bcxrehQB3XNm"
	AccessTokenURL      = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"
	EmailByOpenIDURL    = "https://open.feishu.cn/open-apis/user/v4/info?open_id="
	EmailByUserIDURL    = "https://open.feishu.cn/open-apis/user/v4/info?user_id="
	Email2UserIDURL     = "https://open.feishu.cn/open-apis/user/v4/email2id"
	CreateChatGroupURL  = "https://open.feishu.cn/open-apis/chat/v4/chatter/add/"
	SendMessageURL      = "https://open.feishu.cn/open-apis/message/v4/send/"
	GroupChatListURL    = "https://open.feishu.cn/open-apis/chat/v4/list"
)

var OnCallBot LarkBot = &BotService{}
var ByteSecBot LarkBot = &BotService{}

type LarkBot interface {
	Init(appID, appSecret string)
	GetEmailByOpenID(openID string) string
	GetEmailByUserID(userID string) string
	GetUserIDByEmail(userEmail string) string
	GetLarkGroupChatID() ([]*LarkGroupInfo, error)
	AddUser2Group(email string, chatID string) error
	SendMessageCard(receiver ReceiverInfo, content interface{}) error
	SetConsumerFunction(f func(interface{}))
	ProduceMsg(msg interface{})
	ConsumeMsg()
}

func (b *BotService) Init(appID, appSecret string) {
	b.AppId = &appID
	b.AppSecret = &appSecret
	b.Messages = make(chan interface{}, 100)
}

func (b *BotService) GetUserIDByEmail(userEmail string) string {
	err := b.getValidTenantAccessToken()
	if err != nil {
		logs.Error("【Lark_Bot】Send message failed——Cannot get userId by email, error: %v", err)
		return ""
	}
	var rsp Email2UserInfo
	if !strings.Contains(userEmail, "@") {
		userEmail += "@bytedance.com"
	}
	body := map[string]string{"email": userEmail}
	err = utils.HttpPostWithHeader(Email2UserIDURL, b.makeHeader(), body, &rsp)
	if err != nil {
		logs.Error("【Lark_Bot】Error in GetUserIDByEmail of %s: %v", userEmail, err)
	}
	return rsp.Data.UserId
}

func (b *BotService) GetEmailByOpenID(openID string) string {
	err := b.getValidTenantAccessToken()
	if err != nil {
		logs.Error("【Lark_Bot】 Cannot get tenant access token, error: %v", err)
		return ""
	}
	var rsp UserInfoRsp
	url := EmailByOpenIDURL + openID
	token := fmt.Sprintf("Bearer %s", b.TenantAccessToken.Token)
	err = utils.HttpGetWithToken(url, token, &rsp)
	if err != nil {
		logs.Error("【Lark_Bot】Error in GetChatIdByToken: %v", err)
		return ""
	}
	return rsp.Data.Email
}

func (b *BotService) GetEmailByUserID(userID string) string {
	err := b.getValidTenantAccessToken()
	if err != nil {
		logs.Error("【Lark_Bot】 Cannot get tenant access token, error: %v", err)
		return ""
	}
	var rsp UserInfoRsp
	url := EmailByUserIDURL + userID
	token := fmt.Sprintf("Bearer %s", b.TenantAccessToken.Token)
	err = utils.HttpGetWithToken(url, token, &rsp)
	if err != nil {
		logs.Error("【Lark_Bot】Error in GetChatIdByToken: %v", err)
		return ""
	}
	logs.Infof("【Lark_Bot】GetEmailByUserID done:%s-%s", userID, rsp.Data.Email)
	return rsp.Data.Email
}

func getTenantAccessToken(AppID string, AppSecret string) (TenantAccessTokenRsp, error) {
	var rsp TenantAccessTokenRsp
	params := map[string]string{"app_id": AppID, "app_secret": AppSecret}
	err := utils.HttpPost(AccessTokenURL, params, &rsp)
	if err != nil || rsp.Code != 0 {
		logs.Errorf("getTenantAccessToken failed:%v", err)
		return rsp, fmt.Errorf("【Lark_Bot】get tenantAccessToken responseCode is %v, ErrorMsg is %v", rsp.Code, err)
	}
	return rsp, err
}

func (b *BotService) getToken() error {
	rsp, err := getTenantAccessToken(*(b.AppId), *(b.AppSecret))
	if err == nil {
		expire := time.Duration(rsp.Expire - 5*60) //  离过期时间小于 5分钟
		now := time.Now()
		b.TenantAccessToken = new(TenantAccessToken)
		b.TenantAccessToken.ExpireTime = now.Add(expire * time.Second)
		b.TenantAccessToken.Token = rsp.TenantAccessToken
	}
	return err
}

func (b *BotService) getValidTenantAccessToken() error {
	if b.AppId == nil || b.AppSecret == nil {
		return fmt.Errorf("【Lark_Bot】empty appId or appSecret, init appId and appSecret first")
	}
	if (b.TenantAccessToken == nil) || b.TenantAccessToken.ExpireTime.Before(time.Now()) { // 没有 accessToken
		err := b.getToken()
		return err
	}
	return nil
}

func (b *BotService) makeHeader(data ...map[string]string) map[string]string {
	header := make(map[string]string)
	_ = b.getValidTenantAccessToken()
	header["Authorization"] = fmt.Sprintf("Bearer %s", b.TenantAccessToken.Token)
	header["Content-Type"] = "application/json"
	if len(data) > 0 {
		for _, d := range data {
			for k, v := range d {
				header[k] = v
			}
		}
	}
	return header
}

// AddUser2Group 将用户加入群组
func (b *BotService) AddUser2Group(email string, chatID string) error {
	userID := b.GetUserIDByEmail(email)
	header := b.makeHeader()
	params := map[string]interface{}{"chat_id": chatID, "user_ids": []string{userID}}
	var ret BasicResponse
	err := utils.HttpPostWithHeader(CreateChatGroupURL, header, params, &ret)
	if (err != nil) || (ret.Code != 0) {
		logs.Errorf("【Lark_Bot】Error in AddUser2Group:%v", err)
		return fmt.Errorf("%s", ret.Msg)
	}
	return nil
}

func (b *BotService) GetLarkGroupChatID() ([]*LarkGroupInfo, error) {
	header := b.makeHeader()
	var ret ChatGroup
	err := utils.HttpPostWithHeader(GroupChatListURL, header, nil, &ret)
	if (err != nil) || (ret.Code != 0) {
		logs.Errorf("【Lark_Bot】Error in GetLarkGroupChatID:%v", err)
		return nil, err
	}
	return ret.Data.Groups, nil
}

func (b *BotService) SendMessageCard(receiver ReceiverInfo, content interface{}) error {
	err := b.getValidTenantAccessToken()
	if err != nil {
		return fmt.Errorf("【Lark_Bot】Send message failed——Cannot send message, error= %v", err)
	}

	var msg SendMessageForm
	msgType := "interactive"
	msg.OpenId = receiver.OpenId
	msg.UserId = receiver.UserId
	msg.ChatId = receiver.ChatId
	msg.Email = receiver.Email
	msg.MsgType = &msgType
	msg.Content = content

	var rsp LarkSendResponse
	err = utils.HttpPostWithHeader(SendMessageURL, b.makeHeader(), msg, &rsp)
	if err != nil {
		return err
	}
	if rsp.Code != 0 {
		return fmt.Errorf("【Lark_Bot】send message card error:%v", rsp)
	}
	return err
}

func (b *BotService) SetConsumerFunction(f func(interface{})) {
	b.Consumer = f
}

func (b *BotService) ProduceMsg(data interface{}) {
	b.Messages <- data
}

func (b *BotService) ConsumeMsg() {
	for {
		select {
		case ms := <-b.Messages:
			b.Consumer(ms)
		}
	}
}

func InitOnCallLarkBot() {
	switch utils.TellEnv() {
	case "loc", "boe":
		OnCallBot.Init(BOEBotAPIID, BOEBotAPISecret)
	case "tce":
		OnCallBot.Init(OnCallBotAPIID, OnCallBotAPISecret)
	}
}

func InitByteSecLarkBot() {
	switch utils.TellEnv() {
	case "loc", "boe":
		ByteSecBot.Init(BOEBotAPIID, BOEBotAPISecret)
	case "tce":
		ByteSecBot.Init(ByteSecBotAPIID, ByteSecBotAPISecret)
	}
}
