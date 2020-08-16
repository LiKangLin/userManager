package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"myHome/userManager/conf"
	"time"
)

var redisConn *redis.Client

// init redis connection pool
func initRedisConn(conf *conf.TCPConf) error {
	redisConn = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Passwd,
		DB:       conf.Redis.Db,
		PoolSize: conf.Redis.Poolsize,
	})
	if redisConn == nil {
		return errors.New("Failed to call redis.NewClient")
	}

	_, err := redisConn.Ping().Result()
	if err != nil {
		msg := fmt.Sprintf("Failed to ping redis, err:%s", err.Error())
		return errors.New(msg)
	}
	return nil
}

// cleanup
func closeCache() {
	redisConn.Close()
}

// get cached userinfo
func getUserCacheInfo(username string) (User, error) {

	AesKey := []byte("efd0023fc9ae1bbb") //秘钥长度为16的倍数
	encrypted, err := AesEncrypt([]byte(username), AesKey)
	if err != nil {
		panic(err)
	}
	//logs.Debug("查询缓存中的信息加密后: %s\n", base64.StdEncoding.EncodeToString(encrypted))
	redisKey := base64.StdEncoding.EncodeToString(encrypted)
	val, err := redisConn.Get(redisKey).Result()
	var user User
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

// set cached userinfo
func setUserCacheInfo(user User) error {
	AesKey := []byte("efd0023fc9ae1bbb") //秘钥长度为16的倍数
	//AesKey := []byte("0f90023fc9ae101e") //秘钥长度为16的倍数
	encrypted, err := AesEncrypt([]byte(user.Username), AesKey)
	if err != nil {
		panic(err)
	}
	//logs.Debug("查询缓存中的信息加密后: %s\n", base64.StdEncoding.EncodeToString(encrypted))

	redisKey := base64.StdEncoding.EncodeToString(encrypted)
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.Redis.Cache.Userexpired)
	_, err = redisConn.Set(redisKey, val, expired).Result()
	return err
}

// get token info
func getTokenInfo(token string) (User, error) {
	redisKey := token
	val, err := redisConn.Get(redisKey).Result()
	var user User
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

//连接redis,token值作为key值，user的信息作为value值
func setTokenInfo(user User, token string) error {
	//将加密后的token值写入redis
	redisKey := token
	//redisKey := tokenKeyPrefix + username
	//logs.Debug(redisKey+"为该用户的redis中的key值")
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.Redis.Cache.Tokenexpired)
	_, err = redisConn.Set(redisKey, val, expired).Result()
	return err
}

// update cached userinfo, if failed, try to delete it from cache
func updateCachedUserinfo(user User) error {
	AesKey := []byte("efd0023fc9ae1bbb") //秘钥长度为16的倍数
	encrypted, err := AesEncrypt([]byte(user.Passwd), AesKey)
	if err != nil {
		panic(err)
	}
	//logs.Debug("查询缓存中的信息加密后: %s\n", base64.StdEncoding.EncodeToString(encrypted))
	err1 := setUserCacheInfo(user)
	if err1 != nil {
		redisKey := base64.StdEncoding.EncodeToString(encrypted)
		redisConn.Del(redisKey).Result()
	}
	return err1
}

// delete token cache info
//根据key值删除reids中的缓存
func delTokenInfo(token string) error {

	////密钥匹配
	//AesKey := []byte("0f90023fc9ae1cdm") //秘钥长度为16的倍数
	////fmt.Printf("明文: %s\n秘钥: %s\n", in.Passwd, string(AesKey))
	//encrypted, err := AesEncrypt([]byte(token), AesKey)
	//if err != nil {
	//    panic(err)
	//}
	//logs.Debug("加密后: %s\n", base64.StdEncoding.EncodeToString(encrypted))

	redisKey := token
	_, err1 := redisConn.Del(redisKey).Result()
	return err1
}
