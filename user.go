package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser.ID = len(users) + 1
	users = append(users, newUser)
	c.JSON(http.StatusOK, newUser)
}

func getUserDetailsBasedOnId(id int) User {
	for _, user := range users {
		if user.ID == id {
			return user
		}
	}
	return User{}

}

func updateUserDetailsBasedOnId(newUser User) {
	for i, user := range users {
		if user.ID == newUser.ID {
			users[i] = newUser
			break
		}
	}

}
