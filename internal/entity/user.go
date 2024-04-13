package entity

import "time"

const (
	RoleAdmin Role = iota
	RoleUser
)

type Role int8

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	default:
		return "unknown"
	}
}

func GetRoleFromString(role string) Role {
	switch role {
	case "admin":
		return RoleAdmin
	case "user":
		return RoleUser
	default:
		return RoleUser
	}
}

type UserCredentials struct {
	Phone    string  `json:"phone" db:"phone"`
	Email    *string `json:"email" db:"email"`
	Password string  `json:"password" db:"password"`
}

type User struct {
	Id string `json:"id"`
	UserCredentials
	LastName        *string    `json:"lastName" db:"last_name"`
	FirstName       *string    `json:"firstName" db:"first_name"`
	MiddleName      *string    `json:"middleName" db:"middle_name"`
	Role            string     `json:"role" db:"role"`
	IsPhoneVerified bool       `json:"isPhoneVerified" db:"is_phone_verified"`
	IsEmailVerified bool       `json:"isEmailVerified" db:"is_email_verified"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt       *time.Time `json:"updatedAt" db:"updated_at"`
}

type UserClaims struct {
	Id   string `json:"id"`
	Role string `json:"role"`
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
