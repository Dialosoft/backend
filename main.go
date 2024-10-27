/*
Welcome to the main code of the Dialosoft Forum Software project
This project is licensed under the GPL-3.0 license

For any use of the code, it will be to stay within the open source world.
follow us on github: https://github.com/Dialosoft

This source code was developed by:

  - @Flussen ( https://github.com/Flussen )
  - @matdevcoder ( https://github.com/matdevcoder )

with Golang and ❤️
*/

package main

import (
	"context"
	"log"
	"time"

	_ "github.com/Dialosoft/docs"
	"github.com/Dialosoft/src/app/config"
	"github.com/Dialosoft/src/app/database"
	"github.com/Dialosoft/src/pkg/utils/devconfig"
	"github.com/Dialosoft/src/pkg/utils/logger"
	swagger "github.com/Flussen/swagger-fiber-v3"
	"github.com/redis/go-redis/v9"
)

// @title Dialosoft Swagger API
// @version 0.2.0
// @description this is the swagger implementation for dialosoft, it is a sample docs for fiber
// @termsOfService http://swagger.io/terms/
// @contact.name @Flussen on github
// @contact.email flussen0@gmail.com
// @license.name GPL-3.0 License
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @host localhost:8080
// @BasePath /
func main() {

	var err error
	var conn database.Connection
	var redisConn *redis.Client
	logger.InitLogger()

	conf := config.GetGeneralConfig()
	if conf.Database == "" {
		log.Fatal("not variable setted")
	}

	// Database
	for {
		var count int
		conn, err = database.ConnectToDatabase(conf)
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

	// Auto destroy tokens script
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go database.StartTokenChecker(ctx, conn.Gorm, 24*time.Hour)

	// Auto create default user
	if err := devconfig.StartWithAdmin(conn.Gorm); err != nil {
		log.Fatal(err)
	}

	// Api Setup
	api := config.SetupAPI(conn.Gorm, redisConn, conf)

	// Swagger Docs
	api.Get("/docs/*", swagger.HandlerDefault)

	if err := api.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
