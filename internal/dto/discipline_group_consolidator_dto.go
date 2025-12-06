package dto

type (
	DisciplineGroupConsolidatorRequest struct {
		ID     string `json:"-"`
		UserID string `json:"user_id" binding:"required"`
	}

	DisciplineGroupConsolidatorResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)
