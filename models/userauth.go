package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

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
	err := db.First(user).Error
	if err != nil {
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
	return db.First(&user).Update("enabled", true).Error
}
