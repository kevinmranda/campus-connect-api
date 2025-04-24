package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/group4/campus-connect-api/Controllers"
)

// all web routes defined here
func Routes() {
	r := gin.Default()
	r.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}),
	)

	// Serve static files from the uploads directory (public)
	r.Static("/Images", "./Images")

	// Public routes (no authentication required)
	// User routes
	r.POST("/api/user/register", controllers.CreateUser)
	r.POST("/api/user/login", controllers.Login)
	r.GET("/api/users", controllers.GetUsers)
	r.GET("/api/users/:id", controllers.GetUserByID)
	r.PUT("/api/users/:id/update", controllers.UpdateUser)
	r.DELETE("/api/users/:id/delete", controllers.DeleteUser)

	// Post routes
	r.GET("/api/posts", controllers.GetPosts)
	r.POST("/api/posts", controllers.CreatePost)
	r.GET("/api/posts/:id", controllers.GetPostByID)
	r.PUT("/api/posts/:id/update", controllers.UpdatePost)
	r.DELETE("/api/posts/:id/delete", controllers.DeletePost)

	// Job routes
	r.GET("/api/jobs", controllers.GetJobs)
	r.POST("/api/jobs", controllers.CreateJob)
	r.GET("/api/jobs/:id", controllers.GetJobByID)
	r.PUT("/api/jobs/:id/update", controllers.UpdateJob)
	r.DELETE("/api/jobs/:id/delete", controllers.DeleteJob)

	// Event routes
	r.GET("/api/events", controllers.GetEvents)
	r.POST("/api/events", controllers.CreateEvent)
	r.GET("/api/events/:id", controllers.GetEventByID)
	r.PUT("/api/events/:id/update", controllers.UpdateEvent)
	r.DELETE("/api/events/:id/delete", controllers.DeleteEvent)

	// Timetable routes
	r.GET("/api/timetables", controllers.GetTimetables)
	r.POST("/api/timetables", controllers.CreateTimetable)
	r.GET("/api/timetables/:id", controllers.GetTimetableByID)
	r.PUT("/api/timetables/:id/update", controllers.UpdateTimetable)
	r.DELETE("/api/timetables/:id/delete", controllers.DeleteTimetable)

	r.Run()
}