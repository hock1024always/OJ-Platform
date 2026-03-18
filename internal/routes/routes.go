package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/oj-platform/internal/handlers"
	"github.com/your-org/oj-platform/internal/middleware"
)

func Setup(r *gin.Engine, judgeHandler *handlers.JudgeHandler, userHandler *handlers.UserHandler, leaderboardHandler *handlers.LeaderboardHandler) {
	// 全局CORS中间件
	r.Use(middleware.CORS())

	// 全局限流：100 req/s per IP（令牌桶）
	r.Use(middleware.RateLimit())

	// 健康检查
	r.GET("/health", handlers.HealthCheck)

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", handlers.Welcome)

		// 用户相关
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.GET("/profile", middleware.AuthRequired(), userHandler.GetProfile)

		// 题目相关
		v1.GET("/problems", judgeHandler.ListProblems)
		v1.GET("/problems/:id", judgeHandler.GetProblem)
		v1.POST("/problems", middleware.AuthRequired(), judgeHandler.CreateProblem)       // 简单创建（管理员）
		v1.POST("/problems/import", middleware.AuthRequired(), judgeHandler.ImportProblem) // JSON批量导入（管理员）

		// 提交相关（叠加提交专用限流：5 req/s per IP）
		v1.POST("/submit", middleware.AuthRequired(), middleware.SubmitRateLimit(), judgeHandler.SubmitCode)
		v1.GET("/submissions/:id", middleware.AuthRequired(), judgeHandler.GetSubmission)
		v1.POST("/test", middleware.AuthRequired(), middleware.SubmitRateLimit(), judgeHandler.RunTest) // 运行测试

		// 排行榜
		v1.GET("/leaderboard", leaderboardHandler.GetGlobalLeaderboard)
		v1.GET("/problems/:id/leaderboard", leaderboardHandler.GetProblemLeaderboard)

		// 管理员功能
		v1.GET("/admin/submissions", middleware.AuthRequired(), leaderboardHandler.GetAllSubmissions)
		v1.GET("/admin/submissions/:id", middleware.AuthRequired(), leaderboardHandler.GetSubmissionCode)
	}
}
