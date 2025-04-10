package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	name := c.PostForm("name")
	establishDate := c.PostForm("establish_date")

	// Get logo file
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo is required"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logo must be a .jpg, .jpeg, or .png image"})
		return
	}

	// Create filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename))
	savePath := filepath.Join("media/logos", filename)

	// Create directory (if not exist)
	os.MkdirAll("media/logos", os.ModePerm)

	// Save file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo"})
		return
	}

	// Save to database
	partai := models.Party{
		Name:          name,
		EstablishDate: establishDate,
		Logo:          savePath,
	}
	if err := config.DB.Create(&partai).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create partai"})
		return
	}

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
