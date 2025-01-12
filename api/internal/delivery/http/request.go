package http

type changeLessonDatesRequest struct {
	Dates []string `json:"dates"`
}

type setCriteriaRequest struct {
	StudentId  int `json:"student_id"`
	CriteriaId int `json:"criteria_id"`
	Score      int `json:"score"`
}
