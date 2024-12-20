package entities

type Document struct {
	Base
	Name            string  `json:"name"`
	Location        string  `json:"location"`
	Extension       string  `json:"extension"`
	NoteID          *string `gorm:"type:uuid" json:"noteId"`
	UserID          string  `gorm:"type:uuid" json:"userId"`
	ContextID       string  `gorm:"type:uuid" json:"contextId"`
	IsReadableByAll bool    `json:"isReadableByAll"`
}
