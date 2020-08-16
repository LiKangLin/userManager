package main

import (
	"fmt"
	"math/rand"
	"myHome/userManager/rpcclient"
	"myHome/userManager/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
)

// login
func randomLoginHandler(c *gin.Context) {
	// check params
	uid := rand.Int63n(10000000)
	username := fmt.Sprintf("username%d", uid)
	//logs.Debug("Random中username的值为："+username)
	passwd := "123456"
	//logs.Debug("Random中passwd的值为："+passwd)

	//根据用户名生成一个uuid
	uuid := utils.GenerateToken(username)
	logs.Debug(uuid, " --hahloginHandler access from:", username, "@", passwd)

	//允许跨域
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	c.Header("Access-Control-Allow-Headers", "Action,Module,X-PINGOTHER,Content-Type,Content-Disposition")
	// communicate with rcp server
	//与rcp的服务端开始连接，参数为username、password和uuid
	ret, token, rsp := rpcclient.Login(map[string]string{"username": username, "passwd": passwd, "uuid": uuid})
	// set cookieMD5 将cookie保存在header中
	if ret == http.StatusOK && token != "" {
		c.SetCookie("token", token, config.Logic.Tokenexpire, "/", config.Server.IP, false, true)
		logs.Debug(uuid, " -- Set token ", token, "with expire:", config.Logic.Tokenexpire)
	}

	logs.Debug(uuid, " -- Succ get response from backend with", rsp["code"], " and msg:", rsp["msg"])
	c.JSON(ret, rsp)
}
