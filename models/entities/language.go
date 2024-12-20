package entities

type Language struct {
	Base
	Name       string    `json:"name"`
	Alpha2Code string    `json:"alpha2Code"`
	Alpha3Code string    `json:"alpha3Code"`
	Icon       string    `json:"icon"`
	Notes      []Note    `json:"notes"`
	Users      []*User   `gorm:"many2many:user_languages;" json:"users"`
	Contexts   []Context `json:"contexts"`
}
