package entity

type Rating struct {
	Name        string             `json:"student_name"`
	GroupName   string             `json:"group_name"`
	Performance map[string]float32 `json:"performance"`
}
