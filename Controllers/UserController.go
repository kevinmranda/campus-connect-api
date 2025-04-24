package controllers

import (
	"net/http"
	"strconv"
	"github.com/group4/campus-connect-api/Helpers"
	"github.com/gin-gonic/gin"
	initializers "github.com/group4/campus-connect-api/Initializers"
	models "github.com/group4/campus-connect-api/Models"
	"golang.org/x/crypto/bcrypt"
)

// User Controller (Modified)
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := initializers.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	user.Password = string(hashedPassword)

	base64Image := user.ProfileImage
	user.ProfileImage = "" 

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	// Handle profile image upload if provided
	if base64Image != "" {
		filename := helpers.SanitizeUserIDFilename("profile-picture-UID", user.ID)
		imagePath, err := helpers.SaveImage(base64Image, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "User created but failed to save profile image: " + err.Error(),
				"id":      user.ID,
				"name":    user.Name,
				"email":   user.Email,
				"role":    user.Role,
			})
			return
		}
		user.ProfileImage = imagePath
		if err := initializers.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "User created but failed to update profile image path: " + err.Error(),
				"id":      user.ID,
				"name":    user.Name,
				"email":   user.Email,
				"role":    user.Role,
			})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           user.ID,
		"name":         user.Name,
		"email":        user.Email,
		"role":         user.Role,
		"profileImage": user.ProfileImage,
	})
}

func GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Handle password update
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Handle profile image update if provided
	if updatedUser.ProfileImage != "" && updatedUser.ProfileImage != user.ProfileImage {
		// Delete old image if it exists
		if user.ProfileImage != "" {
			if err := helpers.DeleteImage(user.ProfileImage); err != nil {
				// Log warning but proceed
				c.JSON(http.StatusOK, gin.H{
					"warning": "Failed to delete old profile image: " + err.Error(),
					"user":    user,
				})
			}
		}
		// Save new image with sanitized user ID
		filename := helpers.SanitizeUserIDFilename("profile-picture-UID", uint(id))
		imagePath, err := helpers.SaveImage(updatedUser.ProfileImage, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image: " + err.Error()})
			return
		}
		user.ProfileImage = imagePath
	}

	// Update fields
	user.Name = updatedUser.Name
	user.Role = updatedUser.Role
	user.Course = updatedUser.Course
	user.Year = updatedUser.Year
	user.Email = updatedUser.Email

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete profile image if it exists
	if user.ProfileImage != "" {
		if err := helpers.DeleteImage(user.ProfileImage); err != nil {
			// Log warning but proceed with deletion
			c.JSON(http.StatusOK, gin.H{
				"message": "User deleted successfully",
				"warning": "Failed to delete profile image: " + err.Error(),
			})
			return
		}
	}

	if err := initializers.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func Login(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"role":         user.Role,
			"course":       user.Course,
			"year":         user.Year,
			"profileImage": user.ProfileImage,
		},
	})
}