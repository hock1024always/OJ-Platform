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

func (r *ProblemRepository) Update(problem *models.Problem) error {
	return database.DB.Save(problem).Error
}

func (r *ProblemRepository) Delete(id uint) error {
	// 先删测试用例再删题目（软删除）
	database.DB.Where("problem_id = ?", id).Delete(&models.TestCase{})
	return database.DB.Delete(&models.Problem{}, id).Error
}

func (r *ProblemRepository) GetTestCaseByID(id uint) (*models.TestCase, error) {
	var tc models.TestCase
	if err := database.DB.First(&tc, id).Error; err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *ProblemRepository) UpdateTestCase(tc *models.TestCase) error {
	return database.DB.Save(tc).Error
}

func (r *ProblemRepository) DeleteTestCase(id uint) error {
	return database.DB.Delete(&models.TestCase{}, id).Error
}

func (r *ProblemRepository) DeleteTestCasesByProblemID(problemID uint) error {
	return database.DB.Where("problem_id = ?", problemID).Delete(&models.TestCase{}).Error
}

func (r *ProblemRepository) CreateTestCases(tcs []models.TestCase) error {
	if len(tcs) == 0 {
		return nil
	}
	return database.DB.Create(&tcs).Error
}

func (r *ProblemRepository) CountTestCases(problemID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.TestCase{}).Where("problem_id = ?", problemID).Count(&count).Error
	return count, err
}
