package wechatapi

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/vvlgo/gotools/httptool"
	"github.com/vvlgo/gotools/redisclient"
	"github.com/vvlgo/gotools/tencenttools/sms"
	"strconv"
	"strings"
	"time"
)

type WechatRep struct {
	Errcode     int         `json:"errcode"`
	Errmsg      string      `json:"errmsg"`
	AccessToken string      `json:"access_token"`
	Ticket      string      `json:"ticket"`
	ExpiresIn   int         `json:"expires_in"`
	Invaliduser string      `json:"invaliduser"`
	Data        interface{} `json:"data"`
	UserId      string      `json:"USERID"`
	OpenId      string      `json:"OPENID"`
}

type WechatMsg struct {
	Touser   string   `json:"touser"`
	Msgtype  string   `json:"msgtype"`
	Toall    int      `json:"toall"`
	Agentid  int      `json:"agentid"`
	Textcard Textcard `json:"textcard"`
}
type Textcard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}

type Signature struct {
	AppId         string `json:"appId"`
	AgentId       string `json:"agentId"`
	AppSignature  string `json:"app_signature"`
	BusiSignature string `json:"busi_signature"`
	Noncestr      string `json:"nonceStr"`
	Timestamp     int64  `json:"timestamp"`
}

/*
GetAccessToken 企业微信获取AccessToken
redisConn AccessToken缓存库
*/
func GetAccessToken(corpid string, redisConn redisclient.MyRedisReConn) (string, error) {

	code, _ := redis.String(redisConn.Redo("Get", "access_token"))
	if code != "" {
		return code, nil
	} else {
		url1 := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=aehiE0og9IpDOBG-eKYgSGJA_lClwLgZJ8BQHbPUBOQ"
		resp, err := httptool.GET(url1, nil)
		if err != nil {
			return "", err
		}
		re := WechatRep{}
		err = json.Unmarshal([]byte(resp), &re)
		if err != nil {
			return "", err
		}
		_, err = redisConn.Redo("Set", "access_token", re.AccessToken, "EX", re.ExpiresIn)
		if err != nil {
			return "", err
		}
		return re.AccessToken, nil
	}

}

/*
GetAppTicket 企业微信获取AppTicket
redisConn AppTicket缓存库，和AccessToken同库
*/
func GetAppTicket(corpid string, redisConn redisclient.MyRedisReConn) (string, error) {

	ticket, _ := redis.String(redisConn.Redo("Get", "appticket"))
	if ticket != "" {
		return ticket, nil
	} else {
		url3 := ""
		accessToken, err := GetAccessToken(corpid, redisConn)
		if err != nil {
			return "", err
		}
		if accessToken != "" {
			url3 = url3 + "https://qyapi.weixin.qq.com/cgi-bin/ticket/get?access_token=" + accessToken + "&type=agent_config"
			resp, err := httptool.GET(url3, nil)
			if err != nil {
				return "", err
			}
			re := WechatRep{}
			err = json.Unmarshal([]byte(resp), &re)
			if err != nil {
				return "", err
			}
			_, err = redisConn.Redo("Set", "appticket", re.Ticket, "EX", re.ExpiresIn)
			if err != nil {
				return "", err
			}
			return re.Ticket, nil
		}
		return "", err
	}
}

/*
GetBusiTicket 企业微信获取BusiTicket
redisConn BusiTicket缓存库，和AccessToken同库
*/
func GetBusiTicket(corpid string, redisConn redisclient.MyRedisReConn) (string, error) {
	ticket, _ := redis.String(redisConn.Redo("Get", "busiticket"))
	if ticket != "" {
		return ticket, nil
	} else {
		url3 := ""
		accessToken, err := GetAccessToken(corpid, redisConn)
		if accessToken != "" {
			url3 = url3 + "https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=" + accessToken
			resp, err := httptool.GET(url3, nil)
			if err != nil {
				return "", err
			}
			re := WechatRep{}
			err = json.Unmarshal([]byte(resp), &re)
			if err != nil {
				return "", err
			}
			_, err = redisConn.Redo("Set", "busiticket", re.Ticket, "EX", re.ExpiresIn)
			if err != nil {
				return "", err
			}
			return re.Ticket, nil
		}
		return "", err
	}
}

func GetSignature(url, corpid, agentId string, redisConn redisclient.MyRedisReConn) (*Signature, error) {
	appticket, err := GetAppTicket(corpid, redisConn)
	if err != nil {
		return nil, err
	}
	busiticket, err := GetBusiTicket(corpid, redisConn)
	if err != nil {
		return nil, err
	}
	noncestr := sms.CreateSmsCode()
	unix := time.Now().Unix()
	timestamp := strconv.FormatInt(unix, 10)
	if appticket != "" && busiticket != "" {
		fmt.Println(appticket, "=====", busiticket)
		s := "jsapi_ticket=JSAPITICKET&noncestr=NONCESTR&timestamp=TIMESTAMP&url=URL"
		s = strings.ReplaceAll(s, "JSAPITICKET", appticket)
		s = strings.ReplaceAll(s, "NONCESTR", noncestr)
		s = strings.ReplaceAll(s, "TIMESTAMP", timestamp)
		s = strings.ReplaceAll(s, "URL", url)
		i := Sha1([]byte(s))
		s1 := "jsapi_ticket=JSAPITICKET&noncestr=NONCESTR&timestamp=TIMESTAMP&url=URL"
		s1 = strings.ReplaceAll(s1, "JSAPITICKET", busiticket)
		s1 = strings.ReplaceAll(s1, "NONCESTR", noncestr)
		s1 = strings.ReplaceAll(s1, "TIMESTAMP", timestamp)
		s1 = strings.ReplaceAll(s1, "URL", url)
		i2 := Sha1([]byte(s1))
		signature := Signature{}
		signature.AppId = corpid
		signature.AgentId = agentId
		signature.AppSignature = i
		signature.BusiSignature = i2
		signature.Noncestr = noncestr
		signature.Timestamp = unix
		return &signature, nil
	}
	return nil, err
}
func Sha1(data []byte) string {
	sha1 := sha1.New()
	sha1.Write(data)
	return hex.EncodeToString(sha1.Sum([]byte(nil)))
}

/*
SendMsg 发送审核信息到企业微信,，卡片信息
模板根据自己情景修改
*/
func SendMsg(orderUrl, toUser, corpid, title, cardInfo string, agentid int, redisConn redisclient.MyRedisReConn) (bool, error) {
	url2 := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
	accessToken, err := GetAccessToken(corpid, redisConn)
	if accessToken == "" {
		return false, err
	} else {
		url2 = url2 + accessToken
	}
	wechatMsg := WechatMsg{}
	wechatMsg.Touser = toUser
	wechatMsg.Msgtype = "textcard"
	wechatMsg.Agentid = agentid
	wechatMsg.Toall = 0
	wechatMsg.Textcard.Title = title
	wechatMsg.Textcard.Description = cardInfo
	wechatMsg.Textcard.Url = orderUrl
	wechatMsg.Textcard.Btntxt = "更多"
	body, err := httptool.POST(url2, nil, wechatMsg)
	if err != nil {
		return false, err
	}
	re := WechatRep{}
	err = json.Unmarshal([]byte(body), &re)
	if err != nil {
		return false, err
	}
	if re.Errcode == 0 && re.Errmsg == "ok" {
		return true, nil
	}
	return false, err
}
