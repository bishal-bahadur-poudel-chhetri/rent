package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"

	"strconv"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	service services.ExpenseService
}

func NewExpenseHandler(service services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateExpense(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) ReimburseExpense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	if err := h.service.ReimburseExpense(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "reimbursed"})
}

func (h *ExpenseHandler) AutoReimburse(c *gin.Context) {
	if err := h.service.AutoReimburse(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "auto-reimbursement completed"})
}

func (h *ExpenseHandler) GetPendingExpenses(c *gin.Context) {
	expenses, err := h.service.GetPendingExpenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}
