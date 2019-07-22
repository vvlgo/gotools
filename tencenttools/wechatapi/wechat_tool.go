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
	Errcode      int          `json:"errcode"`
	Errmsg       string       `json:"errmsg"`
	AccessToken  string       `json:"access_token"`
	Ticket       string       `json:"ticket"`
	ExpiresIn    int          `json:"expires_in"`
	Invaliduser  string       `json:"invaliduser"`
	Data         interface{}  `json:"data"`
	UserId       string       `json:"USERID"`
	OpenId       string       `json:"OPENID"`
	DepartmentID int          `json:"id"`
	Department   []Department `json:"department"`
	Userlist     []User       `json:"userlist"`
}

//Department 部门表及字段
type Department struct {
	//企业微信字段
	WechatID       int    `gorm:"column:wechat_id;COMMENT:'企业微信创建的部门id';"                                               form:"wechat_id"       json:"id"`
	DepartmentName string `gorm:"column:department_name;type:varchar(50);COMMENT:'部门名字';"                                    form:"department_name" json:"name"`
	ParentID       int    `gorm:"column:parent_id;COMMENT:'父亲部门id。根部门为1';"                                              form:"parent_id"       json:"parentid"`
	OrderNum       int64  `gorm:"column:order_num;type:bigint;COMMENT:'在父部门中的次序值。order值大的排序靠前。值范围是[0, 2^32)';" form:"order_num"           json:"order"`
}

//User 用户字段
type User struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	//企业微信字段
	UserID          string `gorm:"column:userid;COMMENT:'成员UserID。对应管理端的帐号';"                                              form:"userid"            json:"userid"`
	UserName        string `gorm:"column:user_name;COMMENT:'成员名称';"                                                               form:"user_name"         json:"name"`
	Mobile          string `gorm:"column:mobile;COMMENT:'手机号码';"                                                                  form:"mobile"            json:"mobile"`
	WechatIDs       []int  `gorm:"-"                                                                                                                           json:"department"`
	DeptWechatIDs   string `gorm:"column:dept_wechat_ids;COMMENT:'成员所属部门id列表，仅返回该应用有查看权限的部门id';"               form:"dept_wechat_ids"   json:"dept_wechat_ids"`
	Position        string `gorm:"column:position;COMMENT:'职务信息；';"                                                              form:"position"          json:"position"`
	Gender          string `gorm:"column:gender;COMMENT:'性别。0表示未定义，1表示男性，2表示女性';"                                   form:"gender"            json:"gender"`
	Email           string `gorm:"column:email;COMMENT:'邮箱';"                                                                       form:"email"             json:"email"`
	Avatar          string `gorm:"column:avatar;COMMENT:'头像url。注：如果要获取小图将url最后的”/0”改成”/100”即可';"              form:"avatar"            json:"avatar"`
	Status          int    `gorm:"column:status;COMMENT:'激活状态: 1=已激活，2=已禁用，4=未激活';"                                    form:"status"            json:"status"`
	Enable          int    `gorm:"column:enable;COMMENT:'成员启用状态。1表示启用的成员，0表示被禁用。服务商调用接口不会返回此字段';"  form:"enable"            json:"enable"`
	Isleader        int    `gorm:"column:isleader;COMMENT:'无';"                                                                      form:"isleader"          json:"isleader"`
	HideMobile      int    `gorm:"column:hide_mobile;COMMENT:'无';"                                                                   form:"hide_mobile"       json:"hide_mobile"`
	Telephone       string `gorm:"column:telephone;COMMENT:'座机';"                                                                   form:"telephone"         json:"telephone"`
	Order           []int  `gorm:"-"                                                                                                                           json:"order"`
	Orders          string `gorm:"column:orders;COMMENT:'部门内的排序值，32位整数，默认为0';"                                       form:"orders"         json:"orders"`
	QrCode          string `gorm:"column:qr_code;COMMENT:'员工个人二维码，扫描可添加为外部联系人';"                                   form:"qr_code"           json:"qr_code"`
	Alias           string `gorm:"column:alias;COMMENT:'别名';"                                                                       form:"alias"             json:"alias"`
	IsLeaderInDept  []int  `gorm:"-"                                                                                                                           json:"is_leader_in_dept"`
	IsLeaderInDepts string `gorm:"column:is_leader_in_depts;COMMENT:'表示在所在的部门内是否为上级';"                                   form:"is_leader_in_depts" json:"is_leader_in_depts"`
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
	AuthTplId     string `json:"auth_tpl_id"`
}

