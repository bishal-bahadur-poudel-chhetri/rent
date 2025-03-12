package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	// Extract user ID from the token
	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Failed to parse form", err.Error()))
		return
	}

	// Parse and validate vehicle_id
	vehicleIDStr := c.PostForm("vehicle_id")
	if vehicleIDStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Vehicle ID is required", nil))
		return
	}

	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid vehicle ID", err.Error()))
		return
	}

	// Parse and validate other fields
	totalAmountStr := c.PostForm("total_amount")
	if totalAmountStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Total amount is required", nil))
		return
	}

	totalAmount, err := strconv.ParseFloat(totalAmountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid total amount", err.Error()))
		return
	}

	chargePerDayStr := c.PostForm("charge_per_day")
	if chargePerDayStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Charge per day is required", nil))
		return
	}

	chargePerDay, err := strconv.ParseFloat(chargePerDayStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid charge per day", err.Error()))
		return
	}

	numberOfDaysStr := c.PostForm("number_of_days")
	if numberOfDaysStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Number of days is required", nil))
		return
	}

	numberOfDays, err := strconv.Atoi(numberOfDaysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid number of days", err.Error()))
		return
	}

	// Parse payment-related fields
	amountPaidStr := c.PostForm("amount_paid")
	if amountPaidStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Amount paid is required", nil))
		return
	}

	amountPaid, err := strconv.ParseFloat(amountPaidStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid amount paid", err.Error()))
		return
	}

	paymentDateStr := c.PostForm("payment_date")
	if paymentDateStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Payment date is required", nil))
		return
	}

	paymentDate, err := time.Parse(time.RFC3339, paymentDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payment date", err.Error()))
		return
	}

	verifiedByAdmin := parseBool(c.PostForm("verified_by_admin"))
	paymentType := c.PostForm("payment_type")
	if paymentType == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Payment type is required", nil))
		return
	}

	paymentStatus := c.PostForm("payment_status")
	if paymentStatus == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Payment status is required", nil))
		return
	}

	// Validate required fields
	if c.PostForm("customer_name") == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Customer name is required", nil))
		return
	}
	if c.PostForm("customer_destination") == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Customer Destination is required", nil))
		return
	}
	// Parse dates with error handling
	bookingDate, err := time.Parse(time.RFC3339, c.PostForm("booking_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid booking date", err.Error()))
		return
	}

	dateOfDelivery, err := time.Parse(time.RFC3339, c.PostForm("date_of_delivery"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid delivery date", err.Error()))
		return
	}

	returnDate, err := time.Parse(time.RFC3339, c.PostForm("return_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid return date", err.Error()))
		return
	}

	// Parse boolean values
	isDamaged := parseBool(c.PostForm("is_damaged"))
	isWashed := parseBool(c.PostForm("is_washed"))
	isDelayed := parseBool(c.PostForm("is_delayed"))

	// Set default values for optional fields
	remark := c.PostForm("remark") // Optional field
	status := c.PostForm("status")
	if status == "" {
		status = "pending" // Default status
	}

	// Parse sales charges (JSON array)
	salesChargesJSON := c.PostForm("sales_charges")
	var salesCharges []models.SalesCharge
	if salesChargesJSON != "" {
		if err := json.Unmarshal([]byte(salesChargesJSON), &salesCharges); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sales_charges format", err.Error()))
			return
		}
	}

	// Parse vehicle usage (JSON array)
	vehicleUsageJSON := c.PostForm("vehicle_usage")
	var vehicleUsage []models.VehicleUsage
	if vehicleUsageJSON != "" {
		if err := json.Unmarshal([]byte(vehicleUsageJSON), &vehicleUsage); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid vehicle_usage format", err.Error()))
			return
		}
	}

	// Handle file uploads for images
	var salesImages []models.SalesImage
	imageFiles := c.Request.MultipartForm.File["sales_images"]
	for _, fileHeader := range imageFiles {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to open image file", err.Error()))
			return
		}
		defer file.Close()

		filePath := fmt.Sprintf("uploads/images/%s", fileHeader.Filename)
		if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to save image file", err.Error()))
			return
		}

		salesImages = append(salesImages, models.SalesImage{
			ImageURL: filePath,
		})
	}

	// Handle file uploads for videos
	var salesVideos []models.SalesVideo
	videoFiles := c.Request.MultipartForm.File["sales_videos"]
	for _, fileHeader := range videoFiles {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to open video file", err.Error()))
			return
		}
		defer file.Close()

		filePath := fmt.Sprintf("uploads/videos/%s", fileHeader.Filename)
		if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to save video file", err.Error()))
			return
		}

		salesVideos = append(salesVideos, models.SalesVideo{
			VideoURL: filePath,
		})
	}

	// Create the Sale struct
	sale := models.Sale{
		VehicleID:      vehicleID,
		UserID:         userID,
		CustomerName:   c.PostForm("customer_name"),
		Destination:    c.PostForm("customer_destination"),
		TotalAmount:    totalAmount,
		ChargePerDay:   chargePerDay,
		BookingDate:    bookingDate,
		DateOfDelivery: dateOfDelivery,
		ReturnDate:     returnDate,
		IsDamaged:      isDamaged,
		IsWashed:       isWashed,
		IsDelayed:      isDelayed,
		NumberOfDays:   numberOfDays,
		Remark:         remark,
		Status:         status,
		SalesCharges:   salesCharges,
		SalesImages:    salesImages,
		SalesVideos:    salesVideos,
		VehicleUsage:   vehicleUsage,

		// Payment-related fields
		AmountPaid:      amountPaid,
		PaymentDate:     paymentDate,
		VerifiedByAdmin: verifiedByAdmin,
		PaymentType:     paymentType,
		PaymentStatus:   paymentStatus,
	}

	// Create the sale in the database
	saleID, err := h.saleService.CreateSale(sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create sale", err.Error()))
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Sale created successfully", gin.H{"sale_id": saleID}))
}

func (h *SaleHandler) GetSaleByID(c *gin.Context) {
	// Convert sale ID from string to int
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}

	// Fetch sale
	sale, err := h.saleService.GetSaleByID(saleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch sale", err.Error()))
		return
	}

	// Check if sale exists
	if sale == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Sale not found", nil))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sale fetched successfully", sale))
}
