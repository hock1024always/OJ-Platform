package repository

import (
	"github.com/your-org/oj-platform/internal/database"
	"github.com/your-org/oj-platform/internal/models"
)

type ProblemRepository struct{}

func NewProblemRepository() *ProblemRepository {
	return &ProblemRepository{}
}

func (r *ProblemRepository) GetByID(id uint) (*models.Problem, error) {
	var problem models.Problem
	if err := database.DB.First(&problem, id).Error; err != nil {
		return nil, err
	}
	return &problem, nil
}

func (r *ProblemRepository) GetTestCases(problemID uint) ([]models.TestCase, error) {
	var testCases []models.TestCase
	if err := database.DB.Where("problem_id = ?", problemID).Find(&testCases).Error; err != nil {
		return nil, err
	}
	return testCases, nil
}

func (r *ProblemRepository) GetPublicTestCases(problemID uint) ([]models.TestCase, error) {
	var testCases []models.TestCase
	if err := database.DB.Where("problem_id = ? AND is_public = ?", problemID, true).Find(&testCases).Error; err != nil {
		return nil, err
	}
	return testCases, nil
}

func (r *ProblemRepository) List(limit, offset int) ([]models.Problem, error) {
	var problems []models.Problem
	if err := database.DB.Limit(limit).Offset(offset).Find(&problems).Error; err != nil {
		return nil, err
	}
	return problems, nil
}

func (r *ProblemRepository) Create(problem *models.Problem) error {
	return database.DB.Create(problem).Error
}

func (r *ProblemRepository) CreateTestCase(tc *models.TestCase) error {
	return database.DB.Create(tc).Error
}
