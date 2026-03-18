package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/pkg/response"
	"gorm.io/gorm"
)

type LeaderboardHandler struct {
	db *gorm.DB
}

func NewLeaderboardHandler(db *gorm.DB) *LeaderboardHandler {
	return &LeaderboardHandler{db: db}
}

// LeaderboardEntry 排行榜条目
type LeaderboardEntry struct {
	Rank       int    `json:"rank"`
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	ProblemID  uint   `json:"problem_id"`
	ProblemTitle string `json:"problem_title"`
	TimeUsed   int    `json:"time_used"`
	MemoryUsed int    `json:"memory_used"`
	SubmittedAt string `json:"submitted_at"`
}

// UserStats 用户统计
type UserStats struct {
	Rank         int    `json:"rank"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	SolvedCount  int    `json:"solved_count"`
	TotalSubmissions int `json:"total_submissions"`
	AcceptanceRate float64 `json:"acceptance_rate"`
}

// GetProblemLeaderboard 获取某题目的排行榜（最快通过）
func (h *LeaderboardHandler) GetProblemLeaderboard(c *gin.Context) {
	problemIDStr := c.Param("id")
	problemID, err := strconv.ParseUint(problemIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid problem ID")
		return
	}

	// 获取题目信息
	var problem models.Problem
	if err := h.db.First(&problem, problemID).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Problem not found")
		return
	}

	// 查询该题目通过的最快提交（每个用户只取最快的一次）
	type Result struct {
		UserID     uint
		Username   string
		TimeUsed   int
		MemoryUsed int
		CreatedAt  string
	}

	var results []Result
	h.db.Raw(`
		SELECT s.user_id, u.username, s.time_used, s.memory_used, s.created_at
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		WHERE s.problem_id = ? AND s.status = 'Accepted'
		GROUP BY s.user_id
		HAVING MIN(s.time_used)
		ORDER BY s.time_used ASC
		LIMIT 50
	`, problemID).Scan(&results)

	// 构建排行榜
	entries := make([]LeaderboardEntry, len(results))
	for i, r := range results {
		entries[i] = LeaderboardEntry{
			Rank:         i + 1,
			UserID:       r.UserID,
			Username:     r.Username,
			ProblemID:    uint(problemID),
			ProblemTitle: problem.Title,
			TimeUsed:     r.TimeUsed,
			MemoryUsed:   r.MemoryUsed,
			SubmittedAt:  r.CreatedAt,
		}
	}

	response.Success(c, gin.H{
		"problem": gin.H{
			"id":    problem.ID,
			"title": problem.Title,
		},
		"leaderboard": entries,
	})
}

// GetGlobalLeaderboard 获取全局排行榜（按解题数）
func (h *LeaderboardHandler) GetGlobalLeaderboard(c *gin.Context) {
	type Result struct {
		UserID          uint
		Username        string
		SolvedCount     int
		TotalSubmissions int
	}

	var results []Result
	h.db.Raw(`
		SELECT 
			u.id as user_id,
			u.username,
			COUNT(DISTINCT s.problem_id) as solved_count,
			(SELECT COUNT(*) FROM submissions WHERE user_id = u.id) as total_submissions
		FROM users u
		LEFT JOIN submissions s ON u.id = s.user_id AND s.status = 'Accepted'
		GROUP BY u.id
		HAVING solved_count > 0
		ORDER BY solved_count DESC, total_submissions ASC
		LIMIT 100
	`).Scan(&results)

	entries := make([]UserStats, len(results))
	for i, r := range results {
		acceptanceRate := 0.0
		if r.TotalSubmissions > 0 {
			acceptanceRate = float64(r.SolvedCount) / float64(r.TotalSubmissions) * 100
		}
		entries[i] = UserStats{
			Rank:           i + 1,
			UserID:         r.UserID,
			Username:       r.Username,
			SolvedCount:    r.SolvedCount,
			TotalSubmissions: r.TotalSubmissions,
			AcceptanceRate: acceptanceRate,
		}
	}

	response.Success(c, gin.H{
		"leaderboard": entries,
	})
}

// GetAllSubmissions 获取所有提交记录（管理员用）
func (h *LeaderboardHandler) GetAllSubmissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	problemID := c.Query("problem_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := h.db.Model(&models.Submission{}).Preload("User").Preload("Problem")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if problemID != "" {
		query = query.Where("problem_id = ?", problemID)
	}

	var total int64
	query.Count(&total)

	var submissions []models.Submission
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&submissions)

	// 构建返回数据
	type SubmissionResponse struct {
		ID         uint   `json:"id"`
		UserID     uint   `json:"user_id"`
		Username   string `json:"username"`
		ProblemID  uint   `json:"problem_id"`
		ProblemTitle string `json:"problem_title"`
		Status     string `json:"status"`
		TimeUsed   int    `json:"time_used"`
		MemoryUsed int    `json:"memory_used"`
		CreatedAt  string `json:"created_at"`
	}

	results := make([]SubmissionResponse, len(submissions))
	for i, s := range submissions {
		username := ""
		if s.User != nil {
			username = s.User.Username
		}
		problemTitle := ""
		if s.Problem != nil {
			problemTitle = s.Problem.Title
		}
		results[i] = SubmissionResponse{
			ID:           s.ID,
			UserID:       s.UserID,
			Username:     username,
			ProblemID:    s.ProblemID,
			ProblemTitle: problemTitle,
			Status:       s.Status,
			TimeUsed:     s.TimeUsed,
			MemoryUsed:   s.MemoryUsed,
			CreatedAt:    s.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	response.Success(c, gin.H{
		"submissions": results,
		"total":       total,
		"page":        page,
		"pageSize":    pageSize,
	})
}

// GetSubmissionCode 获取提交的代码（管理员用）
func (h *LeaderboardHandler) GetSubmissionCode(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid submission ID")
		return
	}

	var submission models.Submission
	if err := h.db.Preload("User").Preload("Problem").First(&submission, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Submission not found")
		return
	}

	response.Success(c, gin.H{
		"submission": gin.H{
			"id":          submission.ID,
			"user_id":     submission.UserID,
			"username":    submission.User.Username,
			"problem_id":  submission.ProblemID,
			"problem_title": submission.Problem.Title,
			"code":        submission.Code,
			"language":    submission.Language,
			"status":      submission.Status,
			"time_used":   submission.TimeUsed,
			"memory_used": submission.MemoryUsed,
			"result":      submission.Result,
			"created_at":  submission.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
