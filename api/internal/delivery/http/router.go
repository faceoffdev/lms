package http

import "net/http"

func NewHTTPServer(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /groups", h.GetGroups)
	mux.HandleFunc("GET /groups/{id}/criteria", h.GetCriteria)
	mux.HandleFunc("GET /groups/{id}/students", h.GetStudents)
	mux.HandleFunc("GET /groups/{id}/lessons", h.GetLessons)
	mux.HandleFunc("POST /groups/{id}/lessons", h.CreateLesson)
	mux.HandleFunc("DELETE /groups/{group}/lessons/{lesson}", h.DeleteLesson)
	mux.HandleFunc("PATCH /groups/{group}/lessons/{lesson}/dates", h.ChangeLessonDates)
	mux.HandleFunc("PATCH /groups/{group}/lessons/{lesson}/criteria", h.UpdateLessonScore)
	mux.HandleFunc("DELETE /groups/{group}/lessons/{lesson}/criteria", h.DeleteLessonScore)
	mux.HandleFunc("GET /ratings", h.GetRatings)

	return JSONMiddleware(corsMiddleware(mux))
}
