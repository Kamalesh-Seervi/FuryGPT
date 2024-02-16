package service

import (
	"fmt"

	"github.com/kamalesh-seervi/simpleGPT/models"
	"gorm.io/driver/postgres" 
	"gorm.io/gorm"
)

var db *gorm.DB

func DBSetup() {
	conn := "host=db port=5432 user=root password=changeme dbname=mydb "
	fmt.Println(conn)
	var err error

	db, err = gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println(err)
	}
}
// func CreateUser(user *models.User) error {
// 	return db.Create(user).Error
// }

// func GetUserByUsername(username string) (*models.User, error) {
// 	var user models.User
// 	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func UpdateAPIKey(user *models.User, apiKey string) error {
// 	return db.Model(user).Update("api_key", apiKey).Error
// }
