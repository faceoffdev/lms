package main

import (
	"lms/internal/database"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"lms/internal/config"
	httpDelivery "lms/internal/delivery/http"
	"lms/internal/repository/postgres"
	"lms/internal/usecase"
)

func main() {
	cfg := config.InitConfig()

	db := database.Connect(cfg.Postgres)

	repo := postgres.NewPostgresRepo(db)

	groupUC := usecase.NewGroupUseCase(repo)
	criteriaUC := usecase.NewCriteriaUseCase(repo)
	studentUC := usecase.NewStudentUseCase(repo)
	lessonUC := usecase.NewLessonUseCase(repo)

	handler := httpDelivery.NewHandler(groupUC, criteriaUC, studentUC, lessonUC)

	mux := httpDelivery.NewHTTPServer(handler)

	log.Printf("Starting server on :%s ...", cfg.App.Port)
	if err := http.ListenAndServe(":"+cfg.App.Port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
