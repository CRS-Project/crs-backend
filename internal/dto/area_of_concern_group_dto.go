package dto

type (
	AreaOfConcernGroupRequest struct {
		ID               string `json:"-"`
		ReviewFocus      string `json:"description" binding:"required"`
		UserDisciplineID string `json:"user_discipline_id" binding:"required"`
		PackageID        string `json:"package_id" binding:"required"`
		UserId           string `json:"-"`
	}

	AreaOfConcernGroupResponse struct {
		ID             string `json:"id"`
		ReviewFocus    string `json:"description"`
		Package        string `json:"package"`
		UserDiscipline string `json:"user_discipline"`
	}
)
