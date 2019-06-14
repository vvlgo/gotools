package sms

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/vvlgo/gotools/httptool"
	"github.com/vvlgo/gotools/redisclient"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type SmsRequest struct {
	Ext    string   `json:"ext"`
	Extent string   `json:"extent"`
	Params []string `json:"params"`
	Sig    string   `json:"sig"`
	Sign   string   `json:"sign"`
	Tel    Tel      `json:"tel"`
	Time   int64    `json:"time"`
	TplId  int      `json:"tpl_id"`
}

type Tel struct {
	Mobile     string `json:"mobile"`
	Nationcode string `json:"nationcode"`
}

type SmsResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Fee    int    `json:"fee"`
	Sid    string `json:"sid"`
}
type SmscCallback struct {
	UserReceiveTime string `json:"user_receive_time"`
	Nationcode      string `json:"nationcode"`
	Mobile          string `json:"mobile"`
	ReportStatus    string `json:"report_status"`
	Errmsg          string `json:"errmsg"`
	Description     string `json:"description"`
	Sid             string `json:"sid"`
}

/*
SendSms 发送短信验证码
phone 电话
sdkappid 腾讯短信唯一id,控制台查看
tplId 腾讯短信模板id,控制台查看
redisConnSms redis连接对象,缓存短信验证码的库
redisConnPhoneCount redis连接对象,统计电话请求次数的库
*/
func SendSms(phone, sdkappid, appkey, sigin string, tplId int, redisConnSms redisclient.MyRedisReConn, smsCodeExpiresTime int, redisConnPhoneCount redisclient.MyRedisReConn, phoneCountExpiresTime int) (bool, error) {
	code := CreateSmsCode()
	unix := time.Now().Unix()
	var url1 = `https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=` + sdkappid + `&random=` + code
	smsRe := SmsRequest{}
	smsRe.Params = append(smsRe.Params, code)
	smsRe.Params = append(smsRe.Params, "2")
	smsRe.Tel.Mobile = phone
	smsRe.Tel.Nationcode = "86"
	smsRe.Time = unix
	smsRe.TplId = tplId
	smsRe.Sign = sigin
	smsRe.Sig = sig(smsRe.Tel.Mobile, appkey, code, strconv.FormatInt(unix, 10))
	body, err := httptool.POST(url1, nil, smsRe)
	re := SmsResult{}
	err = json.Unmarshal([]byte(body), &re)
	if err != nil {
		return false, err
	}
	if re.Result == 0 && re.Errmsg == "OK" {
		_, err = redisConnSms.Redo("Set", phone, code, "EX", smsCodeExpiresTime)
		if err != nil {
			return false, err
		}
		err = RecordPhoneUseNums(phone, redisConnPhoneCount, phoneCountExpiresTime)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, errors.New(re.Errmsg)

}

/*
CreateSmsCode 产生随机六位码
*/
func CreateSmsCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result = ""
	for i := 0; i < 6; i++ {
		result += strconv.Itoa(r.Intn(10))
	}
	return result
}

/*
CheckPhoneIsNormal 检查手机是正常请求手机短信码，超过三次请求但又不使用验证码则锁定此号码不得请求，锁定时间24小时
redisConn redis连接对象，统计手机请求次数的库
*/
func CheckPhoneIsNormal(phone string, redisConn redisclient.MyRedisReConn) (bool, error) {
	nums, err := redis.Int(redisConn.Redo("Get", phone))
	x := fmt.Sprintf("%s", err)
	if strings.Index(x, "nil returned") > 0 {
		return true, nil
	}
	if nums >= 3 {
		return false, nil
	}
	return true, nil

}

/*
RecordPhoneUseNums 手机请求验证码，请求一次记录一次
redisConn redis连接对象，统计电话号码的库
*/
func RecordPhoneUseNums(phone string, redisConnPhoneCount redisclient.MyRedisReConn, phoneCountExpiresTime int) error {
	nums, err := redis.Int(redisConnPhoneCount.Redo("Get", phone))
	x := fmt.Sprintf("%s", err)
	if strings.Index(x, "nil returned") > 0 {
		_, err = redisConnPhoneCount.Redo("Set", phone, 1, "EX", phoneCountExpiresTime)
	}
	_, err = redisConnPhoneCount.Redo("Set", phone, nums+1, "EX", phoneCountExpiresTime)
	if err != nil {
		return err
	}
	return nil
}

/*
DelPhoneUseNums 如果手机短信码被使用了，从缓存中删除对此手机号的记录
redisConn redis连接对象，统计电话号码的库
*/
func DelPhoneUseNums(phone string, redisConn redisclient.MyRedisReConn) error {
	_, err := redisConn.Redo("DEL", phone)
	if err != nil {
		return err
	}
	return nil
}

/*
短信签名
*/
func sig(strMobile, strAppKey, strRand, strTime string) string {
	sig := `appkey=` + strAppKey + `&random=` + strRand + `&time=` + strTime + `&mobile=` + strMobile
	h := sha256.New()
	h.Write([]byte(sig))
	return fmt.Sprintf("%x", h.Sum(nil))
}
