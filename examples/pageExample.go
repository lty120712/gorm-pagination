package main

import (
	"github.com/lty120712/gorm-pagination/pagination"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

type User struct {
	gorm.Model
	Username string
	Password string
	Nickname string
}

func main() {
	//修改 username:password
	dsn := "xxx:xxx@tcp(127.0.0.1:3306)/game?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出目标（终端）
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别（Info = 打印 SQL）
			Colorful:      true,        // 彩色打印
		},
	),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	var req pagination.PageRequest
	req.Page = 1
	req.PageSize = 2

	result := &pagination.PageResult[User]{Records: []User{}}

	// 调用分页函数并获取结果
	_, err = pagination.Paginate(Db.Model(&User{}), req.Page, req.PageSize, result)
	if err != nil {
		logrus.Error("分页查询失败:", err)
		return
	}
	logrus.Info("分页查询结果: \n", result)
}
