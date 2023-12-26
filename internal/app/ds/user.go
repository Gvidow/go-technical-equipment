package ds

type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func (u User) TableName() string {
	return "users"
}

type Role uint8

func (r Role) String() string {
	switch r {
	case RegularUser:
		return "user"
	case Moderator:
		return "moderator"
	default:
		return "guest"
	}
}

const (
	_ Role = iota
	Guest
	RegularUser
	Moderator
)

func (u User) GetRole() Role {
	return ParseRole(u.Role)
}

func ParseRole(role string) Role {
	switch role {
	case "user":
		return RegularUser
	case "moderator":
		return Moderator
	default:
		return Guest
	}
}
