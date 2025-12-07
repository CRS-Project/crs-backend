package dto

type (
	DisciplineGroupConsolidatorRequest struct {
		ID     string `json:"-"`
		UserID string `json:"user_id" binding:"required"`
	}

	DisciplineGroupConsolidatorResponse struct {
		ID                            string `json:"id"`
		DisciplineGroupConsolidatorID string `json:"discipline_group_consolidator_id"`
		Name                          string `json:"name"`
	}
)
