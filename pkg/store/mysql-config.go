package store

import (
	"fmt"
	"log"
	"os"

	"github.com/johan-ag/wishlist/internal/users"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic(".env file not found")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	//user:pass@tcp(127.0.0.1:3306)/dbName?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(db)
		log.Println(err)
		panic("Error to connect db")
	}

	err = db.AutoMigrate(&users.User{})
	if err != nil {
		panic("Fail migrations")
	}
	return db
}
