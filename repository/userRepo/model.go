package userRepo

import (
	"smartHomeKit/common"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserSearch struct {
	Limit       int                          `json:"limit"`
	Offset      int                          `json:"offset"`
	Name        string                       `json:"name"`
	Surname     string                       `json:"surname"`
	Email       string                       `json:"email"`
	PhoneNumber string                       `json:"phoneNumber"`
	CreatedAt   common.GraterLess[time.Time] `json:"createdAt"`
	Age         common.GraterLess[int]       `json:"age"`
}

type UserPayload struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Age         int    `json:"age"`
}

type User struct {
	Id           string `json:"id"`
	*UserPayload `gorm:"embedded"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (u *User) HashPassword() (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return nil, err
	}
	u.Password = string(hashed)
	return u, nil
}

func (u *User) CheckPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	if string(hashed) != u.Password {
		return ErrWrongPassword
	}
	return nil
}
