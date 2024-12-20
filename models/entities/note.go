package entities

type Note struct {
	Base
	Header     string     `json:"header"`
	Payload    string     `json:"payload"`
	UserID     string     `gorm:"type:uuid" json:"userId"`
	LanguageID string     `gorm:"type:uuid" json:"languageId"`
	Documents  []Document `json:"documents"`
	ContextID  string     `gorm:"type:uuid" json:"contextId"`
}
