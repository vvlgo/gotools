package example_test

import (
	"encoding/json"
	"fmt"
	"github.com/vvlgo/gotools/redisclient"
	"github.com/vvlgo/gotools/tencenttools/wechatapi"
	"testing"
)

const corpid = "xxxxxxxxx"
const corpsecret = "xxxxxxxxxxxxx"
const agentId = "100000000"

func TestWechtTools(t *testing.T) {
	redisclient.RedisInit("127.0.0.1", "6379", "")
	reConn, err := redisclient.GetRedisConn(0)
	if err != nil {
		panic(err)
	}
	accessToken, _ := wechatapi.GetAccessToken(corpid, corpsecret, *reConn)
	fmt.Println(accessToken)

	busiTicket, _ := wechatapi.GetBusiTicket(corpid, corpsecret, *reConn)
	fmt.Println(busiTicket)

	appTicket, _ := wechatapi.GetAppTicket(corpid, corpsecret, *reConn)
	fmt.Println(appTicket)

	user, err := wechatapi.GetUserByUserID(corpid, corpsecret, "WeiTao", *reConn)
	bytes, err := json.Marshal(user)
	mp := make(map[string]interface{})
	json.Unmarshal(bytes, &mp)
	fmt.Println(mp)
}
