package models

import (
	"fmt"
	"log"
)

// Get .
func (link *Link) Get() error {
	return db.First(&link, "id = ?", link.ID).Error
}

// Get link by slug
func (link *Link) GetBySlug() error {
	return db.Where("slug = ?", link.Slug).First(&link).Error
}

// GetLinksByUser eturns each link owned by given user id
func GetLinksByUser(id string) (*[]Link, error) {
	var links []Link
	err := db.Where("owner_id = ?", id).Find(&links).Error
	log.Println(links)
	return &links, err
}

// Put updates a link object
func (link *Link) Put(uid string, field string, value interface{}) error {
	link.Get()
	if link.OwnerID != uid {
		return fmt.Errorf("User doesn't own resource")
	}
	log.Println("Updating link")
	return db.First(link).Update(field, value).Error
}

// Create a db entry
func (link *Link) Create() error {
	log.Println("Creating link\n", link)
	return db.Create(link).Error
}

// Read a db entry by ID
func (link *Link) Read() error {
	return db.First(&link, "id = ?", link.ID).Error
}

// Search for db entry by slug
func (link *Link) Search() error {
	return db.Where("slug=?", link.Slug).First(link).Error
}

// Update a db entry by ID
func (link *Link) Update() error {
	return db.Save(link).Error
}

// Delete a db entry2
func (link *Link) Delete() error {
	db.First(link)
	return db.Delete(link).Error
}
