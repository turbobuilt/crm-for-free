package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var db *redis.Client

// var sdb *redisearch.Client

func GetDB() *redis.Client {
	if db != nil {
		return db
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	db = rdb

	// sdb = redisearch.NewClient("localhost:6379", "")

	return db
}
