package handlers

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-org/oj-platform/internal/codegen"
	"github.com/your-org/oj-platform/internal/judge"
	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/internal/repository"
	"github.com/your-org/oj-platform/pkg/response"
)

type ProblemAdminHandler struct {
	problemRepo *repository.ProblemRepository
	judge       *judge.Judge
}

func NewProblemAdminHandler(problemRepo *repository.ProblemRepository, j *judge.Judge) *ProblemAdminHandler {
	return &ProblemAdminHandler{
		problemRepo: problemRepo,
		judge:       j,
	}
}

// GenerateCode 根据函数签名预览生成的多语言代码
func (h *ProblemAdminHandler) GenerateCode(c *gin.Context) {
	var req struct {
		FunctionSignature codegen.FunctionSignature `json:"function_signature" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	codes, err := codegen.GenerateAll(&req.FunctionSignature)
	if err != nil {
		response.InternalError(c, "代码生成失败: "+err.Error())
		return
	}

	response.Success(c, codes)
}

// CreateProblemFull 完整创建题目（带函数签名 + 自动生成多语言代码）
func (h *ProblemAdminHandler) CreateProblemFull(c *gin.Context) {
	var req struct {
		Title             string                     `json:"title" binding:"required"`
		Description       string                     `json:"description" binding:"required"`
		Difficulty        string                     `json:"difficulty" binding:"required"`
		Tags              string                     `json:"tags"`
		TimeLimit         int                        `json:"time_limit"`
		MemoryLimit       int                        `json:"memory_limit"`
		FunctionSignature *codegen.FunctionSignature `json:"function_signature"`
		// 如果不用结构化签名，也可以直接提供 Go 的模板和驱动代码
		FunctionTemplate string `json:"function_template"`
		DriverCode       string `json:"driver_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
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
		Status:      "active",
	}

	if req.FunctionSignature != nil {
		// 用结构化签名生成多语言代码
		sigJSON, _ := json.Marshal(req.FunctionSignature)
		problem.FunctionSignature = string(sigJSON)

		codes, err := codegen.GenerateAll(req.FunctionSignature)
		if err != nil {
			response.InternalError(c, "代码生成失败: "+err.Error())
			return
		}

		// Go 的模板和驱动代码放主字段
		if goCode, ok := codes["Go"]; ok {
			problem.FunctionTemplate = goCode.FunctionTemplate
			problem.DriverCode = goCode.DriverCode
		}

		// 所有语言的模板存 TemplatesJSON
		templatesMap := make(map[string]map[string]string)
		for lang, code := range codes {
			templatesMap[lang] = map[string]string{
				"function_template": code.FunctionTemplate,
				"driver_code":       code.DriverCode,
			}
		}
		templatesJSON, _ := json.Marshal(templatesMap)
		problem.TemplatesJSON = string(templatesJSON)
	} else {
		// 直接使用传入的 Go 模板
		problem.FunctionTemplate = req.FunctionTemplate
		problem.DriverCode = req.DriverCode
	}

	if err := h.problemRepo.Create(problem); err != nil {
		response.InternalError(c, "创建题目失败: "+err.Error())
		return
	}

	response.Success(c, problem)
}

// UpdateProblem 编辑题目
func (h *ProblemAdminHandler) UpdateProblem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	problem, err := h.problemRepo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "题目不存在")
		return
	}

	var req struct {
		Title             string                     `json:"title"`
		Description       string                     `json:"description"`
		Difficulty        string                     `json:"difficulty"`
		Tags              string                     `json:"tags"`
		TimeLimit         int                        `json:"time_limit"`
		MemoryLimit       int                        `json:"memory_limit"`
		FunctionSignature *codegen.FunctionSignature `json:"function_signature"`
		FunctionTemplate  string                     `json:"function_template"`
		DriverCode        string                     `json:"driver_code"`
		Status            string                     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Title != "" {
		problem.Title = req.Title
	}
	if req.Description != "" {
		problem.Description = req.Description
	}
	if req.Difficulty != "" {
		problem.Difficulty = req.Difficulty
	}
	if req.Tags != "" {
		problem.Tags = req.Tags
	}
	if req.TimeLimit > 0 {
		problem.TimeLimit = req.TimeLimit
	}
	if req.MemoryLimit > 0 {
		problem.MemoryLimit = req.MemoryLimit
	}
	if req.Status != "" {
		problem.Status = req.Status
	}

	if req.FunctionSignature != nil {
		sigJSON, _ := json.Marshal(req.FunctionSignature)
		problem.FunctionSignature = string(sigJSON)

		codes, err := codegen.GenerateAll(req.FunctionSignature)
		if err != nil {
			response.InternalError(c, "代码生成失败: "+err.Error())
			return
		}
		if goCode, ok := codes["Go"]; ok {
			problem.FunctionTemplate = goCode.FunctionTemplate
			problem.DriverCode = goCode.DriverCode
		}
		templatesMap := make(map[string]map[string]string)
		for lang, code := range codes {
			templatesMap[lang] = map[string]string{
				"function_template": code.FunctionTemplate,
				"driver_code":       code.DriverCode,
			}
		}
		templatesJSON, _ := json.Marshal(templatesMap)
		problem.TemplatesJSON = string(templatesJSON)
	} else {
		if req.FunctionTemplate != "" {
			problem.FunctionTemplate = req.FunctionTemplate
		}
		if req.DriverCode != "" {
			problem.DriverCode = req.DriverCode
		}
	}

	if err := h.problemRepo.Update(problem); err != nil {
		response.InternalError(c, "更新题目失败: "+err.Error())
		return
	}

	response.Success(c, problem)
}

