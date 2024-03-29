package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Get userAuth data
func (user *UserAuth) Get() error {
	return db.First(&user, "uid = ?", user.UID).Error
}

func (user *UserAuth) GetByEmail() error {
	return db.Where("email = ?", user.Email).First(&user).Error
}

// HashPassword returns a hashed version of given password
func HashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword compares provided plaintext password against stored hash
func (user *UserAuth) VerifyPassword(password string) error {
	// fetch user
	err := db.First(&user, "uid = ?", user.UID).Error
	if err != nil {
		log.Println("Failed to fetch user")
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
}

// Create user object
func (user *UserAuth) Create(password string) error {
	log.Println("Creating user credentials\n", user)
	var err error
	user.Hash, err = HashPassword([]byte(password))
	if err != nil {
		log.Println("couldn't hash password")
		return err
	}
	return db.Create(user).Error
}

// Verify sets 'enabled' field to true
func (user *UserAuth) Verify() error {
	return db.Model(&user).Update("enabled", true).Error
}

func (user *UserAuth) ChangePassword(password string) error {
	hash, err := HashPassword([]byte(password))
	if err != nil {
		log.Println("error hashing password")
		return err
	}

	user.Hash = hash
	err = db.Model(&user).Update("hash", user.Hash).Error

	return err
}
