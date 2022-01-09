package models

type User struct {
	Id       uint   `json:"id" ,gorm:"autoIncrement"`
	Username string `json:"username"`
	Email    string `json:"email" ,gorm:"unique"`
	Password []byte `json:"-"`
}

type Token struct {
	Email   string `json:"email"`
	TokenId string `json:"token_id"`
}
