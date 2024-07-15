package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func assignCouponToUser(c *gin.Context) {

	var newCoupon Coupon
	if err := c.BindJSON(&newCoupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCoupon.ID = len(coupons) + 1
	coupons = append(coupons, newCoupon)
	c.JSON(http.StatusOK, newCoupon)
}
