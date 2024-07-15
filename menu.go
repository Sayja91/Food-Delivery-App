package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func getMenu(c *gin.Context) {

	restaurantId := cast.ToInt(c.Query("restaurantId"))
	boolIsRestaurantIdPresent := false

	var menuList []Menu
	for _, menu := range menus {
		if menu.RestaurantID == restaurantId {
			menuList = append(menuList, menu)
			boolIsRestaurantIdPresent = true
		}
	}

	if !boolIsRestaurantIdPresent {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Menu found for the given restaurant id."})
		return
	}
	c.JSON(http.StatusOK, menuList)
}

func addMenu(c *gin.Context) {
	var newMenu Menu
	if err := c.BindJSON(&newMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newMenu.ID = len(menus) + 1
	menus = append(menus, newMenu)
	c.JSON(http.StatusOK, newMenu)
}
