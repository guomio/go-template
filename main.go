package main

import (
	"log"
	"net/http"

	"github.com/guomio/go-template/config"
	"github.com/guomio/go-template/logger"
	"github.com/guomio/go-template/model"
	"github.com/guomio/go-template/service"
)

func init() {
	log.Println("程序开始初始化...")
	log.Println("[1/4] 读取配置...")
	config.Init()
	log.Println("[2/4] 初始化日志服务...")
	logger.Init()
	log.Println("[3/4] 初始化数据库...")
	model.Init()
	log.Println("[4/4] 加载路由...")
	service.Init()
	log.Println("程序初始化完毕")
	log.Println("Running port:", config.C.Port)
}

func main() {
	http.Handle("/", service.Handler(service.DefaultCorsOption()))

	err := http.ListenAndServe(":"+config.GetConfig().Port, nil)
	if err != nil {
		logger.L.Sprintf("http serve error: %s", err)
	}
}
