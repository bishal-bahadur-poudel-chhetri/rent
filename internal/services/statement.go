package services

import (
	"context"
	"renting/internal/models"
	"renting/internal/repositories"
)

type StatementService interface {
	GetOutstandingStatements(ctx context.Context, filters map[string]string, offset, limit int) ([]*models.Statement, error)
}

type statementService struct {
	repo repositories.StatementRepository
}

func NewStatementService(repo repositories.StatementRepository) StatementService {
	return &statementService{repo: repo}
}

func (s *statementService) GetOutstandingStatements(ctx context.Context, filters map[string]string, offset, limit int) ([]*models.Statement, error) {
	return s.repo.GetOutstandingStatements(ctx, filters, offset, limit)
}
