package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"myHome/userManager/conf"
	"myHome/userManager/rpcclient"
	"net/http"
	"os"
	"time"
)

var config conf.HTTPConf

func init() {
	// parser config
	var confFile string
	flag.StringVar(&confFile, "c", "../conf/httpserver.yaml", "config file")
	flag.Parse()

	err := conf.ConfParser(confFile, &config)
	if err != nil {
		logs.Critical("Parser config failed, err:", err.Error())
		os.Exit(-1)
	}

	// init log
	logConfig := fmt.Sprintf(`{"filename":"%s","level":%s,"maxlines":0,"maxsize":0,"daily":true,"maxdays":%s}`,
		config.Log.Logfile, config.Log.Loglevel, config.Log.Maxdays)
	logs.SetLogger(logs.AdapterFile, logConfig)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.Async()

	// 初始化grpc客户端连接池
	err = rpcclient.InitPool(config.Rpcserver.Addr, config.Pool.Initsize, config.Pool.Capacity, time.Duration(config.Pool.Maxidle)*time.Second)
	if err != nil {
		logs.Critical("InitPool failed, err:", err.Error())
		os.Exit(-2)
	}
}

// cleanup global objects
func finalize() {
	rpcclient.DestoryPool()
}

func main() {
	//可以使用关键字defer向函数注册退出调用，即主函数退出时，
	//defer后的函数才被调用。defer语句的作用是不管程序是否
	//出现异常，均在函数退出时自动执行相关代码。
	defer finalize()
	//gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode
	gin.SetMode(gin.DebugMode)
	//创建IO流写文件
	gin.DefaultWriter = ioutil.Discard

	route := gin.Default()
	route.Any("/welcome", webRoot)
	route.POST("/login", loginHandler)
	route.POST("/logout", logoutHandler)
	route.GET("/getuserinfo", getUserinfoHandler)
	route.POST("/editnickname", editNicknameHandler)
	route.POST("/uploadpic", uploadHeadurlHandler)
	route.GET("/randlogin", randomLoginHandler)
	route.Static("/static/", "./static/")
	route.Static("/upload/images/", "./upload/images/")

	route.Run(fmt.Sprintf(":%d", config.Server.Port))
}

func webRoot(context *gin.Context) {
	context.String(http.StatusOK, "hello, world")
}
