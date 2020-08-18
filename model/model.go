package model

import (
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/guomio/go-template/config"
	"github.com/jinzhu/gorm"

	//inject mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// M model 实例
var M *gorm.DB

// Init 初始化
func Init() {
	M = genMysqlDB()
	// 此处注册自动创建表
}

// GetModel 获取 model 实例
func GetModel() *gorm.DB {
	return M
}

func genMysqlDB() *gorm.DB {
	mysqlEnv := config.C.Mysql
	mysql, err := gorm.Open("mysql", mysqlEnv+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("配置: %s 连接mysql数据库失败，程序启动失败(%s)\n", mysqlEnv, err)
	}
	mysql.SingularTable(true)
	return mysql
}

// Model baseType
type Model struct {
	ID        int      `gorm:"primary_key;unique_index;AUTO_INCREMENT;comment:'唯一ID'" json:"id" form:"id"`
	DeletedAt JSONTime `json:"-" form:"-" gorm:"comment:'删除时间'"`
	CreatedAt JSONTime `json:"created" form:"created" gorm:"comment:'创建时间'"`
	UpdatedAt JSONTime `json:"updated" form:"updated" gorm:"comment:'修改时间'"`
}

// JSONTime 自定义数据库时间格式
type JSONTime struct {
	time.Time
}

// MarshalJSON on JSONTime format Time field with %Y-%m-%d %H:%M:%S
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// BeforeCreate manage uuid
// func (M *Model) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("ID", pure.GetUUID())
// 	return nil
// }
