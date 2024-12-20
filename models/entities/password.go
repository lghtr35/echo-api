package entities

type Password struct {
	Base
	UserID string `gorm:"type:uuid" json:"userId"`
	Value  string `json:"password"`
}
