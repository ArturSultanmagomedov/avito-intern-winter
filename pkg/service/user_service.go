package service

import (
	"for_avito_tech_with_gin/pkg/repository"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

// TODO: объединить AddFunds и WriteOffFunds

func (r *UserService) AddFunds(userId int, sum float32) error {
	if sum <= 0 {
		return &NegativeSum{}
	}

	ex, err := r.repo.IsUserExist(userId)
	if err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}
	if !ex {
		err := r.repo.CreateUser(userId, sum)
		if err != nil {
			logrus.Error(err)
			return &InternalServerError{}
		}
	} else {
		_, err := r.repo.UpdateBalance(userId, sum)
		if err != nil {
			logrus.Error(err)
			return &InternalServerError{}
		}
	}

	return nil
}

func (r *UserService) WriteOffFunds(userId int, sum float32) error {
	if sum <= 0 {
		return &NegativeSum{}
	}

	ex, err := r.repo.IsUserExist(userId)
	if err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}
	if !ex {
		return &UserNotFound{Id: userId}
	}

	user, err := r.repo.GetUser(userId)
	if err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}
	if user.Balance < sum {
		return &InsufficientFunds{Id: userId}
	}

	if _, err := r.repo.UpdateBalance(userId, -sum); err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}

	return nil
}

func (r *UserService) FundsTransfer(senderId int, receiverId int, sum float32) error {
	if sum <= 0 {
		return &NegativeSum{}
	}
	if senderId == receiverId {
		return &SameId{}
	}

	// Проверить существует ли отправляющий юзер (если не существует - вернуть ошибку)
	ex, err := r.repo.IsUserExist(senderId)
	if err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}
	if !ex {
		return &UserNotFound{Id: senderId}
	}

	// Проверить достаточно ли средств у отправляющего юзера (если нет - вернуть ошибку)
	user, err := r.repo.GetUser(senderId)
	if err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}
	if user.Balance < sum {
		return &InsufficientFunds{Id: senderId}
	}

	// Проверить существует ли получающий юзер (если не существует - создать)
	ex, err = r.repo.IsUserExist(receiverId)
	if err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}
	if !ex {
		err := r.repo.CreateUser(receiverId, 0)
		if err != nil {
			logrus.Error(err)
			return &InternalServerError{}
		}
	}

	err = r.repo.CreateFundsTransaction(senderId, receiverId, sum)
	if err != nil {
		logrus.Error(err)
		return &InternalServerError{}
	}

	return nil
}

func (r *UserService) GetBalance(userId int) (float32, error) {
	ex, err := r.repo.IsUserExist(userId)
	if err != nil {
		logrus.Error(err)
		return 0, &InternalServerError{}
	}
	if !ex {
		return 0, &UserNotFound{Id: userId}
	}

	user, err := r.repo.GetUser(userId)
	if err != nil {
		logrus.Error(err)
		return 0, &InternalServerError{}
	}

	return user.Balance, nil
}