// DeleteProblem 删除题目
func (h *ProblemAdminHandler) DeleteProblem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	if err := h.problemRepo.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除题目失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "题目已删除"})
}

// ToggleProblemStatus 切换题目上下架
func (h *ProblemAdminHandler) ToggleProblemStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	problem, err := h.problemRepo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "题目不存在")
		return
	}

	if problem.Status == "active" {
		problem.Status = "inactive"
	} else {
		problem.Status = "active"
	}

	if err := h.problemRepo.Update(problem); err != nil {
		response.InternalError(c, "切换状态失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": problem.ID, "status": problem.Status})
}

// GetProblemTemplates 获取题目多语言模板
func (h *ProblemAdminHandler) GetProblemTemplates(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	problem, err := h.problemRepo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "题目不存在")
		return
	}

	result := make(map[string]map[string]string)

	// 优先使用 TemplatesJSON
	if problem.TemplatesJSON != "" {
		json.Unmarshal([]byte(problem.TemplatesJSON), &result)
	}

	// 确保至少有 Go 的模板
	if _, ok := result["Go"]; !ok {
		result["Go"] = map[string]string{
			"function_template": problem.FunctionTemplate,
			"driver_code":       problem.DriverCode,
		}
	}

	response.Success(c, result)
}

// ListTestCases 获取题目所有测试用例
func (h *ProblemAdminHandler) ListTestCases(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	testCases, err := h.problemRepo.GetTestCases(uint(id))
	if err != nil {
		response.InternalError(c, "获取测试用例失败: "+err.Error())
		return
	}

	response.Success(c, testCases)
}

// AddTestCase 添加测试用例
func (h *ProblemAdminHandler) AddTestCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	var req struct {
		Input    string `json:"input" binding:"required"`
		Output   string `json:"output" binding:"required"`
		IsPublic bool   `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	tc := &models.TestCase{
		ProblemID: uint(id),
		Input:     req.Input,
		Output:    req.Output,
		IsPublic:  req.IsPublic,
	}
	if err := h.problemRepo.CreateTestCase(tc); err != nil {
		response.InternalError(c, "添加测试用例失败: "+err.Error())
		return
	}

	response.Success(c, tc)
}

// UpdateTestCase 更新测试用例
func (h *ProblemAdminHandler) UpdateTestCase(c *gin.Context) {
	tcID, err := strconv.ParseUint(c.Param("tc_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的测试用例 ID")
		return
	}

	tc, err := h.problemRepo.GetTestCaseByID(uint(tcID))
	if err != nil {
		response.NotFound(c, "测试用例不存在")
		return
	}

	var req struct {
		Input    string `json:"input"`
		Output   string `json:"output"`
		IsPublic *bool  `json:"is_public"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Input != "" {
		tc.Input = req.Input
	}
	if req.Output != "" {
		tc.Output = req.Output
	}
	if req.IsPublic != nil {
		tc.IsPublic = *req.IsPublic
	}

	if err := h.problemRepo.UpdateTestCase(tc); err != nil {
		response.InternalError(c, "更新测试用例失败: "+err.Error())
		return
	}

	response.Success(c, tc)
}

