package models

import (
	"gorm.io/gorm"

	"github.com/gal/tny/database"
	"github.com/gofrs/uuid"
)

var db *gorm.DB

// InitModels migrates modesls and initiates database connection
func InitModels() {
	db = database.InitDB()

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";") // enable uuid generation on server
	db.AutoMigrate(&Link{}, &Visit{})
}

// User structure containing non-authenticative information
type User struct {
	UID         string `gorm:"primaryKey" json:"uid,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	CreatedAt   int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

// Link structure containing slug/URL pairs
type Link struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id,omitempty"`
	OwnerID     string    `json:"owner_id"`
	Slug        string    `gorm:"unique" json:"slug,omitempty"`
	URL         string    `json:"url,omitempty"`
	DateCreated int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	Lease       int64     `json:"lease,omitempty"`
}

// Visit structure containing time of visitation
type Visit struct {
	LinkID uuid.UUID `json:"link_id"`
	Time   int64     `gorm:"autoCreateTime" json:"time,omitempty"`
}
