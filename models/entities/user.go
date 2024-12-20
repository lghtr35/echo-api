package entities

type User struct {
	Base
	Name      string      `json:"name"`
	Email     string      `gorm:"unique" json:"email"`
	Notes     []Note      `json:"notes"`
	Languages []*Language `gorm:"many2many:user_languages;" json:"languages"`
	Documents []Document  `json:"documents"`
	Contexts  []Context   `json:"contexts"`
	Password  Password    `json:"password"`
	Role      Role        `json:"role"`
}

type Role uint

const (
	Admin Role = iota + 1
	Customer
)

func (r Role) ToString() string {
	switch r {
	case Admin:
		return "Admin"
	case Customer:
		return "Customer"
	default:
		return "Error"
	}
}
