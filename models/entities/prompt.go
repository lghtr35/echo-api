package entities

type Prompt struct {
	Base
	Value     string
	ContextID string  `gorm:"type:uuid" json:"contextId"`
	EntityID  *string `gorm:"type:uuid" json:"entityId"`
}
