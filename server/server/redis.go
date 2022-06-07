package server

import (
	"encoding/json"
	"fmt"
	"server/model"
	"server/utils"
	"time"

	"github.com/garyburd/redigo/redis"
)

var redisPool *redis.Pool

func InitRedisService() error {
	redisPool = &redis.Pool{
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:6379") },
		MaxIdle:     50,              //最大空闲连接数
		MaxActive:   0,               //表示和数据库的最大连接数，0表示不限制
		IdleTimeout: time.Minute * 1, //最大空闲时间，类型为time.Duration
	}
	_, err := redisPool.Dial()
	if err != nil {
		fmt.Println("func InitRedisService,初始化redis连接池失败")
		return err
	}
	redisConn := redisPool.Get()
	defer redisConn.Close()
	testUser := model.User{
		UserId:       0,
		UserName:     "哆啦A梦",
		UserPwd:      "123456",
		Sex:          "男",
		RegisterTime: time.Now().String(),
	}
	testUser.UserPwd = utils.StrEncrypt(testUser.UserPwd)
	data, err := json.Marshal(testUser)
	if err != nil {
		fmt.Println("in func InitRedisService,json.Marshal failed")
		return err
	}
	_, err = redisConn.Do("hset", "users", testUser.UserId, string(data))
	if err != nil {
		fmt.Println("in func InitRedisService,redisConn.Do failed")
		return err
	}
	return nil
}
