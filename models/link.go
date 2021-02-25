package models

import (
	"log"
)

// Get .
func (link *Link) Get() error {
	return db.First(link).Error
}

// Put updates a link object
func (link *Link) Put(field string, value interface{}) {
	db.First(&link).Update(field, value)
}

// Create a db entry
func (link *Link) Create() error {
	log.Println("Creating link\n", link)
	return db.Create(link).Error
}

// Read a db entry by ID
func (link *Link) Read() error {
	return db.First(link).Error
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
	return db.Delete(link).Error
}
