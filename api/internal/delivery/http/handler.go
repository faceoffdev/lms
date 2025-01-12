package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"lms/internal/usecase"
)

type Handler struct {
	groupUC    usecase.GroupUseCase
	criteriaUC usecase.CriteriaUseCase
	studentsUC usecase.StudentUseCase
	lessonUC   usecase.LessonUseCase
}

func NewHandler(guc usecase.GroupUseCase, cus usecase.CriteriaUseCase, suc usecase.StudentUseCase, luc usecase.LessonUseCase) *Handler {
	return &Handler{groupUC: guc, criteriaUC: cus, studentsUC: suc, lessonUC: luc}
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	data, err := h.groupUC.GetGroups(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(&ResponseData{
		Data: data,
	})
}

func (h *Handler) GetCriteria(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.criteriaUC.GetCriteria(r.Context(), groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(&ResponseData{
		Data: data,
	})
}

func (h *Handler) GetStudents(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.studentsUC.GetStudents(r.Context(), groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(&ResponseData{
		Data: data,
	})
}

func (h *Handler) GetLessons(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.lessonUC.GetLessons(r.Context(), groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(&ResponseData{
		Data: data,
	})
}

func (h *Handler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.lessonUC.CreateLesson(r.Context(), groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(ResponseData{
		Data: createLessonResponse{LessonId: id},
	})
}

func (h *Handler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("group"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lessonId, err := strconv.Atoi(r.PathValue("lesson"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.lessonUC.DeleteLesson(r.Context(), groupId, lessonId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ChangeLessonDates(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("group"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lessonId, err := strconv.Atoi(r.PathValue("lesson"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req changeLessonDatesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.lessonUC.ChangeLessonDates(r.Context(), groupId, lessonId, req.Dates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateLessonScore(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("group"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lessonId, err := strconv.Atoi(r.PathValue("lesson"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req setCriteriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var progress int
	if progress, err = h.lessonUC.UpdateLessonScore(r.Context(), groupId, lessonId, req.StudentId, req.CriteriaId, req.Score); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(ResponseData{
		Data: setLessonScoreResponse{Progress: progress},
	})
}

func (h *Handler) DeleteLessonScore(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(r.PathValue("group"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lessonId, err := strconv.Atoi(r.PathValue("lesson"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req setCriteriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var progress int
	progress, err = h.lessonUC.DeleteLessonScore(r.Context(), groupId, lessonId, req.StudentId, req.CriteriaId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(ResponseData{
		Data: setLessonScoreResponse{Progress: progress},
	})
}
