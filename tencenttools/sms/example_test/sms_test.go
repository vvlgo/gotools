package example_test

import (
	"fmt"
	"github.com/vvlgo/gotools/redisclient"
	"github.com/vvlgo/gotools/tencenttools/sms"
	"testing"
)

//腾讯短信平台
const (
	sdkappid = `xxxxxx`
	appkey   = `xxxxxxxxxxxxxxxxxxxxxxxxx`
	tpl_id   = 100000
)

func TestSms(t *testing.T) {
	redisclient.RedisInit("127.0.0.1", "6379", "")
	reConn0, err := redisclient.GetRedisConn(0)
	reConn1, err := redisclient.GetRedisConn(1)
	if err != nil {
		panic(err)
	}
	var phone = "13628005220"
	b, err := sms.SendSms(phone, sdkappid, appkey, "", tpl_id, *reConn0, 60, *reConn1, 60*60)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}
