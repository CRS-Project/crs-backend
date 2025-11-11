package dto

type (
	CreatePackageRequest struct {
		Name string `json:"name" binding:"required"`
	}

	UpdatePackageRequest struct {
		ID   string `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}

	PackageInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)
