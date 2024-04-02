package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"go_crud/utils/mysql_db"
	"go_crud/utils/redis_cache"
	"gorm.io/gorm"
	"log"
	"time"
)

var RDB *redis.Client
var DataBase *gorm.DB
var RedSyncLock *redsync.Redsync

func Init() {
	RDB = redis_cache.ConnectToRedis("user_redis")
	RedSyncLock = redis_cache.NewSync(redis_cache.ConnectToRedis("lock_redis"))
	var err error
	DataBase, err = mysql_db.ConnectToDatabase("user_db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = DataBase.AutoMigrate(&mysql_db.UserList{})
	if err != nil {
		fmt.Println("Error init database:", err)
		return
	}
	ClearEntireRedisCache()
}

func GetFromRedisByName(name string) []mysql_db.UserList {
	ctx := context.Background()
	var adminDataList []mysql_db.UserList
	key := fmt.Sprintf("user:%s", name)
	result, err := RDB.Get(ctx, key).Result()
	if err == nil {
		// 如果Redis中存在缓存，则直接返回
		//fmt.Println("从Redis缓存中获取数据")
		if err := json.Unmarshal([]byte(result), &adminDataList); err == nil {
			return adminDataList
		}
	}
	return nil
}

func SetNameToRedis(name string, adminDataList []mysql_db.UserList) bool {
	// 将查询结果缓存到Redis
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", name)
	data, err := json.Marshal(adminDataList)
	err = RDB.Set(ctx, key, data, 5*time.Minute).Err()
	if err != nil {
		log.Println("Redis缓存失败：", err.Error())
		return false
	}
	return true
}

func ClearNameRedisCache(name string) {
	// 连接到 Redis
	rdb := RDB
	// 删除与用户 ID 相关的缓存
	redisKey := fmt.Sprintf("user:%s", name)
	ctx := context.Background()
	_, err := rdb.Del(ctx, redisKey).Result()
	if err != nil {
		log.Printf("Failed to clear Redis cache for user ID %s: %v", name, err)
	}
}

func ClearEntireRedisCache() {
	// 使用context.Background()作为请求的上下文
	ctx := context.Background()

	log.Println("清除redis")

	// 使用FlushDB命令清空当前数据库
	err := RDB.FlushDB(ctx).Err()
	if err != nil {
		log.Printf("Failed to clear entire Redis cache: %v", err)
	} else {
		log.Println("Successfully cleared entire Redis cache")
	}
}
