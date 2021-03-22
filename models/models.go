package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"

	"github.com/tnyie/tny-api/database"
)

var db *gorm.DB

// InitModels migrates modesls and initiates database connection
func InitModels() {
	db = database.InitDB()

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";") // enable uuid generation on server
	db.AutoMigrate(&User{}, &UserAuth{}, &Link{}, &Visit{})
}

// User structure containing non-authenticative information
type User struct {
	UID       string `gorm:"primaryKey" json:"uid,omitempty"`
	Username  string `json:"display_name,omitempty"`
	Email     string `gorm:"unique" json:"email,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}

// UserAuth contains credentials stored serverside
type UserAuth struct {
	UID      string `gorm:"primaryKey;type:string;default:uuid_generate_v4()" json:"uid,omitempty"`
	Username string `gorm:"unique" json:"username,omitempty"`
	Email    string `gorm:"unique" json:"email,omitempty"`
	Hash     string `json:"hash,omitempty"`
	Enabled  bool   `json:"enabled"`
}

// Link structure containing slug/URL pairs
type Link struct {
	ID          string `gorm:"primaryKey;type:string;default:uuid_generate_v4()" json:"id,omitempty"`
	OwnerID     string `json:"owner_id"`
	Slug        string `gorm:"unique" json:"slug,omitempty"`
	URL         string `json:"url,omitempty"`
	Visits      int    `json:"visits"`
	DateCreated int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	Lease       int64  `json:"lease,omitempty"`
}

// Visit structure containing time of visitation
type Visit struct {
	CreatedAt time.Time `json:"created_at"`
	LinkID    string    `json:"link_id"`
}

type VisitsPerDay struct {
	Count int    `json:"count"`
	Date  string `json:"date"`
}

// JWTClaims claims of the jwt
type JWTClaims struct {
	UserID string
	jwt.StandardClaims
}

// EmailVerification jwt claims to be used to verify an email
type EmailVerification struct {
	Email string
	jwt.StandardClaims
}

// PasswordResetToken contains whole user field and standard claims
type PasswordResetToken struct {
	User
	jwt.StandardClaims
}
