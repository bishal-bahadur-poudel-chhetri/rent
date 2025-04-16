package models

// In your models package
type DisableDateResponse struct {
	ActiveRentals  []DisabledDateResponse `json:"active_rentals"`
	FutureBookings []DisabledDateResponse `json:"future_bookings"` // Changed from FutureBooking
}

type DisabledDateResponse struct {
	DateOfDelivery string `json:"date_of_delivery"` // Format: YYYY-MM-DD
	ReturnDate     string `json:"return_date"`      // Format: YYYY-MM-DD
}

