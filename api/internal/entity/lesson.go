package entity

type Lesson struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
