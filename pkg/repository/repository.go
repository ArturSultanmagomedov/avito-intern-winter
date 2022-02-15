package repository

import (
	"database/sql"
	"for_avito_tech_with_gin/pkg/model"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type User interface {
	CreateUser(userId int, balance float32) error
	GetUser(userId int) (*model.User, error)
	IsUserExist(userId int) (bool, error)
	UpdateBalance(userId int, sum float32) (*model.User, error)
	CreateFundsTransaction(senderId int, receiverId int, sum float32) error
}

type Repository struct {
	User
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
