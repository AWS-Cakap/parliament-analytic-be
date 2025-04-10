package controllers

import (
	"net/http"
	"parliament-analytic-be/config"
	"parliament-analytic-be/models"

	"github.com/gin-gonic/gin"
)

// GET /admin/partai
func GetAllPartai(c *gin.Context) {
	var partai []models.Party
	config.DB.Find(&partai)
	c.JSON(http.StatusOK, gin.H{"data": partai})
}

// GET /admin/partai/:id
func GetPartaiByID(c *gin.Context) {
	id := c.Param("id")
	var partai models.Party

	if err := config.DB.First(&partai, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": partai})
}

// POST /admin/partai
func CreatePartai(c *gin.Context) {
	var input struct {
		Name          string `json:"name"`
		EstablishDate string `json:"establish_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input!"})
		return
	}

	partai := models.Party{
		Name:          input.Name,
		EstablishDate: input.EstablishDate,
	}
	config.DB.Create(&partai)

	c.JSON(http.StatusCreated, gin.H{"data": partai})
}

// PUT /admin/partai/:id
func UpdatePartai(c *gin.Context) {
	id := c.Param("id")
	var partai models.Party

	if err := config.DB.First(&partai, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found!"})
		return
	}

	var input struct {
		Name          string `json:"name"`
		EstablishDate string `json:"establish_date"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input!"})
		return
	}

	partai.Name = input.Name
	partai.EstablishDate = input.EstablishDate
	config.DB.Save(&partai)

	c.JSON(http.StatusOK, gin.H{"data": partai})
}

// DELETE /admin/partai/:id
func DeletePartai(c *gin.Context) {
	id := c.Param("id")
	var partai models.Party

	if err := config.DB.First(&partai, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Party not found!"})
		return
	}

	config.DB.Delete(&partai)
	c.JSON(http.StatusOK, gin.H{"message": "Party has been deleted"})
}