// DeleteTestCase 删除测试用例
func (h *ProblemAdminHandler) DeleteTestCase(c *gin.Context) {
	tcID, err := strconv.ParseUint(c.Param("tc_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的测试用例 ID")
		return
	}

	if err := h.problemRepo.DeleteTestCase(uint(tcID)); err != nil {
		response.InternalError(c, "删除测试用例失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "测试用例已删除"})
}

// GenerateTestCases 基于标程自动生成测试用例
func (h *ProblemAdminHandler) GenerateTestCases(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	problem, err := h.problemRepo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "题目不存在")
		return
	}

	var req struct {
		ReferenceCode string                   `json:"reference_code" binding:"required"`
		Language      string                   `json:"language"`
		Count         int                      `json:"count"`
		Constraints   []codegen.InputConstraint `json:"constraints"`
		PublicCount   int                      `json:"public_count"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Language == "" {
		req.Language = "Go"
	}
	if req.Count <= 0 {
		req.Count = 50
	}
	if req.Count > 200 {
		req.Count = 200
	}
	if req.PublicCount <= 0 {
		req.PublicCount = 3
	}

	// 解析函数签名
	if problem.FunctionSignature == "" {
		response.BadRequest(c, "题目缺少函数签名定义，无法自动生成测试数据")
		return
	}
	sig, err := codegen.ParseSignature(problem.FunctionSignature)
	if err != nil {
		response.InternalError(c, "解析函数签名失败: "+err.Error())
		return
	}

	// 确定驱动代码
	driverCode := problem.DriverCode
	if req.Language != "Go" && problem.TemplatesJSON != "" {
		var templates map[string]map[string]string
		json.Unmarshal([]byte(problem.TemplatesJSON), &templates)
		if t, ok := templates[req.Language]; ok {
			driverCode = t["driver_code"]
		}
	}

	// 编译标程
	prog, compileResult := h.judge.Compile(req.Language, req.ReferenceCode, driverCode)
	if compileResult != nil {
		response.BadRequest(c, "标程编译失败: "+compileResult.Error)
		return
	}
	defer prog.Cleanup()

	// 生成测试用例
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var testCases []models.TestCase
	var failedInputs []string

	for i := 0; i < req.Count; i++ {
		input, err := codegen.GenerateRandomInput(sig, req.Constraints, rng)
		if err != nil {
			continue
		}

		// 运行标程获取输出
		result := h.judge.RunCompiled(prog, input, "")
		if result.Status != "Accepted" && result.Status != "Finished" {
			failedInputs = append(failedInputs, input)
			continue
		}

		tc := models.TestCase{
			ProblemID: uint(id),
			Input:     input,
			Output:    result.Output,
			IsPublic:  i < req.PublicCount,
		}
		testCases = append(testCases, tc)
	}

	// 批量写入
	if len(testCases) > 0 {
		if err := h.problemRepo.CreateTestCases(testCases); err != nil {
			response.InternalError(c, "保存测试用例失败: "+err.Error())
			return
		}
	}

	response.Success(c, gin.H{
		"generated": len(testCases),
		"failed":    len(failedInputs),
		"total":     req.Count,
	})
}

// ValidateTestCases 用标程验证已有测试用例
func (h *ProblemAdminHandler) ValidateTestCases(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的题目 ID")
		return
	}

	problem, err := h.problemRepo.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "题目不存在")
		return
	}

	var req struct {
		ReferenceCode string `json:"reference_code" binding:"required"`
		Language      string `json:"language"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if req.Language == "" {
		req.Language = "Go"
	}

	driverCode := problem.DriverCode
	if req.Language != "Go" && problem.TemplatesJSON != "" {
		var templates map[string]map[string]string
		json.Unmarshal([]byte(problem.TemplatesJSON), &templates)
		if t, ok := templates[req.Language]; ok {
			driverCode = t["driver_code"]
		}
	}

	// 编译标程
	prog, compileResult := h.judge.Compile(req.Language, req.ReferenceCode, driverCode)
	if compileResult != nil {
		response.BadRequest(c, "标程编译失败: "+compileResult.Error)
		return
	}
	defer prog.Cleanup()

	// 获取所有测试用例
	testCases, err := h.problemRepo.GetTestCases(uint(id))
	if err != nil {
		response.InternalError(c, "获取测试用例失败: "+err.Error())
		return
	}

	passed := 0
	failed := 0
	type FailDetail struct {
		TestCaseID uint   `json:"test_case_id"`
		Input      string `json:"input"`
		Expected   string `json:"expected"`
		Actual     string `json:"actual"`
		Status     string `json:"status"`
	}
	var failures []FailDetail

	for _, tc := range testCases {
		result := h.judge.RunCompiled(prog, tc.Input, tc.Output)
		if result.Status == "Accepted" {
			passed++
		} else {
			failed++
			failures = append(failures, FailDetail{
				TestCaseID: tc.ID,
				Input:      tc.Input,
				Expected:   tc.Output,
				Actual:     result.Output,
				Status:     result.Status,
			})
		}
	}

	response.Success(c, gin.H{
		"total":    len(testCases),
		"passed":   passed,
		"failed":   failed,
		"failures": failures,
	})
}
