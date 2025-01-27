package handler

import (
	"context"
	"net/http"
	"strconv"

	"user-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

const (
	userURL          = "/user/:id"
	replenishmentURL = "/replenishment"
	paymentURL       = "/payment"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type userService interface {
	ReplenishmentBalance(ctx context.Context, id int, amount decimal.Decimal) error
	Payment(ctx context.Context, senderID int, recipientID int, amount decimal.Decimal) error
	RecentOperations(ctx context.Context, id int) ([]model.Operation, error)
}

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Handler struct {
	userService
	logger Logger
}

func NewHandler(userService userService, logger Logger) *Handler {
	return &Handler{userService: userService, logger: logger}
}

func (h *Handler) Register(router *gin.Engine) {
	router.POST(paymentURL, h.Payment)
	router.POST(replenishmentURL, h.Replenishment)
	router.GET(userURL, h.CheckUserTransactions)
}

func (h *Handler) Payment(ctx *gin.Context) {
	var transaction model.Payment

	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	if err := h.userService.Payment(ctx.Request.Context(), transaction.SenderID, transaction.RecipientID, transaction.Amount); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"result": "success"})

}

func (h *Handler) Replenishment(ctx *gin.Context) {
	var replenishment model.Replenishment

	if err := ctx.ShouldBindJSON(&replenishment); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	err := h.userService.ReplenishmentBalance(ctx.Request.Context(), replenishment.RecipientID, replenishment.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"result": "success"})

}

func (h *Handler) CheckUserTransactions(ctx *gin.Context) {
	IDStr := ctx.Param("id")
	id, err := strconv.Atoi(IDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	result, err := h.userService.RecentOperations(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
