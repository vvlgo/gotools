package mysqldb

import (
	//mysql数据驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var (
	//DB 数据库
	DB *gorm.DB
)

/*
OpenDB 开启数据库链接
*/
func ConnectDB(dbType string, username string, password string, addr string, dbname string, modeType string) {
	var err error

	DB, err = gorm.Open(dbType, username+":"+password+"@tcp("+addr+")/"+dbname+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.SingularTable(true)
	DB.LogMode(true)
	if modeType == "pro" {
		DB.SetLogger(log.StandardLogger())
	}

}

//DBMethods 数据库基础操作接口
type DBMethods interface {
	Insert() bool
	InsertList() bool
	UpdateById() bool
	Delete() bool
	Select()
}
