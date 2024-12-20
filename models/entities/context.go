package entities

type Context struct {
	Base
	Notes      []Note     `json:"notes"`
	Documents  []Document `json:"documents"`
	Prompts    []Prompt   `json:"prompts"`
	UserID     string     `gorm:"type:uuid" json:"userId"`
	LanguageID string     `gorm:"type:uuid" json:"languageId"`
	ExternalID string     `json:"externalId"`
}
