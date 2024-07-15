package main

import (
	"math"

	"github.com/spf13/cast"
)

type Location struct {
	Lat int `json:"lat"`
	Lon int `json:"lon"`
}

type User struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Location  Location `json:"location"`
	AvgRating float64  `json:"avg_rating"`
}

type Rider struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	Location    Location `json:"location"`
	IsAvailable bool     `json:"is_available"`
	AvgRating   float64  `json:"avg_rating"`
}

type Restaurant struct {
	ID                   int      `json:"id"`
	Name                 string   `json:"name"`
	Location             Location `json:"location"`
	AveragePreparingTime int      `json:"avg_preparing_time"`
}

type SuggestedRestaurant struct {
	Restaurant       Restaurant `json:"restaurant"`
	EstimatedTimeMin int        `json:"estimated_time_min"`
}

type Menu struct {
	ID           int     `json:"id"`
	RestaurantID int     `json:"restaurant_id"`
	ItemName     string  `json:"item_name"`
	ItemPrice    float64 `json:"item_price"`
}

type Coupon struct {
	ID         int     `json:"id"`
	Code       string  `json:"code"`
	Discount   float64 `json:"discount"`
	ValidUntil string  `json:"valid_until"`
	UserID     int     `json:"user_id"`
}

type Order struct {
	ID                   int      `json:"id"`
	UserID               int      `json:"user_id"`
	RestaurantID         int      `json:"restaurant_id"`
	RiderID              int      `json:"rider_id"`
	OrderStatus          string   `json:"order_status"`
	OrderTime            string   `json:"order_time"`
	TotalAmount          float64  `json:"total_amount"`
	ItemNames            []string `json:"item_names"`
	ExpectedDeliveryTime int      `json:"expected_delivery_time"`
	CouponID             int      `json:"coupon_id"`
}

type Rating struct {
	ID      int     `json:"id"`
	UserID  int     `json:"user_id,omitempty"`
	RiderID int     `json:"rider_id,omitempty"`
	Rating  float64 `json:"rating"`
	OrderID int     `json:"order_id"`
}

func calculateDistance(lat1, long1, lat2, long2 int) int {
	return cast.ToInt(math.Sqrt(math.Pow(cast.ToFloat64(lat1)-cast.ToFloat64(lat2), 2) + math.Pow(cast.ToFloat64(long1)-cast.ToFloat64(long2), 2)))
}
