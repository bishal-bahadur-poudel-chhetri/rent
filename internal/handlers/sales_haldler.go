package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ttacon/libphonenumber"
)

type SaleHandler struct {
	saleService *services.SaleService
	jwtSecret   string
}

func NewSaleHandler(saleService *services.SaleService, jwtSecret string) *SaleHandler {
	return &SaleHandler{
		saleService: saleService,
		jwtSecret:   jwtSecret,
	}
}

func parseBool(value string) bool {
	return strings.ToLower(value) == "true"
}

func (h *SaleHandler) CreateSale(c *gin.Context) {
	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	var saleRequest struct {
		VehicleID           int                   `json:"vehicle_id"`
		TotalAmount         float64               `json:"total_amount"`
		ChargePerDay        float64               `json:"charge_per_day"`
		NumberOfDays        int                   `json:"number_of_days"`
		CustomerPhone       string                `json:"customer_phone"`
		AmountPaid          float64               `json:"amount_paid"`
		PaymentStatus       string                `json:"payment_status"`
		CustomerName        string                `json:"customer_name"`
		CustomerDestination string                `json:"customer_destination"`
		DateOfDelivery      string                `json:"date_of_delivery"`
		ReturnDate          string                `json:"return_date"`
		Remark              string                `json:"remark"`
		Status              string                `json:"status"`
		SalesCharges        []models.SalesCharge  `json:"sales_charges"`
		VehicleUsage        []models.VehicleUsage `json:"vehicle_usage"`
		Payments            []models.Payment      `json:"payments"`
		SalesImages         []models.SalesImage   `json:"sales_images"`
		SalesVideos         []models.SalesVideo   `json:"sales_videos"`
	}

	if err := c.ShouldBindJSON(&saleRequest); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid JSON payload", err.Error()))
		return
	}

	dateOfDelivery, err := time.Parse(time.RFC3339, saleRequest.DateOfDelivery)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid delivery date", err.Error()))
		return
	}
	fmt.Println(saleRequest.DateOfDelivery)

	if !isValidNepalesePhoneNumber(saleRequest.CustomerPhone) {
		fmt.Print("hi")
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid phone number", "Phone number must be a valid Nepalese number (e.g., +9779841234567)"))
		return
	}

	returnDate, err := time.Parse(time.RFC3339, saleRequest.ReturnDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid return date", err.Error()))
		return
	}

	paymentDate := time.Now()
	for i := range saleRequest.Payments {
		saleRequest.Payments[i].PaymentDate = paymentDate
	}

	sale := models.Sale{
		VehicleID:      saleRequest.VehicleID,
		UserID:         userID,
		CustomerName:   saleRequest.CustomerName,
		Destination:    saleRequest.CustomerDestination,
		CustomerPhone:  saleRequest.CustomerPhone,
		TotalAmount:    saleRequest.TotalAmount,
		ChargePerDay:   saleRequest.ChargePerDay,
		DateOfDelivery: dateOfDelivery,
		ReturnDate:     returnDate,
		NumberOfDays:   saleRequest.NumberOfDays,
		Remark:         saleRequest.Remark,
		Status:         saleRequest.Status,
		SalesCharges:   saleRequest.SalesCharges,
		SalesImages:    saleRequest.SalesImages,
		SalesVideos:    saleRequest.SalesVideos,
		VehicleUsage:   saleRequest.VehicleUsage,
		Payments:       saleRequest.Payments,
	}

	saleID, err := h.saleService.CreateSale(sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create sale", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Sale created successfully", gin.H{"sale_id": saleID}))
}

func isValidNepalesePhoneNumber(phone string) bool {
	parsedNumber, err := libphonenumber.Parse(phone, "NP")
	if err != nil {
		return false
	}
	return libphonenumber.IsValidNumberForRegion(parsedNumber, "NP")
}

func (h *SaleHandler) GetSaleByID(c *gin.Context) {
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}

	includeParam := c.Query("include")
	var include []string
	if includeParam != "" {
		include = strings.Split(includeParam, ",")
	}

	sale, err := h.saleService.GetSaleByID(saleID, include)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch sale", err.Error()))
		return
	}

	if sale == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Sale not found", nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sale fetched successfully", sale))
}

func (h *SaleHandler) GetSales(c *gin.Context) {
	_, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	filters := make(map[string]string)
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if actualDelivery := c.Query("actual_date_of_delivery"); actualDelivery != "" {
		filters["actual_date_of_delivery"] = actualDelivery
	}
	if dateBefore := c.Query("date_of_delivery_before"); dateBefore != "" {
		filters["date_of_delivery_before"] = dateBefore
	}
	if customerName := c.Query("customer_name"); customerName != "" {
		filters["customer_name"] = customerName
	}
	if vehicleID := c.Query("vehicle_id"); vehicleID != "" {
		filters["vehicle_id"] = vehicleID
	}

	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	includeParam := c.Query("include")
	var include []string
	if includeParam != "" {
		include = strings.Split(includeParam, ",")
	}

	response, err := h.saleService.GetSales(filters, sort, limit, offset, include)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch sales", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sales fetched successfully", response))
}

func (h *SaleHandler) UpdateSaleByUserID(c *gin.Context) {
	body, _ := c.GetRawData()
	fmt.Println("Received body:", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err))
		return
	}

	saleIDStr := c.Param("saleID")
	saleID, err := strconv.Atoi(saleIDStr)
	if err != nil || saleID <= 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", err))
		return
	}

	var req models.UpdateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Binding error:", err)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	// Debug: Check if VehicleID is bound correctly
	fmt.Printf("req.VehicleID: %v\n", req.VehicleID)

	if err := validateUpdateSaleRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Pass the req directly to the service
	if err := h.saleService.UpdateSaleByUserID(saleID, userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sale updated successfully", nil))
}

func validateUpdateSaleRequest(req models.UpdateSaleRequest) error {
	// Only validate if fields are provided (non-nil)
	if req.Status != nil && *req.Status == "" {
		return fmt.Errorf("status cannot be empty if provided")
	}
	if req.CustomerName != nil && *req.CustomerName == "" {
		return fmt.Errorf("customer_name cannot be empty if provided")
	}
	if req.TotalAmount != nil && *req.TotalAmount == 0 {
		return fmt.Errorf("total_amount cannot be zero if provided")
	}
	return nil
}

