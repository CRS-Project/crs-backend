package dto

type (
	DisciplineListDocumentRequest struct {
		ID                string                                      `json:"-"`
		DocumentID        string                                      `json:"document_id" binding:"required"`
		PackageID         string                                      `json:"package_id" binding:"required"`
		Consolidators     []DisciplineListDocumentConsolidatorRequest `json:"consolidators"`
		DisciplineGroupID string                                      `json:"-"`
		UserId            string                                      `json:"-"`
	}

	UpdateDisciplineListDocumentRequest struct {
		ID                string                                      `json:"-"`
		DocumentID        string                                      `json:"document_id" binding:"required"`
		Consolidators     []DisciplineListDocumentConsolidatorRequest `json:"consolidators"`
		DisciplineGroupID string                                      `json:"-"`
		UserId            string                                      `json:"-"`
	}

	DisciplineListDocumentResponse struct {
		ID            string                                       `json:"id"`
		Package       string                                       `json:"package"`
		Document      *DocumentDetailResponse                      `json:"document,omitempty"`
		Consolidators []DisciplineListDocumentConsolidatorResponse `json:"consolidators,omitempty"`
		IsDueDate     bool                                         `json:"is_due_date"`
	}
)
