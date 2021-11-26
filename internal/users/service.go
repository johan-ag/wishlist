package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(user User) (User, error)
	GetUser(id uint64) (User, error)
	GetUserByEmail(email string) (User, error)
}

func NewUsersService(repository Repository) Service {
	return &service{
		repository,
	}
}

type service struct {
	repository Repository
}

func (s *service) CreateUser(user User) (User, error) {
	if user.ID != 0 || user.Name == "" || user.Email == "" || user.Password == "" {
		return user, errors.New("Invalid user")
	}

	hashedPass, err := hashAndSalt(user.Password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPass

	user, err = s.repository.CreateUser(user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) GetUser(id uint64) (User, error) {
	user, err := s.repository.GetUser(id)
	if err != nil {
		return user, err
	}
	user.Password = ""
	return user, nil
}

func (s *service) GetUserByEmail(email string) (User, error) {
	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
