package routes

import (
	"github.com/gin-gonic/gin"
	"task-manager/controllers"
	"task-manager/middleware"
)

func SetupRoutes(router *gin.Engine) {
	// Public routes
	public := router.Group("/api")
	{
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/me", controllers.GetCurrentUser)
		protected.PUT("/profile", controllers.UpdateProfile)
		protected.GET("/users/search", controllers.SearchUsers)

		// Personal Tasks routes
		protected.GET("/tasks", controllers.GetTasks)
		protected.GET("/tasks/:id", controllers.GetTask)
		protected.POST("/tasks", controllers.CreateTask)
		protected.PUT("/tasks/:id", controllers.UpdateTask)
		protected.DELETE("/tasks/:id", controllers.DeleteTask)
		protected.PUT("/tasks/:id/progress", controllers.UpdatePersonalProgress)
		protected.GET("/tasks/:id/progress", controllers.GetPersonalProgress)

		// File routes
		protected.POST("/tasks/:id/files", controllers.UploadFile)
		protected.GET("/tasks/:id/files", controllers.GetTaskFiles)
		protected.GET("/files/:id", controllers.DownloadFile)
		protected.DELETE("/files/:id", controllers.DeleteFile)

		// Analytics routes
		protected.GET("/analytics/streak", controllers.GetStreak)
		protected.GET("/analytics/summary", controllers.GetAnalyticsSummary)

		// Notification routes
		protected.GET("/notifications", controllers.GetNotifications)
		protected.PUT("/notifications/:id/read", controllers.MarkAsRead)
		protected.DELETE("/notifications/:id", controllers.DeleteNotification)

		// ✅ Group routes - اصلاح‌شده
		protected.GET("/groups", controllers.GetUserGroups)
		protected.POST("/groups", controllers.CreateGroup)
		protected.GET("/groups/:id", controllers.GetGroupDetails)
		protected.PUT("/groups/:id", controllers.UpdateGroup)
		protected.DELETE("/groups/:id", controllers.DeleteGroup)

		// Group Members routes
		protected.POST("/groups/:id/members", controllers.AddGroupMembers)
		protected.GET("/groups/:id/members", controllers.GetGroupMembers)
		protected.DELETE("/groups/:id/members/:user_id", controllers.RemoveMember)

		// ✅ Invitations routes - جدید
		protected.GET("/groups/invitations", controllers.GetPendingInvitations)
		protected.POST("/groups/:id/accept-invitation", controllers.AcceptInvitation)
		protected.POST("/groups/:id/reject-invitation", controllers.RejectInvitation)

		// Group Tasks routes
		protected.GET("/groups/:id/tasks", controllers.GetGroupTasks)
		protected.POST("/groups/:id/tasks", controllers.CreateGroupTask)
		protected.PUT("/groups/:id/tasks/:task_id", controllers.UpdateGroupTask)
		protected.DELETE("/groups/:id/tasks/:task_id", controllers.DeleteGroupTask)
		protected.PUT("/groups/:id/tasks/:task_id/progress", controllers.UpdateGroupProgress)
		protected.GET("/groups/:id/tasks/:task_id/progress", controllers.GetGroupProgress)

		// Group Task Files routes
		protected.GET("/groups/:id/tasks/:task_id/files/:user_id", controllers.GetGroupTaskFilesByUser)
		protected.POST("/groups/:id/tasks/:task_id/files/approve", controllers.ApproveGroupTaskFile)
	}

	// Static files
	router.Static("/uploads", "./uploads")
}
