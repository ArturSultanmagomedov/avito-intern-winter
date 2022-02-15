package service

import (
	"for_avito_tech_with_gin/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type User interface {
	AddFunds(userId int, sum float32) error
	WriteOffFunds(userId int, sum float32) error
	FundsTransfer(senderId int, receiverId int, sum float32) error
	GetBalance(userId int) (float32, error)
}

type Service struct {
	User
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		User: NewUserService(r),
	}
}
