package dingdingclient

import (
	"fmt"
	"testing"
)

func TestSendMessage(t *testing.T) {
	bootUrl := "https://oapi.dingtalk.com/robot/send?access_token=ca6441d14175831d2f1e4e5409421b7a1a7859824c2cbb59cd6a705f1f318928"
	secret := "SEC8864fda8960e963bbac681d27aab862957666aa97b5f0e449e9082a6fd6b26ea"
	robot := NewWebhook(bootUrl, secret)
	err := robot.Send(nil, &TextBody{Text: "bingo"})
	if err != nil {
		fmt.Println(err)
	}
}
