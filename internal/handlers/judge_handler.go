package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/oj-platform/internal/judge"
	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/internal/repository"
	"github.com/your-org/oj-platform/internal/services"
	"github.com/your-org/oj-platform/pkg/response"
)

type JudgeHandler struct {
	judgeService *services.JudgeService
	problemRepo  *repository.ProblemRepository
	judge        *judge.Judge
}

func NewJudgeHandler(judgeService *services.JudgeService, problemRepo *repository.ProblemRepository) *JudgeHandler {
	return &JudgeHandler{
		judgeService: judgeService,
		problemRepo:  problemRepo,
	}
}

// SetJudge 设置判题引擎（用于测试功能）
func (h *JudgeHandler) SetJudge(j *judge.Judge) {
	h.judge = j
}

// SubmitCode 提交代码
func (h *JudgeHandler) SubmitCode(c *gin.Context) {
	var req struct {
		ProblemID uint   `json:"problem_id" binding:"required"`
		Code      string `json:"code" binding:"required"`
		Language  string `json:"language" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// 从JWT中获取用户ID
	userID := c.MustGet("user_id").(uint)

	submission, err := h.judgeService.Submit(userID, req.ProblemID, req.Code, req.Language)
	if err != nil {
		response.InternalError(c, "Failed to submit code: "+err.Error())
		return
	}

	response.Success(c, submission)
}

// GetSubmission 获取提交结果
func (h *JudgeHandler) GetSubmission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid submission ID")
		return
	}

	submission, err := h.judgeService.GetSubmission(uint(id))
	if err != nil {
		response.NotFound(c, "Submission not found")
		return
	}

	response.Success(c, submission)
}

// GetProblem 获取题目详情
func (h *JudgeHandler) GetProblem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid problem ID")
		return
	}

	problem, err := h.problemRepo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "Problem not found")
		return
	}

	// 只返回公开的测试用例
	testCases, err := h.problemRepo.GetPublicTestCases(uint(id))
	if err != nil {
		testCases = []models.TestCase{}
	}

	response.Success(c, gin.H{
		"problem":   problem,
		"testCases": testCases,
	})
}

// ListProblems 获取题目列表
func (h *JudgeHandler) ListProblems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	problems, err := h.problemRepo.List(pageSize, (page-1)*pageSize)
	if err != nil {
		response.InternalError(c, "Failed to get problems")
		return
	}

	response.Success(c, gin.H{
		"problems": problems,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateProblem 创建题目（管理员功能）
func (h *JudgeHandler) CreateProblem(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Difficulty  string `json:"difficulty" binding:"required"`
		Tags        string `json:"tags"`
		TimeLimit   int    `json:"time_limit"`
		MemoryLimit int    `json:"memory_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if req.TimeLimit == 0 {
		req.TimeLimit = 5000
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 256
	}

	problem := &models.Problem{
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Tags:        req.Tags,
		TimeLimit:   req.TimeLimit,
		MemoryLimit: req.MemoryLimit,
	}

	if err := h.problemRepo.Create(problem); err != nil {
		response.InternalError(c, "Failed to create problem")
		return
	}

	response.Success(c, problem)
}

// ImportProblem 管理员导入题目（JSON格式，支持100组测试用例）
func (h *JudgeHandler) ImportProblem(c *gin.Context) {
	type TestCaseInput struct {
		Input    string `json:"input"`
		Output   string `json:"output"`
		IsPublic bool   `json:"is_public"`
	}
	var req struct {
		Title            string         `json:"title" binding:"required"`
		Description      string         `json:"description" binding:"required"`
		Difficulty       string         `json:"difficulty" binding:"required"`
		Tags             string         `json:"tags"`
		TimeLimit        int            `json:"time_limit"`
		MemoryLimit      int            `json:"memory_limit"`
		FunctionTemplate string         `json:"function_template"`
		DriverCode       string         `json:"driver_code"`
		TestCases        []TestCaseInput `json:"test_cases"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid JSON: "+err.Error())
		return
	}

	if req.TimeLimit == 0 {
		req.TimeLimit = 5000
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 256
	}
	if len(req.TestCases) > 100 {
		response.BadRequest(c, "最多支持100组测试用例")
		return
	}

	problem := &models.Problem{
		Title:            req.Title,
		Description:      req.Description,
		Difficulty:       req.Difficulty,
		Tags:             req.Tags,
		TimeLimit:        req.TimeLimit,
		MemoryLimit:      req.MemoryLimit,
		FunctionTemplate: req.FunctionTemplate,
		DriverCode:       req.DriverCode,
	}

	if err := h.problemRepo.Create(problem); err != nil {
		response.InternalError(c, "Failed to create problem: "+err.Error())
		return
	}

	for _, tc := range req.TestCases {
		testCase := &models.TestCase{
			ProblemID: problem.ID,
			Input:     tc.Input,
			Output:    tc.Output,
			IsPublic:  tc.IsPublic,
		}
		if err := h.problemRepo.CreateTestCase(testCase); err != nil {
			response.InternalError(c, "Failed to create test case: "+err.Error())
			return
		}
	}

	response.Success(c, gin.H{
		"problem":         problem,
		"test_case_count": len(req.TestCases),
	})
}
func (h *JudgeHandler) RunTest(c *gin.Context) {
	var req struct {
		ProblemID uint   `json:"problem_id" binding:"required"`
		Code      string `json:"code" binding:"required"`
		Language  string `json:"language"`  // 语言，默认 Go
		Input     string `json:"input"` // 用户自定义输入
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if req.Language == "" {
		req.Language = "Go"
	}

	// 获取题目信息
	problem, err := h.problemRepo.GetByID(req.ProblemID)
	if err != nil {
		response.NotFound(c, "Problem not found")
		return
	}

	// 自定义测试：不做输出对比，只运行代码获取输出
	// 传一个特殊标记让 RunGo 不做比较
	result := h.judge.Run(req.Language, req.Code, req.Input, "", problem.DriverCode)

	// 自定义测试下，只要程序正常运行就算成功
	status := result.Status
	if status == "Wrong Answer" {
		status = "Finished"
	}

	response.Success(c, gin.H{
		"status":      status,
		"output":      result.Output,
		"error":       result.Error,
		"time_used":   result.TimeUsed,
		"memory_used": result.MemoryUsed,
	})
}
