package controllers

import (
	"net/http"
	"strconv"
	"github.com/group4/campus-connect-api/Helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	initializers "github.com/group4/campus-connect-api/Initializers"
	models "github.com/group4/campus-connect-api/Models"
)

// Response struct to control user fields in the response
type PostResponse struct {
	gorm.Model
	Image       string       `json:"image"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	UserID      uint         `json:"userID"`
	User        UserResponse `json:"user"`
}

// UserResponse struct to exclude sensitive fields
type UserResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProfileImage string `json:"profileImage"`
	Role         string `json:"role"`
	Course       string `json:"course"`
	Year         string `json:"year"`
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := initializers.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// Convert to response struct to control fields
	var postResponses []PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, PostResponse{
			Model:       post.Model,
			Image:       post.Image,
			Title:       post.Title,
			Description: post.Description,
			UserID:      post.UserID,
			User: UserResponse{
				ID:           post.User.ID,
				Name:         post.User.Name,
				ProfileImage: post.User.ProfileImage,
				Role:         post.User.Role,
				Course:       post.User.Course,
				Year:         post.User.Year,
			},
		})
	}

	c.JSON(http.StatusOK, postResponses)
}

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate UserID exists
	var user models.User
	if err := initializers.DB.First(&user, post.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID: User not found"})
		return
	}

	// Handle image upload if provided
	if post.Image != "" {
		// Sanitize title for filename
		sanitizedTitle := helpers.SanitizeFilename(post.Title)
		filename := "post-" + sanitizedTitle
		imagePath, err := helpers.SaveImage(post.Image, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
			return
		}
		post.Image = imagePath // Store relative path in DB
	}

	if err := initializers.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post: " + err.Error()})
		return
	}

	// Return post with user details
	postResponse := PostResponse{
		Model:       post.Model,
		Image:       post.Image,
		Title:       post.Title,
		Description: post.Description,
		UserID:      post.UserID,
		User: UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			Role:         user.Role,
			Course:       user.Course,
			Year:         user.Year,
		},
	}

	c.JSON(http.StatusCreated, postResponse)
}

func GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := initializers.DB.Preload("User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Convert to response struct
	postResponse := PostResponse{
		Model:       post.Model,
		Image:       post.Image,
		Title:       post.Title,
		Description: post.Description,
		UserID:      post.UserID,
		User: UserResponse{
			ID:           post.User.ID,
			Name:         post.User.Name,
			ProfileImage: post.User.ProfileImage,
			Role:         post.User.Role,
			Course:       post.User.Course,
			Year:         post.User.Year,
		},
	}

	c.JSON(http.StatusOK, postResponse)
}

func UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := initializers.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var updatedPost models.Post
	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Validate UserID if provided
	if updatedPost.UserID != 0 {
		var user models.User
		if err := initializers.DB.First(&user, updatedPost.UserID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID: User not found"})
			return
		}
		post.UserID = updatedPost.UserID
	}

	// Handle image update if provided
	if updatedPost.Image != "" && updatedPost.Image != post.Image {
		// Delete old image if it exists
		if post.Image != "" {
			if err := helpers.DeleteImage(post.Image); err != nil {
				// Log warning but proceed
				c.JSON(http.StatusOK, gin.H{
					"warning": "Failed to delete old image: " + err.Error(),
					"post":    post,
				})
			}
		}
		// Save new image with sanitized title
		sanitizedTitle := helpers.SanitizeFilename(updatedPost.Title)
		filename := "post-" + sanitizedTitle
		imagePath, err := helpers.SaveImage(updatedPost.Image, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image: " + err.Error()})
			return
		}
		post.Image = imagePath
	}

	// Update other fields
	post.Title = updatedPost.Title
	post.Description = updatedPost.Description

	if err := initializers.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Fetch updated user details
	var user models.User
	initializers.DB.First(&user, post.UserID)

	// Return updated post with user details
	postResponse := PostResponse{
		Model:       post.Model,
		Image:       post.Image,
		Title:       post.Title,
		Description: post.Description,
		UserID:      post.UserID,
		User: UserResponse{
			ID:           user.ID,
			Name:         user.Name,
			ProfileImage: user.ProfileImage,
			Role:         user.Role,
			Course:       user.Course,
			Year:         user.Year,
		},
	}

	c.JSON(http.StatusOK, postResponse)
}

func DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := initializers.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Delete image file if it exists
	if post.Image != "" {
		if err := helpers.DeleteImage(post.Image); err != nil {
			// Log the error but proceed with deletion
			c.JSON(http.StatusOK, gin.H{
				"message": "Post deleted successfully",
				"warning": "Failed to delete associated image: " + err.Error(),
			})
			return
		}
	}

	if err := initializers.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}