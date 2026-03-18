package repository

import (
	"github.com/your-org/oj-platform/internal/database"
	"github.com/your-org/oj-platform/internal/models"
)

type SubmissionRepository struct{}

func NewSubmissionRepository() *SubmissionRepository {
	return &SubmissionRepository{}
}

func (r *SubmissionRepository) Create(submission *models.Submission) error {
	return database.DB.Create(submission).Error
}

func (r *SubmissionRepository) GetByID(id uint) (*models.Submission, error) {
	var submission models.Submission
	if err := database.DB.First(&submission, id).Error; err != nil {
		return nil, err
	}
	return &submission, nil
}

func (r *SubmissionRepository) GetByUserID(userID uint, limit, offset int) ([]models.Submission, error) {
	var submissions []models.Submission
	if err := database.DB.Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&submissions).Error; err != nil {
		return nil, err
	}
	return submissions, nil
}

func (r *SubmissionRepository) Update(submission *models.Submission) error {
	return database.DB.Save(submission).Error
}
