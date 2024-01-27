package entities

import (
	"math"
	"time"
)

type User struct {
	ID        string
	FullName  string
	Email     string
	Password  string
	Image     *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeleteAt  *time.Time
}

func (u User) validateFullName(fullName string) error {
	err := emptyStringValidator(fullName, "nome ou sobrenome")
	if err != nil {
		return err
	}
	err = invalidFieldLenghValidator(fullName, "nome ou sobrenome", 3, math.MaxUint8)
	if err != nil {
		return err
	}
	return nil
}

func (u User) validateEmail(email string) error {
	err := emptyStringValidator(email, "email")
	if err != nil {
		return err
	}
	err = emailPatternValidator(email)
	if err != nil {
		return err
	}
	return nil
}

func (u User) validatePass(pass string) error {
	err := emptyStringValidator(pass, "password")
	if err != nil {
		return err
	}
	err = invalidFieldLenghValidator(pass, "password", 8, 32)
	if err != nil {
		return err
	}
	err = passwordPatternValidator(pass)
	if err != nil {
		return err
	}
	return nil
}

func (u User) validator() error {
	err := u.validateFullName(u.FullName)
	if err != nil {
		return err
	}
	err = u.validateEmail(u.Email)
	if err != nil {
		return err
	}
	err = u.validatePass(u.Password)
	if err != nil {
		return err
	}
	return nil
}

func NewUser(user User) (*User, error) {
	err := user.validator()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
