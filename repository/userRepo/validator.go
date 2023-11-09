package userRepo

import (
	validator "github.com/go-playground/validator/v10"
)

var user User

var createRules = map[string]string{
	"Login":       "alphanumeric",
	"Password":    "min=5",
	"Name":        "omitempty,alpha",
	"Surname":     "omitempty,alpha",
	"Email":       "email",
	"PhoneNumber": "omitempty,alphanumeric",
	"Age":         "omitempty,min=1,max=120,numeric",
}

var updateRules = map[string]string{
	"Id":          "required,alphanumeric",
	"Login":       "omitempty,alphanumeric",
	"Password":    "omitempty, min=5",
	"Name":        "omitempty,alpha",
	"Surname":     "omitempty,alpha",
	"Email":       "omitempty,email",
	"PhoneNumber": "omitempty,alphanumeric",
	"Age":         "omitempty,min=1,max=120,numeric",
}

var searchRules = map[string]string{
	"Name":    "omitempty,alphanumeric",
	"Surname": "omitempty,alphanumeric",
}

type userRepoValidator struct {
	validate *validator.Validate
	next     UserRepo
}

func NewUserRepoValidator(validate *validator.Validate, ur UserRepo) *userRepoValidator {
	return &userRepoValidator{
		validate: validate,
		next:     ur,
	}
}

func (urv *userRepoValidator) ValidateStruct(entity interface{}, rules map[string]string) error {

	urv.validate.RegisterStructValidationMapRules(rules, entity)
	err := urv.validate.Struct(entity)
	if err != nil {
		return err
	}
	return nil
}

func (urv *userRepoValidator) Create(up *UserPayload) (out *User, err error) {
	
	if err = urv.ValidateStruct(up, createRules); err != nil {
		return nil, err
	}
	out, err = urv.next.Create(up)
	return
}

func (urv *userRepoValidator) Read(id string) (out *User, err error) {
	if err = urv.validate.Var(id, "required,alphanumeric"); err != nil {
		return nil, err
	}
	out, err = urv.next.Read(id)
	return
}

func (urv *userRepoValidator) Search(us *UserSearch) (out []User, err error) {
	if err = urv.ValidateStruct(us, searchRules); err != nil {
		return nil, err
	}
	out, err = urv.next.Search(us)
	return
}

func (urv *userRepoValidator) Update(up *User) (out *User, err error) {
	if err = urv.ValidateStruct(up, updateRules); err != nil {
		return nil, err
	}
	out, err = urv.next.Update(up)
	return
}

func (urv *userRepoValidator) Delete(id string) (err error) {
	if err = urv.validate.Var(id, "required,alphanumeric"); err != nil {
		return err
	}
	err = urv.next.Delete(id)
	return
}
