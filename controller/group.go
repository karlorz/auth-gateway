package controller

import (
	"net/http"

	"github.com/karlorz/auth-gateway/model"

	"github.com/gin-gonic/gin"
)

// GetAllGroups returns all distinct user groups
func GetAllGroups(c *gin.Context) {
	groups, err := model.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    groups,
	})
}
