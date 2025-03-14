package usecase

import (
	"context"
	"lms/internal/entity"

	"lms/internal/repository"
)

type RatingUseCase interface {
	GetRatings(ctx context.Context, month string) ([]entity.Rating, error)
}

type ratingUseCase struct {
	ratingRepo repository.RatingRepository
}

func NewRatingUseCase(r repository.RatingRepository) RatingUseCase {
	return &ratingUseCase{ratingRepo: r}
}

func (uc *ratingUseCase) GetRatings(ctx context.Context, month string) ([]entity.Rating, error) {
	return uc.ratingRepo.GetRatings(ctx, month)
}
