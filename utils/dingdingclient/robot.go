// 钉钉机器人

package dingdingclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 自定义机器人 https://open.dingtalk.com/document/group/custom-robot-access
type webhook struct {
	url    string // Webhook地址 https://oapi.dingtalk.com/robot/send?access_token=ca6441d14175831d2f1e4e5409421b7a1a7859824c2cbb59cd6a705f1f318928
	secret string // 密钥 SEC8864fda8960e963bbac681d27aab862957666aa97b5f0e449e9082a6fd6b26ea
}

func NewWebhook(url, secret string) *webhook {
	return &webhook{
		url:    url,
		secret: secret,
	}
}

// Webhook签名算法：把timestamp+"\n"+密钥当做签名字符串，使用HmacSHA256算法计算签名，然后进行Base64 encode，最后再把签名参数再进行urlEncode，得到最终的签名（需要使用UTF-8字符集）。
func (webhook *webhook) sign() string {
	timestamp := time.Now().UnixMilli()
	signStr := fmt.Sprintf("%d\n%s", timestamp, webhook.secret)
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(webhook.secret))
	// Write Data to it
	h.Write([]byte(signStr))
	// base64
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// urlEncode
	sign = url.QueryEscape(sign)
	// 拼接最终请求地址
	return fmt.Sprintf(webhook.url+"&timestamp=%d&sign=%s", timestamp, sign)
}

func (webhook *webhook) send(param map[string]interface{}) error {
	bootUrl := webhook.sign()
	b, err := json.Marshal(param)
	if err != nil {
		return err
	}
	reader := strings.NewReader(string(b))
	resp, err := http.Post(bootUrl, "application/json", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var jsonMap map[string]interface{}
	json.Unmarshal(bytes, &jsonMap)
	if errMsg, ok := jsonMap["errmsg"]; ok {
		if errMsg != "ok" {
			fmt.Printf("%#v\n", jsonMap)
			return errors.New(errMsg.(string))
		}
	} else {
		return errors.New(string(bytes))
	}
	return nil
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	AtUserIds []string `json:"atUserIds"`
	IsAtAll   bool     `json:"isAtAll"` // 是否@所有人。
}

// TextBody text类型
type TextBody struct {
	Text string `json:"content"`
}

// LinkBody link类型
type LinkBody struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

// MarkdownBody markdown类型
type MarkdownBody struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// ActionCardWholeBody 整体跳转ActionCard类型
type ActionCardWholeBody struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle    string `json:"singleTitle"`
	SingleUrl      string `json:"singleURL"`
}

// SingleBody 独立跳转ActionCard类型
type SingleBody struct {
	SingleTitle string `json:"title"`
	SingleUrl   string `json:"actionURL"`
}
type ActionCardSingleBody struct {
	Title          string       `json:"title"`
	Text           string       `json:"text"`
	BtnOrientation string       `json:"btnOrientation"`
	Singles        []SingleBody `json:"btns"`
}

// FeedCardBody FeedCard类型
type FeedCardBody struct {
	Title      string `json:"title"`
	PicUrl     string `json:"picURL"`
	MessageUrl string `json:"messageURL"`
}

// Send 发送自定义机器人消息
func (webhook *webhook) Send(at *At, body interface{}) error {
	param := make(map[string]interface{})
	if at != nil {
		param["at"] = at
	}
	messageType := "text"
	switch body.(type) {
	case TextBody:
		messageType = "text"
	case LinkBody:
		messageType = "link"
	case MarkdownBody:
		messageType = "markdown"
	case ActionCardWholeBody:
		messageType = "actionCard"
	case ActionCardSingleBody:
		messageType = "actionCard"
	case FeedCardBody:
		messageType = "feedCard"
	}
	param["msgtype"] = messageType
	param[messageType] = body
	return webhook.send(param)
}
