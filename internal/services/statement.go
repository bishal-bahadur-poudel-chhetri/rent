package services

import (
	"renting/internal/models"
	"renting/internal/repositories"
)

type StatementService interface {
	GetStatements(filter models.StatementFilter) (*models.PaginatedStatements, error)
}

type statementService struct {
	repo *repositories.StatementRepository
}

func NewStatementService(repo *repositories.StatementRepository) *statementService {
	return &statementService{repo: repo}
}

func (s *statementService) GetStatements(filter models.StatementFilter) (*models.PaginatedStatements, error) {
	// Set default limit if not provided
	if filter.Limit <= 0 {
		filter.Limit = 50
	}

	// Get paginated data
	statements, err := s.repo.GetStatements(filter)
	if err != nil {
		return nil, err
	}

	// Get total count
	total, err := s.repo.GetTotalCount(filter)
	if err != nil {
		return nil, err
	}

	// Prepare response
	response := &models.PaginatedStatements{
		Data:       statements,
		TotalCount: total,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		HasMore:    (filter.Offset + filter.Limit) < total,
	}

	return response, nil
}
