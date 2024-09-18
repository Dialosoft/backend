package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Dialosoft/src/app/config"
	"github.com/Dialosoft/src/app/database"
	"gorm.io/gorm"
)

func main() {

	var err error
	var db *gorm.DB

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

	// Api Setup
	api := config.SetupAPI(db, conf)

	if err := api.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
