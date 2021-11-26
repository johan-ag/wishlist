package users

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetUser(id uint64) (User, error)
	CreateUser(user User) (User, error)
	GetUserByEmail(email string) (User, error)
}

func NewUsersRepository(db *gorm.DB) Repository {
	return &repository{
		db,
	}
}

type repository struct {
	db *gorm.DB
}

func (r *repository) GetUser(id uint64) (User, error) {
	var user User

	status := r.db.First(&user, id)
	if status.Error != nil {
		return user, errors.New("User not found")
	}

	return user, nil
}
func (r *repository) GetUserByEmail(email string) (User, error) {
	var user User
	status := r.db.Where("email = ? ", email).First(&user)
	if status.Error != nil {
		return user, errors.New("User not found")
	}

	return user, nil
}

func (r *repository) CreateUser(user User) (User, error) {
	_, err := r.GetUser(user.ID)

	if err == nil {
		return user, errors.New("This user already existed")
	}

	r.db.Save(&user)

	user, err = r.GetUser(user.ID)
	if err != nil {
		return user, err
	}
	return user, nil
}
