package user

type UpdateUserRequest struct {
	ID   string  `json:"id"`
	Name *string `json:"name"`
}

type Role uint

const (
	Student Role = iota
	Teacher
	StudentTeacher
)