//企业微信人员请求对应字段
type RequestUser struct {
	UserID         string `json:"userid,omitempty"`
	UserName       string `json:"name,omitempty"`
	Alias          string `json:"alias,omitempty"`
	Mobile         string `json:"mobile,omitempty"`
	Department     []int  `json:"department,omitempty"`
	Position       string `json:"position,omitempty"`
	Gender         string `json:"gender,omitempty"`
	Email          string `json:"email,omitempty"`
	IsLeaderInDept []int  `json:"is_leader_in_dept,omitempty"`
	Enable         int    `json:"enable,omitempty"`
}

//企业微信部门请求对应字段
type RequestDepartment struct {
	Department   string `json:"name,omitempty"`
	ParentID     int    `json:"parentid,omitempty"`
	DepartmentID int    `json:"id,omitempty"`
}

/*
GetAccessToken 企业微信获取AccessToken
redisConn AccessToken缓存库
*/
func GetAccessToken(corpid, corpsecret string, redisConn redisclient.MyRedisReConn, key string) (string, error) {

	code, _ := redis.String(redisConn.Redo("Get", key+"access_token"))
	if code != "" {
		return code, nil
	} else {
		url1 := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpid + "&corpsecret=" + corpsecret
		resp, err := httptool.GET(url1, nil)
		if err != nil {
			return "", err
		}
		re := WechatRep{}
		err = json.Unmarshal([]byte(resp), &re)
		if err != nil {
			return "", err
		}
		_, err = redisConn.Redo("Set", key+"access_token", re.AccessToken, "EX", re.ExpiresIn)
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
func GetAppTicket(corpid, corpsecret string, redisConn redisclient.MyRedisReConn, key string) (string, error) {

	ticket, _ := redis.String(redisConn.Redo("Get", "appticket"))
	if ticket != "" {
		return ticket, nil
	} else {
		url3 := ""
		accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
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
func GetBusiTicket(corpid, corpsecret string, redisConn redisclient.MyRedisReConn, key string) (string, error) {
	ticket, _ := redis.String(redisConn.Redo("Get", "busiticket"))
	if ticket != "" {
		return ticket, nil
	} else {
		url3 := ""
		accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
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

func GetSignature(url, corpid, corpsecret, agentId string, redisConn redisclient.MyRedisReConn, key string) (*Signature, error) {
	appticket, err := GetAppTicket(corpid, corpsecret, redisConn, key)
	if err != nil {
		return nil, err
	}
	busiticket, err := GetBusiTicket(corpid, corpsecret, redisConn, key)
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
func SendMsg(orderUrl, toUser, corpid, corpsecret, title, cardInfo string, agentid int, redisConn redisclient.MyRedisReConn, key string) (bool, error) {
	url2 := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
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
	if re.Errcode == 0 && (re.Errmsg == "OK" || re.Errmsg == "ok") {
		return true, nil
	}
	return false, err
}

/*
GetUserByCode 企业微信获取用户基础数据
redisConn AccessToken缓存库
*/
func GetUserByCode(corpid, corpsecret, code string, redisConn redisclient.MyRedisReConn, key string) (*WechatRep, error) {
	codeUrl := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=ACCESS_TOKEN&code=CODE"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return nil, err
	}
	codeUrl = strings.ReplaceAll(codeUrl, "ACCESS_TOKEN", accessToken)
	codeUrl = strings.ReplaceAll(codeUrl, "CODE", code)
	res, err := httptool.GET(codeUrl, nil)
	if err != nil {
		return nil, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return nil, err
	}
	if re.Errcode == 0 && (re.Errmsg == "OK" || re.Errmsg == "ok") {
		return &re, nil
	}
	return nil, err
}

/*
GetUserByUserID 企业微信获取用户详细信息数据
redisConn AccessToken缓存库
*/
func GetUserByUserID(corpid, corpsecret, userid string, redisConn redisclient.MyRedisReConn, key string) (*User, error) {
	userUrl := "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=ACCESS_TOKEN&userid=USERID"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return nil, err
	}
	userUrl = strings.ReplaceAll(userUrl, "ACCESS_TOKEN", accessToken)
	userUrl = strings.ReplaceAll(userUrl, "USERID", userid)
	res, err := httptool.GET(userUrl, nil)
	if err != nil {
		return nil, err
	}
	re := User{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return nil, err
	}
	if re.Errcode == 0 && (re.Errmsg == "OK" || re.Errmsg == "ok") {
		return &re, nil
	}
	return nil, err
}

/*
GetDepartmentList 企业微信获取所有部门数据
redisConn AccessToken缓存库
*/
func GetDepartmentList(corpid, corpsecret, userid string, redisConn redisclient.MyRedisReConn, key string) (*WechatRep, error) {
	Url := "https://qyapi.weixin.qq.com/cgi-bin/department/list"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return nil, err
	}
	params := make(map[string]string)
	params["access_token"] = accessToken
	res, err := httptool.GET(Url, params)
	if err != nil {
		return nil, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return nil, err
	}
	if re.Errcode == 0 && (re.Errmsg == "OK" || re.Errmsg == "ok") {
		return &re, nil
	}
	return nil, err
}

/*
GetDepartmentUserList 企业微信获取部门人员信息数据
redisConn AccessToken缓存库
*/
func GetDepartmentUserList(corpid, corpsecret, departmenID string, redisConn redisclient.MyRedisReConn, key string) (*WechatRep, error) {
	Url := "https://qyapi.weixin.qq.com/cgi-bin/user/list"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return nil, err
	}
	params := make(map[string]string)
	params["access_token"] = accessToken
	params["department_id"] = departmenID
	params["fetch_child"] = "1"
	res, err := httptool.GET(Url, params)
	if err != nil {
		return nil, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return nil, err
	}
	if re.Errcode == 0 && (re.Errmsg == "OK" || re.Errmsg == "ok") {
		return &re, nil
	}
	return nil, err
}

/*
AddUser 企业微信新增成员
redisConn AccessToken缓存库
*/
func AddUser(corpid, corpsecret string, user RequestUser, redisConn redisclient.MyRedisReConn, key string) (bool, error) {
	UrL := "https://qyapi.weixin.qq.com/cgi-bin/user/create?access_token=ACCESS_TOKEN"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return false, err
	}
	UrL = strings.ReplaceAll(UrL, "ACCESS_TOKEN", accessToken)
	res, err := httptool.POST(UrL, nil, user)
	if err != nil {
		return false, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return false, err
	}
	if re.Errcode == 0 && (re.Errmsg == "CREATED" || re.Errmsg == "created") {
		return true, nil
	}
	return false, err
}

/*
UpdateUser 企业微信更新成员
redisConn AccessToken缓存库
*/
func UpdateUser(corpid, corpsecret string, user RequestUser, redisConn redisclient.MyRedisReConn, key string) (bool, error) {
	UrL := "https://qyapi.weixin.qq.com/cgi-bin/user/update?access_token=ACCESS_TOKEN"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return false, err
	}
	UrL = strings.ReplaceAll(UrL, "ACCESS_TOKEN", accessToken)
	res, err := httptool.POST(UrL, nil, user)
	if err != nil {
		return false, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return false, err
	}
	if re.Errcode == 0 && (re.Errmsg == "UPDATED" || re.Errmsg == "updated") {
		return true, nil
	}
	return false, err
}

/*
DelUser 企业微信删除成员
redisConn AccessToken缓存库
*/
func DelUser(corpid, corpsecret, userID string, redisConn redisclient.MyRedisReConn, key string) (bool, error) {
	UrL := "https://qyapi.weixin.qq.com/cgi-bin/user/delete?access_token=ACCESS_TOKEN&userid=USERID"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return false, err
	}
	UrL = strings.ReplaceAll(UrL, "ACCESS_TOKEN", accessToken)
	UrL = strings.ReplaceAll(UrL, "USERID", userID)
	res, err := httptool.GET(UrL, nil)
	if err != nil {
		return false, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return false, err
	}
	if re.Errcode == 0 && (re.Errmsg == "DELETED" || re.Errmsg == "deleted") {
		return true, nil
	}
	return false, err
}

