package handler

import (
	"context"
	"net/http"
	"strconv"

	"user-service/internal/model"

	_ "user-service/docs"

	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"

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

	// URL: /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

// Payment обрабатывает запрос на перевод
// @Summary Process a payment
// @Description Processes a payment from one user to another
// @Accept json
// @Produce json
// @Param transaction body model.Payment true "Transaction details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /payment [post]
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

// Replenishment обрабатывает запрос на пополнение баланса
// @Summary Replenish user balance
// @Description Replenishes the balance of a user
// @Accept json
// @Produce json
// @Param replenishment body model.Replenishment true "Replenishment details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /replenishment [post]
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

// CheckUser Transactions обрабатывает запрос на получение транзакций пользователя
// @Summary Get user transactions
// @Description Retrieves a list of recent operations for a user
// @Produce json
// @Param id path int true "User  ID"
// @Success 200 {array} model.Operation
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/{id} [get]
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
