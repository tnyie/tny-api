package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB returns a database connection
func InitDB() *gorm.DB {
	fmt.Println(viper.GetString("db.pass"))
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Europe/Dublin",
		viper.GetString("db.host"), viper.GetString("db.user"),
		viper.GetString("db.pass"), viper.GetString("db.name"),
		viper.GetString("db.port"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("couldn't connect to database", err)
	}
	return db
}
