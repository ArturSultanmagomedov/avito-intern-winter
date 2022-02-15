package service

import (
	"fmt"
	"net/http"
)

type ResponseError interface {
	error
	StatusCode() int
}

// NegativeSum - для ситуаций, когда в запросе пришла сумма <=0
type NegativeSum struct{}

func (r *NegativeSum) Error() string {
	return "sum can't be negative or 0."
}

func (r *NegativeSum) StatusCode() int {
	return http.StatusBadRequest
}

// SameId - для ситуаций, когда в юзер пытается отправить деньги самому себе
type SameId struct{}

func (r *SameId) Error() string {
	return "user cannot send money to himself."
}

func (r *SameId) StatusCode() int {
	return http.StatusBadRequest
}

// UserNotFound - для ситуаций, когда в базе не нашлось нужного юзера
type UserNotFound struct {
	Id int
}

func (r *UserNotFound) Error() string {
	return fmt.Sprintf("user %d does not exist.", r.Id)
}

func (r *UserNotFound) StatusCode() int {
	return http.StatusNotFound
}

// InsufficientFunds - для ситуаций, когда у юзера не хватает денег для перевода
type InsufficientFunds struct {
	Id int
}

func (r *InsufficientFunds) Error() string {
	return fmt.Sprintf("user %d has insufficient funds.", r.Id)
}

func (r *InsufficientFunds) StatusCode() int {
	return http.StatusPreconditionFailed
}

// InternalServerError - для ситуаций, когда черт его знает че там за проблема с бд
type InternalServerError struct{}

func (r *InternalServerError) Error() string {
	return "internal server error."
}

func (r *InternalServerError) StatusCode() int {
	return http.StatusInternalServerError
}

// WrongParam - для ситуаций, когда в запросе указан кривой параметр, например неподдерживаемая валюта
type WrongParam struct {
	Param string
}

func (r *WrongParam) Error() string {
	return fmt.Sprintf("wrong %s param.", r.Param)
}

func (r *WrongParam) StatusCode() int {
	return http.StatusPreconditionFailed
}
