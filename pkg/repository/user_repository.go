package repository

import (
	"database/sql"
	"for_avito_tech_with_gin/pkg/model"
	"github.com/pkg/errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(userId int, balance float32) error {
	_, err := r.db.Query("insert into users (user_id, balance) values ($1, $2);", userId, balance)
	if err != nil {
		return errors.Wrapf(err, "filed to create user %d", userId)
	}
	return nil
}

func (r *UserRepository) GetUser(userId int) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow("select id, user_id, balance from users where user_id = $1", userId).Scan(user.GetFields()...)
	if err != nil {
		return nil, errors.Wrapf(err, "filed to get user %d", userId)
	}

	return &user, err
}

func (r *UserRepository) IsUserExist(userId int) (bool, error) {
	var c int
	err := r.db.QueryRow("select count(1) from users where user_id = $1;", userId).Scan(&c)
	if err != nil {
		return false, errors.Wrapf(err, "filed to check is user %d exist", userId)
	}

	return c > 0, nil
}

func (r *UserRepository) UpdateBalance(userId int, sum float32) (*model.User, error) {
	var user model.User

	tx, err := r.db.Begin()
	if err != nil {
		return nil, errors.Wrapf(err, "filed to begin transaction and update balance for user %d", userId)
	}
	defer tx.Rollback()

	err = tx.QueryRow("select id, user_id, balance from users where user_id = $1;", userId).Scan(user.GetFields()...)
	if err != nil {
		return nil, errors.Wrapf(err, "filed to get user and update balance for user %d", userId)
	}

	err = tx.QueryRow("update users set balance = $1 where user_id = $2;", user.Balance+sum, userId).Scan(user.GetFields()...)
	if err != nil {
		return nil, errors.Wrapf(err, "filed update balance for user %d", userId)
	}

	return &user, tx.Commit()
}

func (r *UserRepository) CreateFundsTransaction(senderId int, receiverId int, sum float32) error {

	var sender model.User
	var receiver model.User

	tx, err := r.db.Begin()
	if err != nil {
		return errors.Wrapf(err, "filed to begin transaction and create transaction between %d and %d users", senderId, receiverId)
	}
	defer tx.Rollback()

	err = tx.QueryRow("select id, user_id, balance from users where user_id = $1;", senderId).Scan(sender.GetFields()...)
	if err != nil {
		return errors.Wrapf(err, "filed to get user %d and create transaction between %d and %d users", senderId, senderId, receiverId)
	}

	err = tx.QueryRow("select id, user_id, balance from users where user_id = $1;", receiverId).Scan(receiver.GetFields()...)
	if err != nil {
		return errors.Wrapf(err, "filed to get user %d and create transaction between %d and %d users", receiverId, senderId, receiverId)
	}

	_, err = tx.Exec("update users set balance = $1 where user_id = $2;", sender.Balance-sum, senderId)
	if err != nil {
		return errors.Wrapf(err, "filed to update user %d and create transaction between %d and %d users", senderId, senderId, receiverId)
	}

	_, err = tx.Exec("update users set balance = $1 where user_id = $2;", receiver.Balance+sum, receiverId)
	if err != nil {
		return errors.Wrapf(err, "filed to update user %d and create transaction between %d and %d users", senderId, senderId, receiverId)
	}

	return tx.Commit()
}
