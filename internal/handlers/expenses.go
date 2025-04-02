package handlers

import (
	"net/http"
	"strconv"

	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils" // Import the utils package

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	service services.ExpenseService
}

func NewExpenseHandler(service services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	if err := h.service.CreateExpense(&expense); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Expense created successfully", expense))
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid expense ID", nil))
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	expense.ExpenseID = id
	if err := h.service.UpdateExpense(id, &expense); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Expense updated successfully", expense))
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid expense ID", nil))
		return
	}

	if err := h.service.DeleteExpense(id); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Expense deleted successfully", nil))
	// Note: Changed from StatusNoContent to JSON response for consistency with StandardResponse
}

func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid expense ID", nil))
		return
	}

	expense, err := h.service.GetExpense(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Expense not found", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Expense retrieved successfully", expense))
}

func (h *ExpenseHandler) GetAllExpenses(c *gin.Context) {
	var filter models.ExpenseFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid query parameters", err.Error()))
		return
	}

	expenses, err := h.service.GetAllExpenses(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Expenses retrieved successfully", expenses))
}
