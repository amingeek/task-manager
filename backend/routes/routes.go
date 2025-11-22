// backend/routes/routes.go

package routes

import (
	"task-manager/controllers"
	"task-manager/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ایجاد پوشه آپلود
	router.Static("/uploads", "./uploads")

	public := router.Group("/api")

	// Authentication routes
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// User
	protected.GET("/me", controllers.GetCurrentUser)
	protected.GET("/users/search", controllers.SearchUsers)
	protected.PUT("/profile", controllers.UpdateProfile)

	// Personal tasks
	protected.GET("/tasks", controllers.GetTasks)
	protected.GET("/tasks/:id", controllers.GetTask)
	protected.POST("/tasks", controllers.CreateTask)
	protected.PUT("/tasks/:id", controllers.UpdateTask)
	protected.DELETE("/tasks/:id", controllers.DeleteTask)

	// Personal task progress
	protected.PUT("/tasks/:id/progress", controllers.UpdatePersonalProgress)
	protected.GET("/tasks/:id/progress", controllers.GetPersonalProgress)

	// Groups
	protected.GET("/groups", controllers.GetGroups)
	protected.GET("/groups/:id", controllers.GetGroup)
	protected.POST("/groups", controllers.CreateGroup)
	protected.PUT("/groups/:id", controllers.UpdateGroup)
	protected.DELETE("/groups/:id", controllers.DeleteGroup)

	// Group members
	protected.POST("/groups/:id/members", controllers.AddMembers)
	protected.DELETE("/groups/:id/members/:user_id", controllers.RemoveMember)
	protected.POST("/groups/:id/members/:user_id/accept", controllers.AcceptInvitation)

	// Group tasks
	protected.POST("/groups/:id/tasks", controllers.CreateGroupTask)
	protected.PUT("/groups/:id/tasks/:task_id", controllers.UpdateGroupTask)
	protected.DELETE("/groups/:id/tasks/:task_id", controllers.DeleteGroupTask)
	protected.GET("/groups/:id/tasks", controllers.GetGroupTasks)

	// Group task progress
	protected.PUT("/groups/:id/tasks/:task_id/progress", controllers.UpdateGroupProgress)
	protected.GET("/groups/:id/tasks/:task_id/progress", controllers.GetGroupProgress)

	// Files
	protected.POST("/tasks/:id/files", controllers.UploadFile)
	protected.GET("/files/:id", controllers.DownloadFile)
	protected.DELETE("/files/:id", controllers.DeleteFile)

	// Notifications
	protected.GET("/notifications", controllers.GetNotifications)
	protected.PUT("/notifications/:id/read", controllers.MarkAsRead)
	protected.DELETE("/notifications/:id", controllers.DeleteNotification)

	// Analytics
	protected.GET("/analytics/streak", controllers.GetStreak)
	protected.GET("/analytics/summary", controllers.GetAnalyticsSummary)
}
