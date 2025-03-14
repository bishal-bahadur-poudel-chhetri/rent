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

	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Failed to parse form", err.Error()))
		return
	}

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

	amountPaidStr := c.PostForm("amount_paid")
	if amountPaidStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Amount paid is required", nil))
		return
	}

	paymentDateStr := c.PostForm("payment_date")
	if paymentDateStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Payment date is required", nil))
		return
	}

	paymentStatus := c.PostForm("payment_status")
	if paymentStatus == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Payment status is required", nil))
		return
	}

	if c.PostForm("customer_name") == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Customer name is required", nil))
		return
	}
	if c.PostForm("customer_destination") == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Customer Destination is required", nil))
		return
	}

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

	isDamaged := parseBool(c.PostForm("is_damaged"))
	isWashed := parseBool(c.PostForm("is_washed"))
	isDelayed := parseBool(c.PostForm("is_delayed"))

	remark := c.PostForm("remark")
	status := c.PostForm("status")
	if status == "" {
		status = "pending"
	}

	salesChargesJSON := c.PostForm("sales_charges")
	fmt.Println("Raw sales_charges JSON:", salesChargesJSON)

	var salesCharges []models.SalesCharge
	if salesChargesJSON != "" {
		if err := json.Unmarshal([]byte(salesChargesJSON), &salesCharges); err != nil {
			fmt.Println("Error unmarshaling sales_charges:", err)
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sales_charges format", err.Error()))
			return
		}
	}

	vehicleUsageJSON := c.PostForm("vehicle_usage")
	var vehicleUsage []models.VehicleUsage
	if vehicleUsageJSON != "" {
		if err := json.Unmarshal([]byte(vehicleUsageJSON), &vehicleUsage); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid vehicle_usage format", err.Error()))
			return
		}
	}

	paymentJSON := c.PostForm("payments")
	var payments []models.Payment
	if paymentJSON != "" {
		if err := json.Unmarshal([]byte(paymentJSON), &payments); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payments format", err.Error()))
			return
		}
	}

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
		Payments:       payments,
	}

	saleID, err := h.saleService.CreateSale(sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create sale", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Sale created successfully", gin.H{"sale_id": saleID}))
}

func (h *SaleHandler) GetSaleByID(c *gin.Context) {

	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}

	sale, err := h.saleService.GetSaleByID(saleID)
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
