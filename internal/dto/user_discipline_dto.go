package dto

type (
	UserDisciplineInfo struct {
		Discipline       string `json:"discipline"`
		DisciplineNumber int    `json:"number"`
		Initial          string `json:"initial"`
	}
)
