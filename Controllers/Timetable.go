package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	initializers "github.com/group4/campus-connect-api/Initializers"
	models "github.com/group4/campus-connect-api/Models"
)

func GetTimetables(c *gin.Context) {
	var timetables []models.Timetable
	if err := initializers.DB.Find(&timetables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timetables"})
		return
	}
	c.JSON(http.StatusOK, timetables)
}

func CreateTimetable(c *gin.Context) {
	var timetable models.Timetable
	if err := c.ShouldBindJSON(&timetable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := initializers.DB.Create(&timetable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create timetable: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, timetable)
}

func GetTimetableByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timetable ID"})
		return
	}

	var timetable models.Timetable
	if err := initializers.DB.First(&timetable, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Timetable not found"})
		return
	}
	c.JSON(http.StatusOK, timetable)
}

func UpdateTimetable(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timetable ID"})
		return
	}

	var timetable models.Timetable
	if err := initializers.DB.First(&timetable, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Timetable not found"})
		return
	}

	var updatedTimetable models.Timetable
	if err := c.ShouldBindJSON(&updatedTimetable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	timetable.Day = updatedTimetable.Day
	timetable.Subject = updatedTimetable.Subject
	timetable.SubjectCode = updatedTimetable.SubjectCode
	timetable.Faculty = updatedTimetable.Faculty
	timetable.Room = updatedTimetable.Room
	timetable.Time = updatedTimetable.Time
	timetable.Instructor = updatedTimetable.Instructor

	if err := initializers.DB.Save(&timetable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update timetable"})
		return
	}
	c.JSON(http.StatusOK, timetable)
}

func DeleteTimetable(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timetable ID"})
		return
	}

	var timetable models.Timetable
	if err := initializers.DB.First(&timetable, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Timetable not found"})
		return
	}

	if err := initializers.DB.Delete(&timetable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete timetable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Timetable deleted successfully"})
}