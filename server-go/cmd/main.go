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
	r.Use(middleware.RateLimitMiddleware())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/api/ws"},
	}))

	r.Static("/uploads", cfg.UploadPath)

	authHandler := handlers.NewAuthHandler(database.DB, cfg)
	bloggerHandler := handlers.NewBloggerHandler(database.DB)
	teamHandler := handlers.NewTeamHandler(database.DB)
	forumHandler := handlers.NewForumHandler(database.DB)
	miscHandler := handlers.NewMiscHandler(database.DB, cfg)
	adminHandler := handlers.NewAdminHandler(database.DB)
	messageHandler := handlers.NewMessageHandler(database.DB)
	calendarHandler := handlers.NewCalendarHandler()
	templateHandler := handlers.NewFollowupTemplateHandler()
	wsHandler := handlers.NewWSHandler()
	analyticsHandler := handlers.NewAnalyticsHandler()
	teamCollabHandler := handlers.NewTeamCollabHandler()
	importHandler := handlers.NewImportHandler()
	filterHandler := handlers.NewFilterHandler(database.DB)
	tagRecommendHandler := handlers.NewTagRecommendationHandler(database.DB)
	workflowHandler := handlers.NewWorkflowHandler(database.DB)

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
	authorized.Use(middleware.CSRFMiddleware())
	{
		authorized.GET("/ws", wsHandler.HandleWebSocket)

		authorized.GET("/user/profile", authHandler.GetProfile)
		authorized.PUT("/user/profile", authHandler.UpdateProfile)
		authorized.POST("/user/change-password", authHandler.ChangePassword)
		authorized.POST("/user/avatar", miscHandler.UploadUserAvatar)
		authorized.GET("/user/detail/:username", authHandler.GetUserDetail)
		authorized.GET("/user/:username/bloggers", authHandler.GetUserBloggers)
		authorized.GET("/user/:username/followup-stats", authHandler.GetUserFollowupStats)

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
		authorized.POST("/blogger/batch/status", bloggerHandler.BatchUpdateStatus)
		authorized.POST("/blogger/batch/tags", bloggerHandler.BatchUpdateTags)
		authorized.POST("/blogger/batch/delete", bloggerHandler.BatchDelete)
		authorized.GET("/saved-filters", filterHandler.GetSavedFilters)
		authorized.POST("/saved-filters", filterHandler.CreateSavedFilter)
		authorized.PUT("/saved-filters/:id", filterHandler.UpdateSavedFilter)
		authorized.DELETE("/saved-filters/:id", filterHandler.DeleteSavedFilter)
		authorized.POST("/saved-filters/:id/default", filterHandler.SetDefaultFilter)
		authorized.GET("/blogger/recommend-tags", tagRecommendHandler.GetRecommendations)
		authorized.GET("/blogger/similar", tagRecommendHandler.GetSimilarBloggers)
		authorized.GET("/workflow/rules", workflowHandler.GetRules)
		authorized.POST("/workflow/rules", workflowHandler.CreateRule)
		authorized.PUT("/workflow/rules/:id", workflowHandler.UpdateRule)
		authorized.DELETE("/workflow/rules/:id", workflowHandler.DeleteRule)
		authorized.POST("/workflow/rules/:id/toggle", workflowHandler.ToggleRule)
		authorized.POST("/workflow/trigger", workflowHandler.TriggerWorkflow)
		authorized.GET("/blogger/transfers", bloggerHandler.GetTransferRequests)
		authorized.POST("/blogger/transfers/:id/confirm", bloggerHandler.ConfirmTransferRequest)
		authorized.POST("/blogger/:id/invalid", bloggerHandler.SetInvalid)
		authorized.POST("/blogger/:id/valid", bloggerHandler.SetValid)
		authorized.POST("/blogger/evaluation", bloggerHandler.SubmitEvaluation)
		authorized.GET("/blogger/evaluation/:blogger_id", bloggerHandler.GetEvaluation)

		tagHandler := handlers.NewTagHandler()
		authorized.GET("/tags", tagHandler.GetTags)
		authorized.POST("/tags", tagHandler.CreateTag)
		authorized.PUT("/tags/:id", tagHandler.UpdateTag)
		authorized.DELETE("/tags/:id", tagHandler.DeleteTag)
		authorized.GET("/blogger/:id/tags", tagHandler.GetBloggerTags)
		authorized.POST("/blogger/:id/tags", tagHandler.SetBloggerTags)
		authorized.GET("/bloggers/by-tag", tagHandler.GetBloggersByTag)

		authorized.GET("/analytics/overview", analyticsHandler.GetOverview)
		authorized.GET("/analytics/blogger/:id", analyticsHandler.GetBloggerAnalytics)

		authorized.GET("/followup/list", bloggerHandler.GetFollowupList)
		authorized.POST("/followup/add", bloggerHandler.CreateFollowupV2)
		authorized.PUT("/followup/:id", bloggerHandler.UpdateFollowupV2)
		authorized.DELETE("/followup/:id", bloggerHandler.DeleteFollowupV2)

		authorized.GET("/cooperation/list", bloggerHandler.GetCooperationList)
		authorized.POST("/cooperation/add", bloggerHandler.CreateCooperationV2)
		authorized.PUT("/cooperation/:id", bloggerHandler.UpdateCooperationV2)
		authorized.DELETE("/cooperation/:id", bloggerHandler.DeleteCooperationV2)

		authorized.GET("/blogger/stat", miscHandler.GetBloggerStat)
		authorized.GET("/blogger/charts", miscHandler.GetBloggerCharts)
		authorized.GET("/blogger/status/list", miscHandler.GetBloggerStatusList)
		authorized.GET("/blogger/invalid/count", miscHandler.GetInvalidBloggerCount)
		authorized.GET("/blogger/invalid/list", bloggerHandler.GetInvalidBloggers)
		authorized.GET("/blogger/expiring-without-contact", bloggerHandler.GetExpiringBloggers)
		authorized.GET("/blogger/hot/:nickname", miscHandler.GetBloggerHotData)
		authorized.GET("/blogger/news/:nickname", miscHandler.GetBloggerNews)

		authorized.GET("/calendar/events", calendarHandler.GetCalendarEvents)

		authorized.GET("/team/members", teamCollabHandler.GetTeamMembers)
		authorized.POST("/blogger/transfer/request", teamCollabHandler.RequestTransfer)
		authorized.GET("/blogger/transfer/requests", teamCollabHandler.GetTransferRequests)
		authorized.POST("/blogger/transfer/handle", teamCollabHandler.HandleTransferRequest)
		authorized.GET("/team/tasks", teamCollabHandler.GetTeamTasks)
		authorized.POST("/team/tasks", teamCollabHandler.CreateTeamTask)
		authorized.PUT("/team/tasks/status", teamCollabHandler.UpdateTaskStatus)

		authorized.GET("/templates", templateHandler.GetTemplates)
		authorized.POST("/templates", templateHandler.CreateTemplate)
		authorized.PUT("/templates/:id", templateHandler.UpdateTemplate)
		authorized.DELETE("/templates/:id", templateHandler.DeleteTemplate)
		authorized.POST("/templates/:id/use", templateHandler.UseTemplate)
		authorized.POST("/templates/:id/set-default", templateHandler.SetDefault)

		authorized.POST("/upload/blogger-avatar", miscHandler.UploadBloggerAvatar)

		authorized.GET("/categories", miscHandler.GetCategories)
		authorized.POST("/categories", miscHandler.AddCategory)
		authorized.PUT("/categories/:id", miscHandler.UpdateCategory)
		authorized.DELETE("/categories/:id", miscHandler.DeleteCategory)
		authorized.POST("/products", miscHandler.AddProduct)
		authorized.PUT("/products/:id", miscHandler.UpdateProduct)
		authorized.DELETE("/products/:id", miscHandler.DeleteProduct)

		authorized.POST("/import/preview", importHandler.PreviewImport)
		authorized.POST("/import/do", importHandler.DoImport)
		authorized.GET("/import/template", importHandler.DownloadTemplate)
		authorized.GET("/import/history", importHandler.GetImportHistory)
		authorized.POST("/import/export-errors", importHandler.ExportFailedRows)

		authorized.GET("/users", authHandler.GetUsers)
		authorized.PUT("/users/:id", authHandler.UpdateUser)
		authorized.POST("/users/approve", authHandler.ApproveUser)
		authorized.DELETE("/users/:id", authHandler.DeleteUser)
		authorized.POST("/users/:id/reset-password", authHandler.ResetPassword)

		authorized.GET("/teams", teamHandler.GetTeams)
		authorized.GET("/team/:id", teamHandler.GetTeam)
		authorized.POST("/teams", teamHandler.CreateTeam)
		authorized.PUT("/teams/:id", teamHandler.UpdateTeam)
		authorized.DELETE("/teams/:id", teamHandler.DeleteTeam)
		authorized.GET("/team/posts", teamHandler.GetTeamPosts)
		authorized.GET("/team/posts/search", teamHandler.SearchTeamPosts)
		authorized.GET("/team/:id/posts", teamHandler.GetTeamPosts)
		authorized.GET("/team/posts/:id", teamHandler.GetTeamPost)
		authorized.POST("/team/posts", teamHandler.CreateTeamPost)
		authorized.POST("/team/:id/posts", teamHandler.CreateTeamPost)
		authorized.DELETE("/team/posts/:id", teamHandler.DeleteTeamPost)
		authorized.POST("/team/posts/:id/like", teamHandler.LikeTeamPost)
		authorized.GET("/team/posts/:id/like-status", teamHandler.GetLikeStatus)
		authorized.POST("/team/posts/:id/collect", teamHandler.CollectTeamPost)
		authorized.GET("/team/posts/collect-status/:id", teamHandler.GetCollectStatus)
		authorized.GET("/team/posts/collected", teamHandler.GetCollectedPosts)
		authorized.POST("/team/posts/:id/comments", teamHandler.CreateTeamPostComment)
		authorized.DELETE("/team/posts/:id/comments/:commentId", teamHandler.DeleteTeamPostComment)
		authorized.POST("/team/posts/:id/pin", teamHandler.PinTeamPost)
		authorized.POST("/team/posts/:id/feature", teamHandler.FeatureTeamPost)
		authorized.GET("/team/messages", teamHandler.GetTeamMessages)
		authorized.GET("/team/:id/messages", teamHandler.GetTeamMessages)
		authorized.POST("/team/messages", teamHandler.CreateTeamMessage)
		authorized.POST("/team/:id/messages", teamHandler.CreateTeamMessage)
		authorized.PUT("/team/:id/messages/:messageId/read", teamHandler.MarkTeamMessageRead)
		authorized.DELETE("/team/:id/messages/:messageId", teamHandler.DeleteTeamMessage)
		authorized.PUT("/team/:id/admins/:userId", teamHandler.SetTeamAdmin)
		authorized.DELETE("/team/:id/members/:userId", teamHandler.RemoveTeamMember)
		authorized.GET("/team/:id/logs", teamHandler.GetTeamOperationLogs)
		authorized.POST("/upload/team-logo", miscHandler.UploadTeamLogo)
		authorized.POST("/upload/team-bg", miscHandler.UploadTeamBgImage)

		authorized.POST("/forum/posts", forumHandler.CreatePost)
		authorized.GET("/forum/posts/search", forumHandler.SearchForumPosts)
		authorized.GET("/forum/posts/hot", forumHandler.GetHotPosts)
		authorized.GET("/forum/posts/collected", forumHandler.GetCollectedPosts)
		authorized.GET("/forum/posts/:id/like-status", forumHandler.GetLikeStatus)
		authorized.GET("/forum/posts/:id/collect-status", forumHandler.GetCollectStatus)
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
		authorized.DELETE("/notifications/:id", miscHandler.DeleteNotification)
		authorized.POST("/notifications/batch-delete", miscHandler.BatchDeleteNotifications)

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

		authorized.GET("/team/blogger/stat", miscHandler.GetTeamBloggerStat)
		authorized.GET("/team/blogger/charts", miscHandler.GetTeamBloggerCharts)

		authorized.GET("/messages/conversations", messageHandler.GetConversations)
		authorized.GET("/messages", messageHandler.GetMessages)
		authorized.POST("/messages", messageHandler.SendMessage)
		authorized.POST("/messages/read", messageHandler.MarkAsRead)
		authorized.GET("/messages/unread", messageHandler.GetUnreadCount)
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
		admin.GET("/export", miscHandler.ExportData)
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
