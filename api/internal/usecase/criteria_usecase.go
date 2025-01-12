package usecase

import (
	"context"

	"lms/internal/entity"
	"lms/internal/repository"
)

type CriteriaUseCase interface {
	GetCriteria(ctx context.Context, groupId int) ([]entity.Criteria, error)
}

type criteriaUseCase struct {
	repo repository.CriteriaRepository
}

func NewCriteriaUseCase(r repository.CriteriaRepository) CriteriaUseCase {
	return &criteriaUseCase{repo: r}
}

func (uc *criteriaUseCase) GetCriteria(ctx context.Context, groupId int) ([]entity.Criteria, error) {
	return uc.repo.GetCriteriaByGroup(ctx, groupId)
}
