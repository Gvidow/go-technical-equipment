package ds

type User struct {
	ID       int `gorm:"primary_key"`
	Username string
	Role     string
	Email    string
	Password string
	role     Role
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetRole() Role {
	if u.role != 0 {
		return u.role
	}
	u.role = ParseRole(u.Role)
	return u.role
}

func (u *User) SetRole(r Role) {
	u.role = r
	u.Role = r.String()
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
