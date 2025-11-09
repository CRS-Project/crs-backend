package dto

type (
	AreaOfConcernRequest struct {
		ID              string                             `json:"-"`
		AreaOfConcernId string                             `json:"area_of_concern_id" binding:"required"`
		Description     string                             `json:"description" binding:"required"`
		PackageID       string                             `json:"package_id" binding:"required"`
		Consolidators   []AreaOfConcernConsolidatorRequest `json:"consolidators"`

		AreaOfConcernGroupID string `json:"-"`
		UserId               string `json:"-"`
	}

	UpdateAreaOfConcernRequest struct {
		ID              string                             `json:"-"`
		AreaOfConcernId string                             `json:"area_of_concern_id" binding:"required"`
		Description     string                             `json:"description" binding:"required"`
		PackageID       string                             `json:"package_id" binding:"required"`
		Consolidators   []AreaOfConcernConsolidatorRequest `json:"consolidators"`

		UserId string `json:"-"`
	}

	AreaOfConcernResponse struct {
		ID              string                              `json:"id"`
		AreaOfConcernId string                              `json:"area_of_concern_id"`
		Description     string                              `json:"description"`
		Package         string                              `json:"package"`
		Consolidators   []AreaOfConcernConsolidatorResponse `json:"consolidators,omitempty"`
	}
)