/*
DelUserList 企业微信批量删除成员
redisConn AccessToken缓存库
*/
func DelUserList(corpid, corpsecret string, userList []string, redisConn redisclient.MyRedisReConn, key string) (bool, error) {
	UrL := "https://qyapi.weixin.qq.com/cgi-bin/user/batchdelete?access_token=ACCESS_TOKEN"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return false, err
	}
	UrL = strings.ReplaceAll(UrL, "ACCESS_TOKEN", accessToken)
	data := make(map[string][]string)
	data["useridlist"] = userList
	res, err := httptool.POST(UrL, nil, data)
	if err != nil {
		return false, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return false, err
	}
	if re.Errcode == 0 && (re.Errmsg == "DELETED" || re.Errmsg == "deleted") {
		return true, nil
	}
	return false, err
}

/*
AddDepartment 企业微信新增部门
redisConn AccessToken缓存库
*/
func AddDepartment(corpid, corpsecret string, dep RequestDepartment, redisConn redisclient.MyRedisReConn, key string) (id int, err error) {
	UrL := "https://qyapi.weixin.qq.com/cgi-bin/department/create?access_token=ACCESS_TOKEN"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return -1, err
	}
	UrL = strings.ReplaceAll(UrL, "ACCESS_TOKEN", accessToken)
	res, err := httptool.POST(UrL, nil, dep)
	if err != nil {
		return -1, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return -1, err
	}
	if re.Errcode == 0 && (re.Errmsg == "CREATED" || re.Errmsg == "created") {
		return re.DepartmentID, nil
	}
	return -1, err
}

