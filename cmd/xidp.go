package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/colin-404/logx"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xid-protocol/xidp/biz"
	"github.com/xid-protocol/xidp/db"
)

var sig = make(chan os.Signal, 1)

func initConfig() string {
	confPath := flag.String("c", "/opt/xidp/conf/config.yml", "config file path")
	flag.Parse()
	// confPath_str := common.NormalizePath(*confPath)
	//如果配置文件不存在，则报错并关闭程序
	if _, err := os.Stat(*confPath); os.IsNotExist(err) {
		logx.Errorf("config file not found: %s", *confPath)
		os.Exit(1)
	}
	//加载配置
	return *confPath
}

func initLog() {
	//log配置为空，则退出
	logPath := viper.GetString("Log.path")
	maxSize := viper.GetInt("Log.max_size")
	maxAge := viper.GetInt("Log.max_age")
	maxBackups := viper.GetInt("Log.max_backups")

	if logPath == "" || maxSize == 0 || maxAge == 0 || maxBackups == 0 {
		logx.Errorf("log.path, log.max_size, log.max_age, log.max_backups is required")
		os.Exit(1)
	}
	//初始化日志
	logOpts := logx.Options{
		LogFile:    viper.GetString("Log.path"),
		MaxSize:    viper.GetInt("Log.max_size"),
		MaxAge:     viper.GetInt("Log.max_age"),
		MaxBackups: viper.GetInt("Log.max_backups"),
		TimeFormat: logx.TimeFormats.RFC3339,
	}
	loger := logx.NewLoger(&logOpts)
	logx.InitLogger(loger)
}

func init() {
	//优雅关闭
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	//获取配置路径
	confPath := initConfig()

	//使用viper加载配置
	viper.SetConfigFile(confPath)
	viper.ReadInConfig()

	initLog()

	// 初始化MongoDB连接
	err := db.InitMongoDB()
	if err != nil {
		logx.Fatalf("Failed to initialize MongoDB: %v", err)
	}
}

func main() {
	defer func() {
		// 关闭数据库连接
		if err := db.CloseMongoDB(); err != nil {
			logx.Errorf("Failed to close MongoDB connection: %v", err)
		}
	}()

	go ServerStart()
	//go sealsuite.SealsuiteAcountInit()
	//go accounts.AccountMonitor()
	<-sig
}

func ServerStart() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	biz.RegisterRouter(router)

	//获取端口配置，如果获取不到，则退出
	port := viper.GetInt("Server.port")
	if port == 0 {
		logx.Errorf("server.port is not set")
		os.Exit(1)
	}

	logx.Infof("Listening and serving on %d", port)
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logx.Errorf("SRV_ERROR %s", err.Error())
	}

	<-sig
}
