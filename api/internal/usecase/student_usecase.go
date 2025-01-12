package usecase

import (
	"context"

	"lms/internal/entity"
	"lms/internal/repository"
)

type StudentUseCase interface {
	GetStudents(ctx context.Context, groupId int) ([]entity.Student, error)
}

type studentUseCase struct {
	repo repository.StudentRepository
}

func NewStudentUseCase(r repository.StudentRepository) StudentUseCase {
	return &studentUseCase{repo: r}
}

func (uc *studentUseCase) GetStudents(ctx context.Context, groupId int) ([]entity.Student, error) {
	return uc.repo.GetStudentsByGroup(ctx, groupId)
}
