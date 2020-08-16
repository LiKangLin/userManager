package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"myHome/userManager/code"
	"myHome/userManager/conf"
	pb "myHome/userManager/proto"
	"myHome/userManager/utils"
	"net"
)

// UserServer for rcpclient
type UserServer struct {
}

func getUUID(ctx context.Context) string {
	var uuid string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok == false {
		return uuid
	}
	uuids := md.Get("uuid")
	if len(uuids) == 1 {
		uuid = uuids[0]
	}
	return uuid
}

// Login login handler
func (server *UserServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// get uuid，从context中取出uuid
	uuid := getUUID(ctx)
	// query userinfo 通过用户名查询用户信息
	user := getUserInfo(in.Username)

	//passwd密钥匹配
	AesKey := []byte("0f90023fc9ae101e") //秘钥长度为16的倍数
	encrypted, err := AesEncrypt([]byte(in.Passwd), AesKey)
	if err != nil {
		panic(err)
	}
	//与数据库中的密码匹配
	if base64.StdEncoding.EncodeToString(encrypted) != user.Passwd {
		//logs.Error("password的值为："+in.Passwd+",user.password的值为"+user.Passwd)
		logs.Error(uuid+" ", in.Username+" -- Failed to match passwd ", in.Passwd)
		return &pb.LoginResponse{Code: code.CodeTCPPasswdErr, Msg: code.CodeMsg[code.CodeTCPPasswdErr]}, nil
	} else {
		logs.Info(uuid+" ", "-- "+in.Username+"密码匹配成功，密码为："+user.Passwd)
	}

	token := utils.GenerateToken(user.Username)
	logs.Info(uuid+" ", "-- 开始写入缓存，以"+token+"作为key值")
	// set cache 做缓存，用户名作为缓存的key值，user信息作为value值写入redis
	err1 := setTokenInfo(user, token)
	if err1 != nil {
		logs.Error(uuid, " -- Failed to set token for user:", user.Username, " err:", err.Error())
		return &pb.LoginResponse{Code: code.CodeTCPInternelErr, Msg: code.CodeMsg[code.CodeTCPInternelErr]}, nil
	}
	logs.Info(uuid+" ", "-- ls"+
		""+in.Username+" Login succesfully")
	//将生成的token值返回给前端
	return &pb.LoginResponse{Username: user.Username, Nickname: user.Nickname, Headurl: user.Headurl, Token: token, Code: code.CodeSucc}, nil
}

// GetUserInfo get user info
func (server *UserServer) GetUserInfo(ctx context.Context, in *pb.CommRequest) (*pb.LoginResponse, error) {
	// get uuid
	uuid := getUUID(ctx)
	logs.Debug(uuid, " -- GetUserInfo access from:", in.Username)
	//get userinfo and compare username  传入username作为key值查询，从缓存中获取user信息
	user, err := getTokenInfo(in.Token)
	if err != nil {
		logs.Error(uuid, " -- Failed to get key值:", " with err:", err.Error())
		return &pb.LoginResponse{Code: code.CodeTCPTokenExpired, Msg: code.CodeMsg[code.CodeTCPTokenExpired]}, nil
	}
	//logs.Debug("GetUserInfo查询缓存")
	// check if username is the same
	if user.Username != in.Username {
		logs.Error(uuid, " -- Error: token info not match:", in.Username, " while cache:", user.Username)
		return &pb.LoginResponse{Code: code.CodeTCPUserInfoNotMatch, Msg: code.CodeMsg[code.CodeTCPUserInfoNotMatch]}, nil
	}
	logs.Debug(uuid, " -- Succ to GetUserInfo :", in.Username)
	return &pb.LoginResponse{Username: user.Username, Nickname: user.Nickname, Headurl: user.Headurl, Code: code.CodeSucc}, nil
}

// EditUserInfo edit userinfo
func (server *UserServer) EditUserInfo(ctx context.Context, in *pb.EditRequest) (*pb.EditResponse, error) {
	// get uuid
	uuid := getUUID(ctx)
	logs.Debug(uuid, " -- EditUserInfo access from:", in.Username, in.Nickname)
	// auth
	authResult := auth(in.Username, in.Token)
	if authResult == false {
		logs.Error(uuid, " -- Failed to auth for user:", in.Username)
		return &pb.EditResponse{Code: code.CodeTCPTokenExpired, Msg: code.CodeMsg[code.CodeTCPTokenExpired]}, nil
	}
	affectRows := editUserInfo(in.Username, in.Nickname, in.Headurl, in.Token, in.Mode)
	logs.Info(uuid, " -- Succ to edit userinfo, affected rows is:", affectRows)
	return &pb.EditResponse{Code: code.CodeSucc, Msg: code.CodeMsg[code.CodeSucc]}, nil
}

// Logout logout
func (server *UserServer) Logout(ctx context.Context, in *pb.CommRequest) (*pb.EditResponse, error) {
	// get uuid
	uuid := getUUID(ctx)
	logs.Info(uuid+" ", "-- 删除缓存中的token值")
	//删除redis中的缓存信息
	logs.Debug(uuid, " -- Logout access from:", in.Token)
	err := delTokenInfo(in.Token)
	if err != nil {
		logs.Error(uuid, " -- Failed to delTokenInfo :", err.Error())
	}
	logs.Info(uuid+" ", "-- "+in.Username+"----logout")
	return &pb.EditResponse{Code: code.CodeSucc, Msg: code.CodeMsg[code.CodeSucc]}, nil
}

// start userserver
func start(config *conf.TCPConf) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		logs.Critical("Listen failed, err:", err.Error())
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserServer{})

	logs.Info("start to listen on localhost:%d", config.Server.Port)
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Println("Server failed, err:", err.Error())
	}
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//AES加密
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
