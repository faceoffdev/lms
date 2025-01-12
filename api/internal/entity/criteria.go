package entity

type Criteria struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	MaxScore int    `json:"max_score"`
}
