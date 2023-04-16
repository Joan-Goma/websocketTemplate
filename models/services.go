package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

func NewServices(connectionInfo string, logMode bool) error {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return err
	}

	db.LogMode(logMode)
	DBCONNECTION = db
	NewUserService(db)
	return nil
}

var DBCONNECTION *gorm.DB

type Services struct {
	User UserService
	db   *gorm.DB
}

func (s *Services) Close() error {
	return s.db.Close()
}

func DestructiveReset() error {
	if err := DBCONNECTION.DropTableIfExists(&UserMessage{}, &pwReset{}, &User{}).Error; err != nil {
		return err
	}
	return AutoMigrate()
}

func AutoMigrate() error {

	if err := DBCONNECTION.AutoMigrate(&User{}, &UserMessage{}, &pwReset{}).Error; err != nil {
		return err
	}
	return nil
}

type NeftModel struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
