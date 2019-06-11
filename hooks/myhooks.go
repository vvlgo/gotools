package hooks

import (
	"encoding/json"
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/vvlgo/gotools/utils"
	"github.com/weekface/mgorus"
	"gopkg.in/mgo.v2"
	"log"

	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// ContextHook for log the call context
type contextHook struct {
	Field  string
	Skip   int
	levels []logrus.Level
}

// NewContextHook use to make an hook
// 根据上面的推断, 我们递归深度可以设置到5即可.
func NewContextHook(levels ...logrus.Level) logrus.Hook {
	hook := contextHook{
		Field:  "line",
		Skip:   5,
		levels: levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

// Levels implement levels
func (hook contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implement fire
func (hook contextHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = findCaller(hook.Skip)
	return nil
}

// 对caller进行递归查询, 直到找到非logrus包产生的第一个调用.
// 因为filename我获取到了上层目录名, 因此所有logrus包的调用的文件名都是 logrus/...
// 因此通过排除logrus开头的文件名, 就可以排除所有logrus包的自己的函数调用
func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// 这里其实可以获取函数名称的: fnName := runtime.FuncForPC(pc).Name()
// 但是我觉得有 文件名和行号就够定位问题, 因此忽略了caller返回的第一个值:pc
// 在标准库log里面我们可以选择记录文件的全路径或者文件名, 但是在使用过程成并发最合适的,
// 因为文件的全路径往往很长, 而文件名在多个包中往往有重复, 因此这里选择多取一层, 取到文件所在的上层目录那层.
func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	//fmt.Println(file)
	//fmt.Println(line)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}

//ConfigLocalFilesystemLogger config logrus log to local filesystem, with file rotation
func ConfigLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) *lfshook.LfsHook {
	baseLogPaht := path.Join(logPath, logFileName)
	b, _ := utils.PathExists(baseLogPaht)
	if !b {
		file, err := os.Create(baseLogPaht)
		logrus.Error(err)
		file.Close()
	}
	writer, err := rotatelogs.New(
		logPath+"/%Y-%m-%d."+logFileName,
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
	return lfHook
}

func MongodbHooks(addrs []string, database, source, username, password, collection string) logrus.Hook {

	dialInfo := &mgo.DialInfo{
		Addrs:     addrs, //远程(或本地)服务器地址及端口号
		Direct:    false,
		Timeout:   time.Second * 5,
		Database:  database, //数据库
		Source:    source,
		Username:  username,
		Password:  password,
		PoolLimit: 4096, // Session.SetPoolLimit
	}
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalf("can't create session: %s\n", err)
	}
	c := s.DB(database).C(collection)
	hooker := mgorus.NewHookerFromCollection(c)
	return hooker
}

type LogOption struct {
	IsMgLog       bool          //是否monogodb数据库记录
	Addrs         []string      //如果开启mg就要设置addr，ip:port
	Database      string        //数据库名字
	Source        string        //数据库名字
	Username      string        //如果开启安全连接，用户名
	Password      string        //如果开启安全连接，密码
	Collection    string        //所选表
	Logpath       string        //日志路径 PS：logs/
	Logname       string        //日志名字 PS：log.log
	RetentionTime time.Duration //保存时间
	CutTime       time.Duration //切割时间
}

func LogConf(modeType string, option *LogOption) {
	logrus.AddHook(NewContextHook())
	if strings.Contains(modeType, "dev") {
		//logrus.AddHook(MongodbHooks(yamlconf.Config.MgDataBase))
	} else {
		if option.IsMgLog {
			logrus.AddHook(MongodbHooks(option.Addrs, option.Database, option.Source, option.Username, option.Password, option.Collection))
		} else {
			logrus.AddHook(ConfigLocalFilesystemLogger(option.Logpath, option.Logname, option.RetentionTime*time.Hour, option.CutTime*time.Hour))
		}
	}

}

//mg记录字段，随需求变化

func ToMap(v interface{}) map[string]interface{} {
	bytes, _ := json.Marshal(v)
	mp := make(map[string]interface{})
	json.Unmarshal(bytes, &mp)
	return mp
}
