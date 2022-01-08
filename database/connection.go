package connection

import (
	"github.com/iamtakdir/jwt-auth-go/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//For using DB outside this file we need to have a global variable

var DB *gorm.DB

func Connect() {
	connect, err := gorm.Open(mysql.Open("root:root@/jwtauthdb"), &gorm.Config{})

	if err != nil {
		log.Fatal("Couldn't connect to the database error := ", err)
	}

	//referencing DB
	DB = connect

	connect.AutoMigrate(models.User{})

}
