package dto

type (
	UserDisciplineInfo struct {
		ID               string `json:"id"`
		Discipline       string `json:"discipline"`
		DisciplineNumber int    `json:"number"`
		Initial          string `json:"initial"`
	}
)
