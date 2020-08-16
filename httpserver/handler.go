package main

import (
	codeModule "myHome/userManager/code"
	"myHome/userManager/rpcclient"
	"myHome/userManager/utils"
	"net/http"
	"path"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
)

// generate upload image file name
func generateImgName(fname, postfix string) string {
	ext := path.Ext(fname)
	fileName := strings.TrimSuffix(fname, ext)
	fileName = utils.Md5String(fileName + postfix)

	return fileName + ext
}

// login
func loginHandler(c *gin.Context) {
	username := c.PostForm("username")
	passwd := c.PostForm("passwd")
	uuid := utils.GenerateToken(username)
	logs.Info(uuid+" ", "-- loginHandler access from:", username)
	//允许跨域
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	c.Header("Access-Control-Allow-Headers", "Action,Module,X-PINGOTHER,Content-Type,Content-Disposition")
	//与grpc的服务端开始连接，参数为username、password和uuid
	ret, token, rsp := rpcclient.Login(map[string]string{"username": username, "passwd": passwd, "uuid": uuid})

	if ret == http.StatusOK && token != "" {
		c.SetCookie("token", token, config.Logic.Tokenexpire, "/", config.Server.IP, false, true)
		logs.Debug(uuid, " -- Set token ", token, "with expire:", config.Logic.Tokenexpire)
	}

	logs.Debug(uuid, " -- Succ get response from backend with", rsp["code"], " and msg:", rsp["msg"])
	c.JSON(ret, rsp)
}

// logout
func logoutHandler(c *gin.Context) {

	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	c.Header("Access-Control-Allow-Headers", "Action,Module,X-PINGOTHER,Content-Type,Content-Disposition")
	// check params
	username := c.PostForm("username")
	token := c.PostForm("logoutToken")
	uuid := utils.GenerateToken(username)
	//logs.Debug(uuid, " -- logoutHandler access from:", username)
	logs.Info(uuid+" ", "-- logoutHandler access from:", username, "，token值为："+token)
	ret, rsp := rpcclient.Logout(map[string]string{"username": username, "token": token, "uuid": uuid})
	//logs.Debug(uuid, " -- Succ to get response from backend with ", rsp["code"], " and msg:", rsp["msg"])
	c.JSON(ret, rsp)
}

// edit nickname
func editNicknameHandler(c *gin.Context) {
	// check params
	username := c.PostForm("username")
	nickname := c.PostForm("nickname")
	token := c.PostForm("editToken")
	//logs.Debug(" -- username为:", username, "new nickname:", nickname)

	uuid := utils.GenerateToken(username)
	logs.Info(uuid, " -- editNicknameHandler access from:", username, " new nickname:", nickname)
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	c.Header("Access-Control-Allow-Headers", "Action,Module,X-PINGOTHER,Content-Type,Content-Disposition")
	// communicate with rcp server
	ret, rsp := rpcclient.EditUserinfo(map[string]string{"username": username, "token": token, "nickname": nickname, "headurl": "", "mode": "1", "uuid": uuid})

	logs.Debug(uuid, " -- Succ to get response from backend with ", rsp["code"], " and msg:", rsp["msg"])
	c.JSON(ret, rsp)
}

// uploadHeadurlHandle
func uploadHeadurlHandler(c *gin.Context) {
	// check params
	username := c.Query("username")
	token := c.Query("uploadToken")

	uuid := utils.GenerateToken(username)
	logs.Info(uuid, " -- uploadHeadurlHandler access from:", username)
	// step 1 : auth
	httpCode, tcpCode, msg := rpcclient.Auth(map[string]string{"username": username, "token": token, "uuid": uuid})
	if httpCode != http.StatusOK || tcpCode != 0 {
		logs.Error(uuid, " -- uploadHeadurlHandler Auth failed, msg:", msg)
		c.JSON(httpCode, rpcclient.FormatResponse(tcpCode, msg, nil))
		return
	}
	logs.Debug(uuid, " -- uploadHeadurlHandler Auth succ")
	// step 2 : save upload picture into file
	file, image, err := c.Request.FormFile("picture")
	if err != nil {
		logs.Error(uuid, " -- Failed to FormFile, err:", err.Error())
		c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFormFileFailed, "", nil))
	}
	// check image
	if image == nil {
		logs.Error(uuid, " -- Failed to get image from formfile!")
		c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFormFileFailed, "", nil))
		return
	}
	// check filesize
	size, err := utils.GetFileSize(file)
	if err != nil {
		logs.Error(uuid, " -- Failed to get filesize, err:", err.Error())
		c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFileSizeErr, "", nil))
		return
	}
	if size == 0 || size > config.Image.Maxsize*1024*1024 {
		logs.Error(uuid, " -- Filesize illegal, size:", size)
		c.JSON(http.StatusOK, rpcclient.FormatResponse(codeModule.CodeFileSizeErr, "", nil))
		return
	}
	logs.Debug(uuid, " -- uploadHeadurlHandler CheckImage succ")
	// save
	imageName := generateImgName(image.Filename, username)
	fullPath := config.Image.Savepath + imageName

	if err = c.SaveUploadedFile(image, fullPath); err != nil {
		logs.Error(uuid, " -- Failed to save file, err:", err.Error())
		c.JSON(http.StatusInternalServerError, rpcclient.FormatResponse(codeModule.CodeInternalErr, "", nil))
		return
	}
	logs.Debug(uuid, " -- Succ to save upload image, path:", fullPath)
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	c.Header("Access-Control-Allow-Headers", "Action,Module,X-PINGOTHER,Content-Type,Content-Disposition")
	//update picture info
	imageURL := config.Image.Prefixurl + "/" + fullPath
	ret, editRsp := rpcclient.EditUserinfo(map[string]string{"username": username, "token": token, "nickname": "", "headurl": imageURL, "mode": "2", "uuid": uuid})
	logs.Debug(uuid, " -- editUserInfo response:", ret)
	c.JSON(ret, editRsp)
}

// get user info
func getUserinfoHandler(c *gin.Context) {
	// check params
	username := c.Query("username")

	token := c.Query("userToken")
	//logs.Info("access from:", username, " with token:", token)
	//logs.Debug("access from:", username)

	uuid := utils.GenerateToken(username)
	logs.Debug(uuid, " -- getUserinfoHandler access from:", username)
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
	c.Header("Access-Control-Allow-Headers", "Action,Module,X-PINGOTHER,Content-Type,Content-Disposition")
	// communicate with rcp server
	ret, rsp := rpcclient.GetUserinfo(map[string]string{"username": username, "token": token, "uuid": uuid})
	logs.Debug(uuid, " -- Succ to get response from backend with ", rsp["code"], " and msg:", rsp["msg"])
	c.JSON(ret, rsp)
}
