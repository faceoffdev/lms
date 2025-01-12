package entity

type LessonScore struct {
	LessonID   int `json:"lesson_id"`
	StudentID  int `json:"student_id"`
	CriteriaID int `json:"criteria_id"`
	Score      int `json:"score"`
}
