package ds

type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (u *User) TableName() string {
	return "users"
}
