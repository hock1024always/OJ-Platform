package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/oj-platform/internal/handlers"
	"github.com/your-org/oj-platform/internal/middleware"
)

func Setup(r *gin.Engine, judgeHandler *handlers.JudgeHandler, userHandler *handlers.UserHandler, leaderboardHandler *handlers.LeaderboardHandler, problemAdminHandler *handlers.ProblemAdminHandler) {
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

		// 管理员功能（需要认证）
		admin := v1.Group("/admin", middleware.AuthRequired())
		{
			admin.GET("/submissions", leaderboardHandler.GetAllSubmissions)
			admin.GET("/submissions/:id", leaderboardHandler.GetSubmissionCode)

			// 代码生成预览
			admin.POST("/codegen/preview", problemAdminHandler.GenerateCode)

			// 题目 CRUD
			admin.POST("/problems", problemAdminHandler.CreateProblemFull)
			admin.PUT("/problems/:id", problemAdminHandler.UpdateProblem)
			admin.DELETE("/problems/:id", problemAdminHandler.DeleteProblem)
			admin.PUT("/problems/:id/status", problemAdminHandler.ToggleProblemStatus)
			admin.GET("/problems/:id/templates", problemAdminHandler.GetProblemTemplates)

			// 测试用例管理
			admin.GET("/problems/:id/testcases", problemAdminHandler.ListTestCases)
			admin.POST("/problems/:id/testcases", problemAdminHandler.AddTestCase)
			admin.PUT("/testcases/:tc_id", problemAdminHandler.UpdateTestCase)
			admin.DELETE("/testcases/:tc_id", problemAdminHandler.DeleteTestCase)

			// 标程自动生成 & 验证测试用例
			admin.POST("/problems/:id/generate-testcases", problemAdminHandler.GenerateTestCases)
			admin.POST("/problems/:id/validate-testcases", problemAdminHandler.ValidateTestCases)
		}
	}
}
