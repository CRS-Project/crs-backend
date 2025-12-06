package dto

type (
	DisciplineGroupRequest struct {
		ID                           string                               `json:"-"`
		ReviewFocus                  string                               `json:"review_focus" binding:"required"`
		UserDiscipline               string                               `json:"user_discipline" binding:"required"`
		DisciplineInitial            string                               `json:"discipline_initial" binding:""`
		PackageID                    string                               `json:"package_id" binding:"required"`
		DisciplineGroupConsolidators []DisciplineGroupConsolidatorRequest `json:"discipline_group_consolidators"`
		UserId                       string                               `json:"-"`
	}

	DisciplineGroupResponse struct {
		ID                           string                                `json:"id"`
		ReviewFocus                  string                                `json:"review_focus"`
		Package                      string                                `json:"package"`
		UserDiscipline               string                                `json:"user_discipline"`
		DisciplineInitial            string                                `json:"discipline_initial"`
		DisciplineGroupConsolidators []DisciplineGroupConsolidatorResponse `json:"discipline_group_consolidators,omitempty"`
	}

	DisciplineGroupStatistic struct {
		TotalDisciplineGroup        int `json:"total_discipline_group"`
		TotalDisciplineListDocument int `json:"total_discipline_list_document"`
		TotalComment                int `json:"total_comment"`
	}
)
