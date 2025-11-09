package dto

type (
	AreaOfConcernConsolidatorRequest struct {
		ID     string `json:"-"`
		UserID string `json:"user_id" binding:"required"`
	}

	AreaOfConcernConsolidatorResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)
