package models

import "log"

// Get user object
func (user *User) Get() error {
	log.Println("getting user\n", user.UID)
	return db.First(&user, user.UID).Error
}

// Create user object
func (user *User) Create() error {
	log.Println("Creating user\n", user)
	return db.Create(user).Error
}

// Put updates a user object
func (user *User) Put(field string, value interface{}) {
	db.First(user).Update(field, value)
}

// Delete deletes a user's db entry
func (user *User) Delete() error {
	return db.Delete(user).Error
}
