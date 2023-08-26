package todo

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ErrInvalidCredentials struct{}

func (e *ErrInvalidCredentials) Error() string {
	return "Invalid credentials"
}
