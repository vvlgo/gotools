package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data

}

// StringsJoin 字符串拼接
func StringsJoin(strs ...string) string {
	var str string
	var b bytes.Buffer
	strsLen := len(strs)
	if strsLen == 0 {
		return str
	}
	for i := 0; i < strsLen; i++ {
		b.WriteString(strs[i])
	}
	str = b.String()
	return str

}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*
复制map函数
*/
func CopyMap(oldMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})

	// Copy from the original map to the target map
	for key, value := range oldMap {
		newMap[key] = value
	}
	return newMap
}

/*
将各种类型转换为字符串
*/
func TypeToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

/*
计算时间，返回time，输入参数xxxx-xx-xx xx:xx:xx
*/
func CalculateTime(stime string, year, month, day int) time.Time {
	ntime, err := time.Parse("2006-01-02 15:04:05", strings.Replace(stime, "\\", "", -1))
	if err != nil {
		log.Error("string to time err,", err)
	}
	caltime := ntime.AddDate(year, month, day)
	return caltime
}

/*
xxxx-xx-xx字符串转日期
*/
func StrToDay(stime string) time.Time {
	stime = strings.Replace(stime, "\\", "", -1)

	ntime, err := time.Parse("2006-01-02", stime[:10])
	if err != nil {
		log.Error("string to day err,", err)
	}
	return ntime
}

/*
xxxx-xx字符串转日期
*/
func StrToMonth(stime string) time.Time {
	stime = strings.Replace(stime, "\\", "", -1)
	ntime, err := time.Parse("2006-01", stime[:7])
	if err != nil {
		log.Error("string to month err,", err)
	}
	return ntime
}

/*
提取员工编号的数字,并生成新的编号
*/
func GetMnum(str *string, i int, flag bool) *string {
	runes := []rune(*str)
	var s string
	var s2 string
	var num int
	for k, v := range runes {
		if v > 48 && v <= 57 {
			s = string(runes[k:])
			num, _ = strconv.Atoi(s)

			if flag {
				num = num + 1
			}
			n := len(strconv.Itoa(num+i)) - len(s)
			if num+i < 10 {
				s2 = string(runes[:k])
			}

			if num+i >= 10 && num+i < 100 {
				s2 = string(runes[:k-n])
			}
			if num+i >= 100 && num+i < 1000 {
				s2 = string(runes[:k-n])
			}
			if num+i >= 1000 && num+i < 10000 {
				s2 = string(runes[:k-n])
			}
			if num+i >= 10000 && num+i < 100000 {
				s2 = string(runes[:k-n])
			}
			break
		}
	}
	new := ""
	if num+i >= 100000 {
		new = "M" + strconv.Itoa(num+i)
	} else {
		new = s2 + strconv.Itoa(num+i)
	}
	return &new
}

/*
对计算后的浮点数处理，n保留几位小数
*/
func Decimal(value float64, n string) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+n+"f", value), 64)
	return value
}

/*
发送邮件
*/
func SendEmail(to []string, title, context string) {
	for i := 0; i < len(to); i++ {
		go func(i int) {
			m := gomail.NewMessage()
			email := ""
			emailpwd := ""
			m.SetAddressHeader("From", email, "ERP系统")
			m.SetHeader("To", to[i])
			m.SetHeader("Subject", title)
			m.SetBody("text/html", context)
			d := gomail.NewDialer("smtp.exmail.qq.com", 465, email, emailpwd)

			if err := d.DialAndSend(m); err != nil {
				log.Error("send "+to[i]+" email err,", err)
			}

			log.Info("done.发送成功")
		}(i)
	}

}

/*
float64 to string
*/
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

/*
float32 to string
*/
func Float32ToString(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

/*
计算两个时间之间相差的天数,n小数位数
*/
func TimeDValue(startTime, endTime time.Time, n string) float64 {

	return Decimal(float64(endTime.Unix()-startTime.Unix())/(60.00*60*24), n)
}

/*
md5加密
*/
func Md5PWD(pwd string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
}

/*
计算所提供的日期的月份最后一天
*/
func MonthLastDay(t time.Time) (start, last string) {
	const DATE_FORMAT = "2006-01-02"
	year, month, _ := t.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	start = thisMonth.AddDate(0, 0, 0).Format(DATE_FORMAT)
	last = thisMonth.AddDate(0, 1, -1).Format(DATE_FORMAT)
	return
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

/*
10位字符长度的uuid
*/
func GetUUID(deviceId, offerCode string) string {
	uuids := uuid.UUID{}
	v3 := uuid.NewV3(uuids, deviceId+offerCode)
	ss := strings.Split(v3.String(), "-")
	s := ss[0] + ss[1]

	return s[:10]
}

/*
TimeToDate 输入时间，结果xxxx-xx-xx
*/
func TimeToDate(date time.Time) string {
	return date.Format("2006-01-02")
}

/*
TimeToTimestamp 输入时间，结果xxxx-xx-xx xx:xx:xx:
*/
func TimeToTimestamp(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

/*
GetOrderNum 计算订单号，订单号格式yyyymmddDDmmss+电话后四位+5位随机数
*/
func GetOrderNum(phone string) string {
	ntime := time.Now()
	s := ntime.Format("20060102150405")
	phone = phone[7:]
	s = s + phone
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var result = ""
	for i := 0; i < 5; i++ {
		result += strconv.Itoa(r.Intn(10))
	}
	s = s + result
	return s
}

/*
CreateDir 创建文件夹
*/
func CreateDir(dir string) error {
	b, _ := PathExists(dir)
	if !b {
		err := os.Mkdir(dir, 07777)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
