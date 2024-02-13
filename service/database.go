package service

import "github.com/kamalesh-seervi/simpleGPT/models"

var db *gorm.DB

func InitDB() {
	// Open a database connection
	var err error
	db, err = gorm.Open("sqlite3", "test.db") // Replace with your database configuration
	if err != nil {
		panic("Failed to connect to the database")
	}

	// AutoMigrate creates tables based on the User model
	db.AutoMigrate(&models.User{})
}

func CreateUser(user *models.User) error {
	return db.Create(user).Error
}

// GetUserByUsername retrieves a user by their username
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateAPIKey updates the API key for a user
func UpdateAPIKey(user *models.User, apiKey string) error {
	return db.Model(user).Update("api_key", apiKey).Error
}
