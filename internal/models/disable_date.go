package models

type DisableDateResponse struct {
	FutureBooking []DisabledDateResponse `json:"futureBooking"`
	TodaySales    []DisabledDateResponse `json:"todaySales"`
}

type DisabledDateResponse struct {
	DateOfDelivery string `json:"date_of_delivery"` // Format: YYYY-MM-DD
	ReturnDate     string `json:"return_date"`      // Format: YYYY-MM-DD
}
