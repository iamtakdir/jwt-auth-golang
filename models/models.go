package models

type User struct {
	Id       uint `gorm:"autoIncrement"`
	Username string
	Email    string `gorm:"unique"`
	Password []byte
}
