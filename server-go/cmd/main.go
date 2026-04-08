package main

import (
	"log"
	"server-go/internal/config"
	"server-go/internal/database"
	"server-go/internal/handlers"
	"server-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if err := database.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	r.Static("/uploads", cfg.UploadPath)

	authHandler := handlers.NewAuthHandler(database.DB, cfg)
	bloggerHandler := handlers.NewBloggerHandler(database.DB)
	teamHandler := handlers.NewTeamHandler(database.DB)
	forumHandler := handlers.NewForumHandler(database.DB)
	miscHandler := handlers.NewMiscHandler(database.DB, cfg)
	adminHandler := handlers.NewAdminHandler(database.DB)

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		api.GET("/public/posts", forumHandler.GetPosts)
		api.GET("/public/posts/:id", forumHandler.GetPost)
		api.GET("/public/forums", forumHandler.GetPosts)
		api.GET("/public/forums/:id", forumHandler.GetPost)

		api.GET("/category/list", miscHandler.GetCategories)
		api.GET("/platform/list", miscHandler.GetPlatforms)
		api.GET("/product/list", miscHandler.GetProducts)

		api.GET("/users/public", miscHandler.GetPublicUsers)
	}

	authorized := api.Group("/")
	authorized.Use(middleware.AuthMiddleware(cfg))
	{
		authorized.GET("/user/profile", authHandler.GetProfile)
		authorized.PUT("/user/profile", authHandler.UpdateProfile)
		authorized.POST("/user/change-password", authHandler.ChangePassword)
		authorized.POST("/user/avatar", miscHandler.UploadUserAvatar)

		authorized.GET("/blogger/list", bloggerHandler.GetBloggers)
		authorized.GET("/blogger/my", bloggerHandler.GetMyBloggers)
		authorized.GET("/blogger/:id", bloggerHandler.GetBlogger)
		authorized.POST("/blogger/add", bloggerHandler.CreateBlogger)
		authorized.PUT("/blogger/:id", bloggerHandler.UpdateBlogger)
		authorized.DELETE("/blogger/:id", bloggerHandler.DeleteBlogger)
		authorized.GET("/blogger/:id/history", bloggerHandler.GetStatusHistory)
		authorized.GET("/blogger/:id/followups", bloggerHandler.GetFollowups)
		authorized.POST("/blogger/:id/followups", bloggerHandler.CreateFollowup)
		authorized.DELETE("/blogger/:id/followups/:followupId", bloggerHandler.DeleteFollowup)
		authorized.GET("/blogger/:id/cooperations", bloggerHandler.GetCooperations)
		authorized.POST("/blogger/:id/cooperations", bloggerHandler.CreateCooperation)
		authorized.PUT("/blogger/:id/cooperations/:coopId", bloggerHandler.UpdateCooperation)
		authorized.DELETE("/blogger/:id/cooperations/:coopId", bloggerHandler.DeleteCooperation)
		authorized.POST("/blogger/:id/transfer", bloggerHandler.CreateTransferRequest)
		authorized.GET("/blogger/transfers", bloggerHandler.GetTransferRequests)
		authorized.POST("/blogger/transfers/:id/confirm", bloggerHandler.ConfirmTransferRequest)
		authorized.POST("/blogger/:id/invalid", bloggerHandler.SetInvalid)
		authorized.POST("/blogger/:id/valid", bloggerHandler.SetValid)

		authorized.GET("/blogger/stat", miscHandler.GetBloggerStat)
		authorized.GET("/blogger/charts", miscHandler.GetBloggerCharts)
		authorized.GET("/blogger/status/list", miscHandler.GetBloggerStatusList)
		authorized.GET("/blogger/invalid/count", miscHandler.GetInvalidBloggerCount)
		authorized.GET("/blogger/expiring-without-contact", bloggerHandler.GetExpiringBloggers)
		authorized.GET("/blogger/hot/:nickname", miscHandler.GetBloggerHotData)
		authorized.GET("/blogger/news/:nickname", miscHandler.GetBloggerNews)

		authorized.POST("/upload/blogger-avatar", miscHandler.UploadBloggerAvatar)

		authorized.GET("/categories", miscHandler.GetCategories)
		authorized.POST("/categories", miscHandler.AddCategory)
		authorized.PUT("/categories/:id", miscHandler.UpdateCategory)
		authorized.DELETE("/categories/:id", miscHandler.DeleteCategory)
		authorized.POST("/products", miscHandler.AddProduct)

		authorized.GET("/users", authHandler.GetUsers)
		authorized.PUT("/users/:id", authHandler.UpdateUser)
		authorized.POST("/users/:id/approve", authHandler.ApproveUser)
		authorized.DELETE("/users/:id", authHandler.DeleteUser)
		authorized.POST("/users/:id/reset-password", authHandler.ResetPassword)

		authorized.GET("/teams", teamHandler.GetTeams)
		authorized.GET("/team/:id", teamHandler.GetTeam)
		authorized.POST("/teams", teamHandler.CreateTeam)
		authorized.PUT("/teams/:id", teamHandler.UpdateTeam)
		authorized.DELETE("/teams/:id", teamHandler.DeleteTeam)
		authorized.GET("/team/posts", teamHandler.GetTeamPosts)
		authorized.GET("/team/:id/posts", teamHandler.GetTeamPosts)
		authorized.GET("/team/posts/:id", teamHandler.GetTeamPost)
		authorized.POST("/team/posts", teamHandler.CreateTeamPost)
		authorized.POST("/team/:id/posts", teamHandler.CreateTeamPost)
		authorized.DELETE("/team/posts/:id", teamHandler.DeleteTeamPost)
		authorized.POST("/team/posts/:id/like", teamHandler.LikeTeamPost)
		authorized.POST("/team/posts/:id/collect", teamHandler.CollectTeamPost)
		authorized.POST("/team/posts/:id/comments", teamHandler.CreateTeamPostComment)
		authorized.DELETE("/team/posts/:id/comments/:commentId", teamHandler.DeleteTeamPostComment)
		authorized.POST("/team/posts/:id/pin", teamHandler.PinTeamPost)
		authorized.POST("/team/posts/:id/feature", teamHandler.FeatureTeamPost)
		authorized.GET("/team/messages", teamHandler.GetTeamMessages)
		authorized.GET("/team/:id/messages", teamHandler.GetTeamMessages)
		authorized.POST("/team/messages", teamHandler.CreateTeamMessage)
		authorized.POST("/team/:id/messages", teamHandler.CreateTeamMessage)
		authorized.POST("/team/logo", miscHandler.UploadTeamLogo)
		authorized.POST("/team/bgimage", miscHandler.UploadTeamBgImage)
		authorized.POST("/upload/team-logo", miscHandler.UploadTeamLogo)
		authorized.POST("/upload/team-bg", miscHandler.UploadTeamBgImage)

		authorized.POST("/forum/posts", forumHandler.CreatePost)
		authorized.POST("/forum/posts/:id/comments", forumHandler.CreateComment)
		authorized.POST("/forum/posts/:id/like", forumHandler.LikePost)
		authorized.POST("/forum/posts/:id/collect", forumHandler.CollectPost)
		authorized.DELETE("/forum/posts/:id", forumHandler.DeletePost)
		authorized.DELETE("/forum/posts/:id/comments/:commentId", forumHandler.DeleteComment)
		authorized.POST("/forum/posts/:id/pin", forumHandler.PinPost)
		authorized.POST("/forum/posts/:id/feature", forumHandler.FeaturePost)

		authorized.GET("/notifications", miscHandler.GetNotifications)
		authorized.POST("/notifications/:id/read", miscHandler.MarkNotificationRead)
		authorized.POST("/notifications/mark-all-read", miscHandler.MarkAllNotificationsRead)

		authorized.GET("/announcements", miscHandler.GetAnnouncements)
		authorized.POST("/announcements", miscHandler.CreateAnnouncement)

		authorized.GET("/logs", miscHandler.GetOperationLog)

		authorized.GET("/snapshots", miscHandler.GetSnapshots)
		authorized.POST("/snapshots", miscHandler.CreateSnapshot)
		authorized.DELETE("/snapshots/:id", miscHandler.DeleteSnapshot)
		authorized.GET("/snapshots/settings", miscHandler.GetSnapshotSettings)
		authorized.POST("/snapshots/settings", miscHandler.SetSnapshotSettings)
		authorized.GET("/snapshots/file/:filename", miscHandler.DownloadSnapshot)
		authorized.POST("/snapshots/file/:filename/restore", miscHandler.RestoreSnapshot)
		authorized.POST("/backup", miscHandler.CreateBackup)
		authorized.POST("/clear", miscHandler.ClearData)

		authorized.GET("/users/pending", adminHandler.GetPendingUsers)
		authorized.GET("/users/list", adminHandler.GetAllUsers)

		authorized.GET("/my/team", teamHandler.GetMyTeam)
		authorized.PUT("/my/team", teamHandler.UpdateMyTeam)
	}

	admin := authorized.Group("/admin")
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/registration-mode", adminHandler.GetRegistrationMode)
		admin.POST("/registration-mode", adminHandler.SetRegistrationMode)
		admin.GET("/log", adminHandler.GetLogs)
		admin.DELETE("/log/clear", adminHandler.ClearLogs)
		admin.DELETE("/log/clear/:count", adminHandler.ClearEarlyLogs)
		admin.POST("/blogger/delete", bloggerHandler.DeleteBlogger)
		admin.POST("/purge", miscHandler.PurgeData)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Go Backend Server starting on port %s", cfg.ServerPort)
	log.Printf("API endpoints configured completely")

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
