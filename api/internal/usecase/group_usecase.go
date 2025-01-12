package usecase

import (
	"context"

	"lms/internal/entity"
	"lms/internal/repository"
)

type GroupUseCase interface {
	GetGroups(ctx context.Context) ([]entity.Group, error)
}

type groupUseCase struct {
	groupRepo repository.GroupRepository
}

func NewGroupUseCase(gr repository.GroupRepository) GroupUseCase {
	return &groupUseCase{groupRepo: gr}
}

func (uc *groupUseCase) GetGroups(ctx context.Context) ([]entity.Group, error) {
	return uc.groupRepo.GetAllGroups(ctx)
}
