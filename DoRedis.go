package main

import (
	"fmt"
	"strconv"
	"github.com/go-redis/redis"
)

func RedisConnect() (*redis.Client, int32) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return client, -1
	}
	return client, 1
}

func RedisAddEntity(client *redis.Client, redis_package InternetPackage) int32 {
	key := strconv.Itoa(int(redis_package.id))

	if client.Exists(key).Val() != 0 {
		fmt.Println("Already exists !")
		return -1
	}

	infomap := make(map[string]interface{})
	infomap["ptype"] = redis_package.ptype
	infomap["state"] = redis_package.state
	infomap["reserve"] = redis_package.reserve
	infomap["X"] = redis_package.X
	infomap["Y"] = redis_package.Y
	infomap["Z"] = redis_package.Z
	infomap["toX"] = redis_package.toX
	infomap["toY"] = redis_package.toY
	infomap["toZ"] = redis_package.toZ

	err := client.HMSet(key, infomap).Err()
	if err != nil {
		fmt.Println(err)
		return -2
	}
	return 1
}

func RedisGetEntity(client *redis.Client, key string) (map[string]string, int32) {
	var val  map[string]string

	if client.Exists(key).Val() == 0 {
		fmt.Println("Not exists !")
		return val, -1
	}

	val, err := client.HGetAll(key).Result()
	if err != nil {
		fmt.Println(err)
		return val, -2
	}

	return val, 1
}

func RedisChangeEntity(client *redis.Client, redis_package InternetPackage) int32 {
	key := strconv.Itoa(int(redis_package.id))
	if client.Exists(key).Val() == 0 {
		fmt.Println("Not exists !")
		return -1
	}

	infomap := make(map[string]interface{})
    infomap["ptype"] = redis_package.ptype
    infomap["state"] = redis_package.state
    infomap["reserve"] = redis_package.reserve
    infomap["X"] = redis_package.X
    infomap["Y"] = redis_package.Y
    infomap["Z"] = redis_package.Z
    infomap["toX"] = redis_package.toX
    infomap["toY"] = redis_package.toY
    infomap["toZ"] = redis_package.toZ

	err := client.HMSet(key, infomap).Err()
	if err != nil {
		fmt.Println(err)
		return -2
	}

	return 1
}

func RedisDeleteEntity(client *redis.Client, key string) int32 {
    if client.Exists(key).Val() == 0 {
        fmt.Println("Not exists !")
        return -1
    }

	err := client.Del(key).Err()
	if err != nil {
		fmt.Println(err)
		return -2
	}

	return 1
}

func RedisGetAllEntityKeys(client *redis.Client) ([]string, int32) {
	val, err := client.Keys("*").Result()
	if err != nil {
		fmt.Println(err)
		return val, -1
	}
	return val, 1
}
