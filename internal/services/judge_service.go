package services

import (
	"fmt"
	"log"

	"github.com/your-org/oj-platform/internal/judge"
	"github.com/your-org/oj-platform/internal/models"
	"github.com/your-org/oj-platform/internal/queue"
	"github.com/your-org/oj-platform/internal/repository"
)

type JudgeService struct {
	judge      *judge.Judge
	problemRepo *repository.ProblemRepository
	submitRepo *repository.SubmissionRepository
	queue      *queue.TaskQueue
}

func NewJudgeService(judgeInstance *judge.Judge, problemRepo *repository.ProblemRepository, submitRepo *repository.SubmissionRepository, taskQueue *queue.TaskQueue) *JudgeService {
	service := &JudgeService{
		judge:      judgeInstance,
		problemRepo: problemRepo,
		submitRepo: submitRepo,
		queue:      taskQueue,
	}

	// 启动队列worker
	taskQueue.Start(service.handleTask)

	return service
}

// Submit 提交代码
func (s *JudgeService) Submit(userID, problemID uint, code, language string) (*models.Submission, error) {
	// 创建提交记录
	submission := &models.Submission{
		UserID:    userID,
		ProblemID: problemID,
		Code:      code,
		Language:  language,
		Status:    "Pending",
	}

	if err := s.submitRepo.Create(submission); err != nil {
		return nil, err
	}

	// 创建任务
	task := &queue.Task{
		ID:        fmt.Sprintf("%d", submission.ID),
		ProblemID: problemID,
		UserID:    userID,
		Code:      code,
		Language:  language,
	}

	// 异步执行
	go func() {
		resultChan := make(chan *queue.TaskResult, 1)
		task.ResultChan = resultChan
		s.queue.Submit(task)

		result := <-resultChan
		s.updateSubmission(submission.ID, result)
	}()

	return submission, nil
}

// handleTask 处理任务（一次编译，多次执行）
func (s *JudgeService) handleTask(task *queue.Task) *queue.TaskResult {
	// 获取题目和测试用例
	problem, err := s.problemRepo.GetByID(task.ProblemID)
	if err != nil {
		return &queue.TaskResult{
			TaskID: task.ID,
			Status: "System Error",
			Error:  "Problem not found",
		}
	}

	testCases, err := s.problemRepo.GetTestCases(task.ProblemID)
	if err != nil || len(testCases) == 0 {
		return &queue.TaskResult{
			TaskID: task.ID,
			Status: "System Error",
			Error:  "No test cases found",
		}
	}

	// 一次编译
	prog, compileResult := s.judge.Compile(task.Language, task.Code, problem.DriverCode)
	if compileResult != nil {
		return &queue.TaskResult{
			TaskID:     task.ID,
			Status:     compileResult.Status,
			Error:      compileResult.Error,
			TimeUsed:   compileResult.TimeUsed,
			MemoryUsed: compileResult.MemoryUsed,
		}
	}
	defer prog.Cleanup()

	// 多次执行测试用例
	totalTime := 0
	maxMemory := 0

	for _, tc := range testCases {
		result := s.judge.RunCompiled(prog, tc.Input, tc.Output)

		if result.Status != "Accepted" {
			return &queue.TaskResult{
				TaskID:     task.ID,
				Status:     result.Status,
				Output:     result.Output,
				Expected:   result.Expected,
				Error:      result.Error,
				TimeUsed:   result.TimeUsed,
				MemoryUsed: result.MemoryUsed,
			}
		}

		totalTime += result.TimeUsed
		if result.MemoryUsed > maxMemory {
			maxMemory = result.MemoryUsed
		}
	}

	// 所有测试用例通过
	return &queue.TaskResult{
		TaskID:     task.ID,
		Status:     "Accepted",
		TimeUsed:   totalTime,
		MemoryUsed: maxMemory,
	}
}

// updateSubmission 更新提交记录
func (s *JudgeService) updateSubmission(submissionID uint, result *queue.TaskResult) {
	submission, err := s.submitRepo.GetByID(submissionID)
	if err != nil {
		log.Printf("Failed to get submission %d: %v", submissionID, err)
		return
	}

	submission.Status = result.Status
	submission.TimeUsed = result.TimeUsed
	submission.MemoryUsed = result.MemoryUsed
	submission.Result = result.Error

	if result.Output != "" {
		submission.Result = fmt.Sprintf("Output: %s", result.Output)
		if result.Expected != "" {
			submission.Result = fmt.Sprintf("%s\nExpected: %s", submission.Result, result.Expected)
		}
	}

	if err := s.submitRepo.Update(submission); err != nil {
		log.Printf("Failed to update submission %d: %v", submissionID, err)
	}
}

// GetSubmission 获取提交记录
func (s *JudgeService) GetSubmission(id uint) (*models.Submission, error) {
	return s.submitRepo.GetByID(id)
}
