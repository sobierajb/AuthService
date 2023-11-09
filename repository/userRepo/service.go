package userRepo

import (
	"fmt"

	"gorm.io/gorm"
)

type UserRepo interface {
	Create(user *UserPayload) (*User, error)
	Read(id string) (*User, error)
	GetByLogin(login string) (*User, error)
	Search(findPayload *UserSearch) ([]User, error)
	Update(user *User) (*User, error)
	Delete(id string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (ur *userRepo) Create(up *UserPayload) (*User, error) {

	user := &User{
		UserPayload: up,
	}
	// hash  User password
	user.HashPassword()
	result := ur.db.Create(&user)
	return user, result.Error
}

func (ur *userRepo) Read(id string) (*User, error) {
	var user User
	result := ur.db.Select("*").Omit("Password").First(&user, id)
	return &user, result.Error
}

func (ur *userRepo) GetByLogin(login string) (*User, error) {
	var user User
	result := ur.db.Where("Login = ?", login).First(&user)
	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil

}

func (ur *userRepo) Search(uf *UserSearch) ([]User, error) {
	var users []User

	if uf.Limit > 0 {
		ur.db = ur.db.Limit(uf.Limit)
	}
	if uf.Offset > 0 {
		ur.db = ur.db.Offset(uf.Offset)
	}
	if uf.Name != "" {
		ur.db = ur.db.Where("Name LIKE ?", fmt.Sprintf("%%%s%%", uf.Name))
	}
	if uf.Surname != "" {
		ur.db = ur.db.Where("Surname Like ?", fmt.Sprintf("%%%s%%", uf.Surname))
	}
	if uf.Email != "" {
		ur.db = ur.db.Where("Email Like ?", fmt.Sprintf("%%%s%%", uf.Email))
	}

	if !uf.CreatedAt.GraterThen.IsZero() {
		ur.db = ur.db.Where("CreatedAt > ?", fmt.Sprintf("%v", uf.CreatedAt.GraterThen))
	}

	if !uf.CreatedAt.LessThen.IsZero() {
		ur.db = ur.db.Where("CreatedAt < ?", fmt.Sprintf("%v", uf.CreatedAt.LessThen))
	}
	ur.db = ur.db.Select("*").Omit("password").Find(&users)
	return users, ur.db.Error
}

func (ur *userRepo) Update(user *User) (*User, error) {
	if user.Password != "" {
		user.HashPassword()
	} else {
		ur.db = ur.db.Omit("Password")
	}
	ur.db = ur.db.Model(&User{}).Updates(user)
	return user, ur.db.Error
}

func (ur *userRepo) Delete(id string) error {
	result := ur.db.Delete(&User{}, id)
	return result.Error
}
