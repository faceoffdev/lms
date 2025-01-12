package usecase

import (
	"context"

	"lms/internal/repository"
)

type LessonUseCase interface {
	GetLessons(ctx context.Context, groupID int) ([]repository.LessonWithScores, error)
	CreateLesson(ctx context.Context, groupId int) (int, error)
	DeleteLesson(ctx context.Context, groupId, lessonId int) error
	ChangeLessonDates(ctx context.Context, groupId, lessonId int, dates []string) error
	UpdateLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId, score int) (int, error)
	DeleteLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId int) (int, error)
}

type lessonUseCase struct {
	lessonRepo repository.LessonRepository
}

func NewLessonUseCase(lr repository.LessonRepository) LessonUseCase {
	return &lessonUseCase{lessonRepo: lr}
}

func (uc *lessonUseCase) GetLessons(ctx context.Context, groupID int) ([]repository.LessonWithScores, error) {
	return uc.lessonRepo.GetLessonsWithScores(ctx, groupID)
}

func (uc *lessonUseCase) CreateLesson(ctx context.Context, groupId int) (int, error) {
	return uc.lessonRepo.CreateLesson(ctx, groupId)
}

func (uc *lessonUseCase) DeleteLesson(ctx context.Context, groupId, lessonId int) error {
	return uc.lessonRepo.DeleteLesson(ctx, groupId, lessonId)
}

func (uc *lessonUseCase) ChangeLessonDates(ctx context.Context, groupId, lessonId int, dates []string) error {
	return uc.lessonRepo.ChangeLessonDates(ctx, groupId, lessonId, dates)
}

func (uc *lessonUseCase) UpdateLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId int, score int) (int, error) {
	return uc.lessonRepo.UpdateLessonScore(ctx, groupId, lessonId, studentId, criteriaId, score)
}

func (uc *lessonUseCase) DeleteLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId int) (int, error) {
	return uc.lessonRepo.DeleteLessonScore(ctx, groupId, lessonId, studentId, criteriaId)
}
