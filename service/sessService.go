package service

import (
	"chat/service/dao"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

type sessService struct{}

// User Token Validate for 8 hours
func (s *sessService) RedisSaveUser(token string, user dao.User) {
	userData := map[string]interface{}{"UserName": user.UserName, "PassWord": user.PassWord}
	redisClient.HMSet(token, userData)
	redisClient.Expire(token, time.Duration(time.Hour*8))
}

// 检查token和User是否对应
func (s *sessService) RedisCheckAuth(token string, user dao.User) bool {
	userData, err := redisClient.HGetAll(token).Result()
	if err != nil {
		log.Println("在HGETALL时发生error ", err)
	}
	if userName, ok := userData["userName"]; !ok {
		return false
	} else if userName == user.UserName {
		return true
	}
	return false
}

func (s *sessService) RedisCheckUserByToken(token string) dao.User {
	userData, err := redisClient.HGetAll(token).Result()
	if err != nil {
		log.Println("在HGETALL时发生error ", err)
	}
	return dao.User{
		UserName: userData["UserName"],
	}

}
