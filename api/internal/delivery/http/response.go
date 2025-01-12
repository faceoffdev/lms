package http

type ResponseData struct {
	Data interface{} `json:"data"`
}

type createLessonResponse struct {
	LessonId int `json:"lesson_id"`
}

type setLessonScoreResponse struct {
	Progress int `json:"progress"`
}
