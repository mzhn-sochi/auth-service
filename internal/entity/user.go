package entity

import "time"

const (
	RoleAdmin Role = iota
	RoleUser
)

type Role int8

type UserCredentials struct {
	Phone    string  `json:"phone"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

type User struct {
	Id string `json:"id"`
	UserCredentials
	LastName   *string    `json:"last_name"`
	FirstName  *string    `json:"first_name"`
	MiddleName *string    `json:"middle_name"`
	Role       Role       `json:"role"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type UserClaims struct {
	Id   string `json:"id"`
	Role Role   `json:"role"`
}

type Tokens struct {
	Refresh string `json:"refresh_token"`
	Access  string `json:"access_token"`
}

func (u *User) GetClaims() *UserClaims {
	return &UserClaims{
		Id:   u.Id,
		Role: u.Role,
	}
}
