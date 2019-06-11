package ecode

/*
系统常量
Web* web页面返回code
Sql* 数据库常量状态
Sys* 系统常量
*/
const (
	WebOk                    = 0  // 正确
	WebUsernameOrPasswordErr = -1 // 用户名或密码错误
	WebCaptchaErr            = -2 // 验证码错误
	WebDataexist             = -3 // 数据存在

	WebSessionExpires   = -4 //session过期
	WebSmsNotExpires    = -5 // 短信码未过期不能请求
	WebSmsExpires       = -6 // 短信码失效
	WebSmsErr           = -7 // 短信码错误
	WebSmsPhoneExcption = -8 // 超过次数被锁定

	WebServiceUpdate = -100 // 系统升级中
	WebSysExcption   = -101 // 系统异常，具体信息具体说明

	WebNotModified        = -304 // 木有改动
	WebTemporaryRedirect  = -307 // 撞车跳转
	WebRequestErr         = -400 // 请求错误
	WebUnauthorized       = -401 // 未认证
	WebAccessDenied       = -403 // 访问权限不足
	WebNothingFound       = -404 // 啥都木有
	WebMethodNotAllowed   = -405 // 不支持该方法
	WebConflict           = -409 // 冲突
	WebServerErr          = -500 // 服务器错误
	WebServiceUnavailable = -503 // 过载保护,服务暂不可用
	WebDeadline           = -504 // 服务调用超时
	WebLimitExceed        = -509 // 超出限制

	WebFileTypeErr   = -615 // 上传文件类型错误
	WebFileNotExists = -616 // 上传文件不存在
	WebFileTooLarge  = -617 // 上传文件太大

	WebFailedTooManyTimes = -625 // 登录失败次数太多
	WebUserNotExist       = -626 // 用户不存在
	WebUserDisabled       = -627 // 用户停用
	WebPasswordTooLeak    = -628 // 密码太弱
	WebPasswordOriginal   = -629 // 原始密码

	WebTargetNumberLimit  = -632 // 操作对象数量限制
	WebTargetBlocked      = -643 // 被锁定
	WebAccessTokenExpires = -658 // Token 过期

	SqlUserNormal   = 0  //账户使用状态：正常
	SqlUserDel      = -1 //账户使用状态：删除
	SqlUserStop     = 1  //账户使用状态： 停用
	SqlUserOriginal = 2  //账户使用状态： 原始密码

	SysTokenExpiresTime   = 10 * 24 * 60 * 60 //token 失效时间，单位秒 ，10天
	SysSmsCodeExpiresTime = 2 * 60            //token 失效时间，单位秒 ，15分钟
	SysPhoneExpiresTime   = 24 * 60 * 60      //token 失效时间，单位秒 ，1天

	Enable_0 = 0 //企业微信成员启用状态。1表示启用的成员
	Enable_1 = 1 //0表示被禁用

	WechatBusUserStatus_1 = 1 //激活状态: 1=已激活
	WechatBusUserStatus_2 = 2 //激活状态: 2=已禁用
	WechatBusUserStatus_4 = 4 //激活状态: 4=未激活
)
