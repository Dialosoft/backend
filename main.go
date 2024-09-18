package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Dialosoft/src/app/config"
	"github.com/Dialosoft/src/app/database"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func main() {

	var err error
	var db *gorm.DB
	var redisConn *redis.Client

	conf := config.GetGeneralConfig()
	if conf.Database == "" {
		log.Fatal("not variable setted")
	}

	fmt.Println(conf)

	// Database
	for {
		var count int
		db, err = database.ConnectToDatabase(conf)
		if err == nil {
			break
		} else {
			count++
			time.Sleep(3 * time.Second)
		}
	}

	for {
		redisConn = database.NewRedisClient()
		if redisConn != nil {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}

	// Api Setup
	api := config.SetupAPI(db, redisConn, conf)

	if err := api.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
