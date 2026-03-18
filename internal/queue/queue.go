package queue

import (
	"sync"
)

// Task 表示一个判题任务
type Task struct {
	ID          string
	ProblemID   uint
	UserID      uint
	Code        string
	Language    string
	ResultChan  chan *TaskResult
}

// TaskResult 表示任务执行结果
type TaskResult struct {
	TaskID      string
	Status      string // "Accepted", "Wrong Answer", "Compile Error", "Time Limit Exceeded", etc
	Output      string
	Expected    string
	TimeUsed    int // 毫秒
	MemoryUsed  int // KB
	Error       string
}

// TaskQueue 任务队列
type TaskQueue struct {
	tasks  chan *Task
	workers int
	wg     sync.WaitGroup
}

// NewTaskQueue 创建任务队列
func NewTaskQueue(maxSize, workers int) *TaskQueue {
	return &TaskQueue{
		tasks:   make(chan *Task, maxSize),
		workers: workers,
	}
}

// Submit 提交任务
func (q *TaskQueue) Submit(task *Task) error {
	q.tasks <- task
	return nil
}

// Start 启动worker
func (q *TaskQueue) Start(handler func(*Task) *TaskResult) {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()
			for task := range q.tasks {
				result := handler(task)
				if task.ResultChan != nil {
					task.ResultChan <- result
				}
			}
		}()
	}
}

// Stop 停止队列
func (q *TaskQueue) Stop() {
	close(q.tasks)
	q.wg.Wait()
}
