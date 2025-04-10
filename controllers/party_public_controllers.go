package controllers

import (
	"net/http"
	"parliament-analytic-be/config"
	"parliament-analytic-be/models"

	"github.com/gin-gonic/gin"
)

// Get all party available
func GetAllPartaiPublic(c *gin.Context) {
	var parties []models.Party
	// Party finder
	if err := config.DB.Find(&parties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching party data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": parties})
}

// GetPartaiDetailPublic to get party detail
func GetPartaiDetailPublic(c *gin.Context) {
	id := c.Param("id")
	var party models.Party

	// Fetch based on ID
	if err := config.DB.First(&party, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": party})
}
