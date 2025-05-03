package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"renting/internal/models"
	"renting/internal/repositories"
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

// CheckAccountingPermission middleware to verify if user has accounting permission
func CheckAccountingPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get userID from the context (set by JWT middleware)
		userID, exists := c.Get("userID")
		if !exists {
			log.Printf("[ERROR] User ID not found in context")
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
			c.Abort()
			return
		}

		log.Printf("[DEBUG] CheckAccountingPermission middleware - User ID: %d", userID.(int))

		// Get userRepo from the app context or create a new one
		db, exists := c.Get("sqlDB")
		if !exists {
			log.Printf("[ERROR] Database connection not found in context")
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Database connection not available", nil))
			c.Abort()
			return
		}

		// Get the user from the database
		userRepo := repositories.NewUserRepository(db.(*sql.DB))
		user, err := userRepo.GetUserByID(c.Request.Context(), userID.(int))
		if err != nil {
			log.Printf("[ERROR] Failed to get user by ID: %v", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Error fetching user data", nil))
			c.Abort()
			return
		}

		if user == nil {
			log.Printf("[ERROR] User not found with ID: %d", userID.(int))
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not found", nil))
			c.Abort()
			return
		}

		log.Printf("[DEBUG] User permissions - IsAdmin: %v, HasAccounting: %v", user.IsAdmin, user.HasAccounting)

		// Check permissions
		if !user.HasAccounting && !user.IsAdmin {
			log.Printf("[INFO] User %d denied access - no accounting permission", userID.(int))
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "User does not have accounting permission", nil))
			c.Abort()
			return
		}

		log.Printf("[INFO] User %d granted accounting access - IsAdmin: %v, HasAccounting: %v",
			userID.(int), user.IsAdmin, user.HasAccounting)
		c.Next()
	}
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
