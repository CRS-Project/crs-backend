package dto

type (
	UserDisciplineInfo struct {
		Number  int     `json:"number"`
		Initial string  `json:"initial"`
		Package *string `json:"package"`
	}
)
