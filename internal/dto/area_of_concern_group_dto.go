package dto

type (
	AreaOfConcernGroupRequest struct {
		ID               string `json:"-"`
		ReviewFocus      string `json:"review_focus" binding:"required"`
		UserDisciplineID string `json:"user_discipline_id" binding:"required"`
		PackageID        string `json:"package_id" binding:"required"`
		UserId           string `json:"-"`
	}

	AreaOfConcernGroupResponse struct {
		ID             string `json:"id"`
		ReviewFocus    string `json:"review_focus"`
		Package        string `json:"package"`
		UserDiscipline string `json:"user_discipline"`
	}

	AreaOfConcernGroupStatistic struct {
		TotalAreaOfConcernGroup int `json:"total_area_of_concern_group"`
		TotalAreaOfConcern      int `json:"total_area_of_concern"`
		TotalComment            int `json:"total_comment"`
	}
)
