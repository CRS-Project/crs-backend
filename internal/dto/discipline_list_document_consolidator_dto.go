package dto

type (
	DisciplineListDocumentConsolidatorRequest struct {
		ID                            string `json:"-"`
		DisciplineGroupConsolidatorID string `json:"discipline_group_consolidator_id" binding:"required"`
		DisciplineListDocumentID      string `json:"discipline_list_document_id" binding:"required"`
	}

	DisciplineListDocumentConsolidatorResponse struct {
		ID                       string `json:"id"`
		DisciplineListDocumentID string `json:"discipline_list_document_id"`
		Name                     string `json:"name"`
	}
)
