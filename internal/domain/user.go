package domain

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Age      uint8  `json:"age"`
}
