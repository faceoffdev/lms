package entity

type Rating struct {
	Name        string  `json:"student_name"`
	GroupName   string  `json:"group_name"`
	Month       string  `json:"month"`
	Performance float32 `json:"performance_percentage"`
}
