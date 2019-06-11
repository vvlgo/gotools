package example_test

import (
	"fmt"
	"github.com/vvlgo/gotools/redisclient"
	"github.com/vvlgo/gotools/tencenttools/wechatapi"
	"testing"
)

const corpid = "xxxxxxxxxxx"
const agentId = "100000000"

func TestWechtTools(t *testing.T) {
	redisclient.RedisInit("127.0.0.1", "6379", "")
	reConn, err := redisclient.GetRedisConn(0)
	if err != nil {
		panic(err)
	}
	accessToken := wechatapi.GetAccessToken(corpid, *reConn)
	fmt.Println(accessToken)

	busiTicket := wechatapi.GetBusiTicket(corpid, *reConn)
	fmt.Println(busiTicket)

	appTicket := wechatapi.GetAppTicket(corpid, *reConn)
	fmt.Println(appTicket)
}
