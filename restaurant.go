package main

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func registerRestaurant(c *gin.Context) {
	var newRestaurant Restaurant
	if err := c.BindJSON(&newRestaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newRestaurant.ID = len(restaurants) + 1
	restaurants = append(restaurants, newRestaurant)
	c.JSON(http.StatusOK, newRestaurant)
}

func suggestRestaurants(c *gin.Context) {

	itemName := c.Query("itemName")
	maxTimeExpected := cast.ToInt(c.Query("maxTimeExpected"))
	userId := cast.ToInt(c.Query("userId"))

	boolItemExists := false

	var suggested []SuggestedRestaurant

	/*
		Here we are assuming delivery boy is closer to the restaurant and any one of the rider is always available
		so time taken will be user distance to restaurant distance + average preparation time
	*/

	for _, menu := range menus {

		if menu.ItemName == itemName {

			restaurant := getRestaurantDetailsBasedOnId(menu.RestaurantID)

			estimatedDeliveryTime := FindDeliveryTimeStrategy(riderCloseToRestaurant{}).calculateDeliveryTime(menu.RestaurantID, userId, 0)

			if estimatedDeliveryTime <= maxTimeExpected {

				boolItemExists = true

				suggested = append(suggested, SuggestedRestaurant{
					Restaurant:       restaurant,
					EstimatedTimeMin: estimatedDeliveryTime,
				})
			}
		}
	}

	if !boolItemExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Restaurant found for the given item with expected time."})
		return
	}

	/*
		Sort the suggested restaurants
		based on the estimated time to deliver
	*/

	sort.Slice(suggested, func(i, j int) bool {
		return suggested[i].EstimatedTimeMin < suggested[j].EstimatedTimeMin
	})

	c.JSON(http.StatusOK, suggested)
}

func getRestaurantDetailsBasedOnId(restaurantID int) Restaurant {

	for _, restaurant := range restaurants {
		if restaurant.ID == restaurantID {
			return restaurant
		}
	}
	return Restaurant{}

}
