package handler

import (
	avito_tech "for_avito_tech_with_gin/pkg"
	"for_avito_tech_with_gin/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// @Summary Add Funds
// @Description add funds (sum) for user (id)
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "input"
// @Success 200 {integer} integer
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /add_funds [post]
func (h *Handler) addFundsHandler(ctx *gin.Context) {
	s := &struct {
		UserId int     `json:"id" binding:"required"`
		Sum    float32 `json:"sum" binding:"required"`
	}{}
	if err := ctx.BindJSON(s); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, "invalid body.")
		return
	}

	if err := h.services.AddFunds(s.UserId, s.Sum); err != nil {
		responseError, ok := err.(service.ResponseError)
		if !ok {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		newErrorResponse(ctx, responseError.StatusCode(), responseError.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Write Off Funds
// @Description writes off funds (sum) for user (id)
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "input"
// @Success 200 {integer} integer
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 412 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /write_off_funds [post]
func (h *Handler) writeOffFundsHandler(ctx *gin.Context) {
	s := &struct {
		UserId int     `json:"id" binding:"required"`
		Sum    float32 `json:"sum" binding:"required"`
	}{}
	if err := ctx.BindJSON(s); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, "invalid body.")
		return
	}

	if err := h.services.WriteOffFunds(s.UserId, s.Sum); err != nil {
		responseError, ok := err.(service.ResponseError)
		if !ok {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		newErrorResponse(ctx, responseError.StatusCode(), responseError.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Funds Transfer
// @Description transfer funds (sum) from user (sender_id) to user (receiver_id)
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "input"
// @Param currency path string false "balance will convert from RUB to currency"
// @Success 200 {integer} integer
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 412 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /funds_transfer [post]
func (h *Handler) fundsTransferHandler(ctx *gin.Context) {
	s := &struct {
		SenderId   int     `json:"sender_id" binding:"required"`
		ReceiverId int     `json:"receiver_id" binding:"required"`
		Sum        float32 `json:"sum" binding:"required"`
	}{}
	if err := ctx.BindJSON(s); err != nil {
		logrus.Error(err)
		newErrorResponse(ctx, http.StatusBadRequest, "invalid body.")
		return
	}

	if err := h.services.FundsTransfer(s.SenderId, s.ReceiverId, s.Sum); err != nil {
		responseError, ok := err.(service.ResponseError)
		if !ok {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		newErrorResponse(ctx, responseError.StatusCode(), responseError.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary Get Balance
// @Description get user balance for user (id)
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /get_balance [get]
func (h *Handler) getBalanceHandler(calculator avito_tech.CurrencyCalculator) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		s := &struct {
			UserId int `json:"id" binding:"required"`
		}{}
		if err := ctx.BindJSON(s); err != nil {
			logrus.Error(err)
			newErrorResponse(ctx, http.StatusBadRequest, "invalid body.")
			return
		}

		balance, err := h.services.GetBalance(s.UserId)
		if err != nil {
			responseError, ok := err.(service.ResponseError)
			if !ok {
				newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
				return
			}
			newErrorResponse(ctx, responseError.StatusCode(), responseError.Error())
			return
		}

		currency := ctx.Query("currency")
		if currency == "" {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"balance": balance,
			})
		} else {
			if tt, err := calculator.ConvertRubTo(currency, balance); err != nil {
				responseError, ok := err.(service.ResponseError)
				if !ok {
					newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
					return
				}
				newErrorResponse(ctx, responseError.StatusCode(), responseError.Error())
				return
			} else {
				ctx.JSON(http.StatusOK, map[string]interface{}{
					"balance": tt,
				})
			}
		}
	}
}

// TODO: довести до ума документацию
