package example_test

import (
	"fmt"
	conf2 "github.com/vvlgo/gotools/yamlconf"
	"testing"
)

var (
	conf *Config
)

type Config struct {
	RunModel     string           `yaml:"runmode"`
	Port         string           `yaml:"port"`
	Gormlog      bool             `yaml:"gormlog"`
	CertFile     string           `yaml:"cert"`
	KeyFile      string           `yaml:"key"`
	IP           string           `yaml:"ip"`
	InviterPhone []string         `yaml:"inviterphone"`
	PdfPath      string           `yaml:"pdfpath"`
	PdfTemp      string           `yaml:"pdftemp"`
	WechatToUser string           `yaml:"wechattouser"`
	Redis        RedisConfig      `yaml:"redis"`
	DataBase     DataBase         `yaml:"db"`
	MgDataBase   MgDataBase       `yaml:"mgdb"`
	BusWechat    TencentBusWechat `yaml:"buswechat"`
	Sms          TencentSMS       `yaml:"sms"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

type DataBase struct {
	Type     string `yaml:"type"`
	DbName   string `yaml:"dbname"`
	Addr     string `yaml:"addr"`
	UserName string `yaml:"username"`
	PassWord string `yaml:"password"`
}
type MgDataBase struct {
	Addrs      []string `yaml:"addrs"`
	Datebase   string   `yaml:"datebase"`
	Source     string   `yaml:"source"`
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
	Collection string   `yaml:"collection"`
}

type TencentBusWechat struct {
	Corpid  string `yaml:"corpid"`
	AgentId string `yaml:"agentId"`
}
type TencentSMS struct {
	SdkappID string `yaml:"sdkappid"`
	Appkey   string `yaml:"appkey"`
	TplID    int    `yaml:"tplid"`
}

func GetConfig() *Config {
	return conf
}

func TestConf(t *testing.T) {
	config := Config{}
	err := conf2.Read("dev.yaml", &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Port)
}
