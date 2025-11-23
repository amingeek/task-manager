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

	// ==================== User Routes ====================
	protected.GET("/me", controllers.GetCurrentUser)
	protected.GET("/users/search", controllers.SearchUsers)
	protected.PUT("/profile", controllers.UpdateProfile)

	// ==================== Personal Task Routes ====================
	protected.GET("/tasks", controllers.GetTasks)
	protected.GET("/tasks/:id", controllers.GetTask)
	protected.POST("/tasks", controllers.CreateTask)
	protected.PUT("/tasks/:id", controllers.UpdateTask)
	protected.DELETE("/tasks/:id", controllers.DeleteTask)

	// Personal task progress
	protected.PUT("/tasks/:id/progress", controllers.UpdatePersonalProgress)
	protected.GET("/tasks/:id/progress", controllers.GetPersonalProgress)

	// ==================== Group Routes ====================
	protected.GET("/groups", controllers.GetUserGroups)
	protected.POST("/groups", controllers.CreateGroup)        // اضافه شد
	protected.GET("/groups/:id", controllers.GetGroupDetails) // اضافه شد
	protected.PUT("/groups/:id", controllers.UpdateGroup)     // اضافه شد
	protected.DELETE("/groups/:id", controllers.DeleteGroup)  // اضافه شد

	// Group members
	protected.POST("/groups/:id/members", controllers.AddGroupMembers)
	protected.DELETE("/groups/:id/members/:user_id", controllers.RemoveMember)
	protected.POST("/groups/:id/members/:user_id/accept", controllers.AcceptInvitation)

	// ==================== Group Task Routes ====================
	protected.POST("/groups/:id/tasks", controllers.CreateGroupTask)
	protected.PUT("/groups/:id/tasks/:task_id", controllers.UpdateGroupTask)
	protected.DELETE("/groups/:id/tasks/:task_id", controllers.DeleteGroupTask)
	protected.GET("/groups/:id/tasks", controllers.GetGroupTasks)

	// Group task progress (تسک گروهی)
	protected.PUT("/groups/:id/tasks/:task_id/progress", controllers.UpdateGroupProgress)
	protected.GET("/groups/:id/tasks/:task_id/progress", controllers.GetGroupProgress)
	protected.GET("/groups/invitations", controllers.GetPendingInvitations) // جدید
	protected.POST("/groups/:id/accept", controllers.AcceptInvitation)      // تغییر مسیر برای سادگی

	// Approve group task files (تایید فایل توسط مدیر)
	protected.POST("/groups/:id/tasks/:task_id/files/approve", controllers.ApproveGroupTaskFile)

	// ==================== File Routes ====================

	// آپلود فایل برای تسک (شخصی یا گروهی)
	protected.POST("/tasks/:id/files", controllers.UploadFile)

	// دریافت فایلهای یک تسک
	protected.GET("/tasks/:id/files", controllers.GetTaskFiles)

	// دریافت فایلهای آپلود شده توسط یک کاربر برای تسک گروهی
	protected.GET("/groups/:id/tasks/:task_id/files/:user_id", controllers.GetGroupTaskFilesByUser)

	// دانلود فایل
	protected.GET("/files/:id", controllers.DownloadFile)

	// حذف فایل
	protected.DELETE("/files/:id", controllers.DeleteFile)

	// ==================== Notification Routes ====================
	protected.GET("/notifications", controllers.GetNotifications)
	protected.PUT("/notifications/:id/read", controllers.MarkAsRead)
	protected.DELETE("/notifications/:id", controllers.DeleteNotification)

	// ==================== Analytics Routes ====================
	protected.GET("/analytics/streak", controllers.GetStreak)
	protected.GET("/analytics/summary", controllers.GetAnalyticsSummary)
}
