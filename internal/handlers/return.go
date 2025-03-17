package handlers

import (
	"fmt"
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReturnHandler struct {
	returnService *services.ReturnService
	jwtSecret     string
}

func NewReturnHandler(returnService *services.ReturnService, jwtSecret string) *ReturnHandler {
	return &ReturnHandler{returnService: returnService, jwtSecret: jwtSecret}
}

func (h *ReturnHandler) CreateReturn(c *gin.Context) {
	// Extract userID from the JWT token
	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Debugging: Log that the handler is called
	fmt.Println("CreateReturn handler called")

	// Extract sale_id from the URL path parameter
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}

	// Debugging: Log the sale ID and user ID
	fmt.Println("Sale ID:", saleID)
	fmt.Println("User ID:", userID)

	// Parse the request body into the ReturnRequest struct
	var returnRequest models.ReturnRequest
	if err := c.ShouldBindJSON(&returnRequest); err != nil {
		// Debugging: Log the error
		fmt.Println("Error parsing JSON body:", err)

		// Return a 400 Bad Request response with the error details
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	// Debugging: Log the parsed request body
	fmt.Println("Parsed request body:", returnRequest)

	// Call the service to create the return record, passing the userID
	if err := h.returnService.CreateReturn(saleID, userID, returnRequest); err != nil {
		// Debugging: Log the error
		fmt.Println("Error creating return record:", err)

		// Return a 500 Internal Server Error response with the error details
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create return record", err.Error()))
		return
	}

	// Return a 201 Created response with a success message
	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Return record created successfully", nil))
}