/*
UpdateDepartment 企业微信更新部门
redisConn AccessToken缓存库
*/
func UpdateDepartment(corpid, corpsecret string, dep Department, redisConn redisclient.MyRedisReConn, key string) (bool, error) {
	UrL := "https://qyapi.weixin.qq.com/cgi-bin/department/update?access_token=ACCESS_TOKEN"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return false, err
	}
	UrL = strings.ReplaceAll(UrL, "ACCESS_TOKEN", accessToken)
	res, err := httptool.POST(UrL, nil, dep)
	if err != nil {
		return false, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return false, err
	}
	if re.Errcode == 0 && (re.Errmsg == "UPDATED" || re.Errmsg == "updated") {
		return true, nil
	}
	return false, err
}

/*
DelDepartment 企业微信删除部门
redisConn AccessToken缓存库
*/
func DelDepartment(corpid, corpsecret, departmenID string, redisConn redisclient.MyRedisReConn, key string) (bool, error) {
	UrL := "https://qyapi.weixin.qq.com/cgi-bin/department/delete?access_token=ACCESS_TOKEN&id=ID"
	accessToken, err := GetAccessToken(corpid, corpsecret, redisConn, key)
	if accessToken == "" {
		return false, err
	}
	UrL = strings.ReplaceAll(UrL, "ACCESS_TOKEN", accessToken)
	UrL = strings.ReplaceAll(UrL, "ID", departmenID)
	res, err := httptool.GET(UrL, nil)
	if err != nil {
		return false, err
	}
	re := WechatRep{}
	err = json.Unmarshal(res, &re)
	if err != nil {
		return false, err
	}
	if re.Errcode == 0 && (re.Errmsg == "DELETED" || re.Errmsg == "deleted") {
		return true, nil
	}
	return false, err
}
