package repository

import (
	"context"

	"lms/internal/entity"
)

type GroupRepository interface {
	GetAllGroups(ctx context.Context) ([]entity.Group, error)
}

type CriteriaRepository interface {
	GetCriteriaByGroup(ctx context.Context, groupId int) ([]entity.Criteria, error)
}

type StudentRepository interface {
	GetStudentsByGroup(ctx context.Context, groupId int) ([]entity.Student, error)
}

type LessonRepository interface {
	GetLessonsWithScores(ctx context.Context, groupID int) ([]LessonWithScores, error)
	CreateLesson(ctx context.Context, groupId int) (int, error)
	DeleteLesson(ctx context.Context, groupId, lessonId int) error
	ChangeLessonDates(ctx context.Context, groupId, lessonId int, dates []string) error
	UpdateLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId, score int) (int, error)
	DeleteLessonScore(ctx context.Context, groupId, lessonId, studentId, criteriaId int) (int, error)
}

type RatingRepository interface {
	GetRatings(ctx context.Context, month string) ([]entity.Rating, error)
}

type LessonWithScores struct {
	entity.Lesson
	Students []StudentWithScores `json:"students"`
}

type StudentWithScores struct {
	entity.Student
	Progress int            `json:"progress"`
	Criteria map[string]int `json:"criteria"`
}

type ScoreDetails struct {
	Criteria entity.Criteria `json:"criteria"`
	Score    int             `json:"score"`
}
