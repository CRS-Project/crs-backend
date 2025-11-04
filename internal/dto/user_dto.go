package dto

type (
	CreateUserRequest struct {
		Name         string  `json:"name" binding:"required"`
		Email        string  `json:"email" binding:"required,email"`
		Password     string  `json:"password" binding:"required"`
		Initial      string  `json:"initial" binding:"required"`
		Institution  string  `json:"institution" binding:"required"`
		Role         string  `json:"role" binding:"required,oneof=CONTRACTOR REVIEWER"`
		PackageID    string  `json:"package_id" binding:"required"`
		DisciplineID *string `json:"discipline_id" binding:""`
	}

	CreateUserResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		Initial     string `json:"initial"`
		Institution string `json:"institution"`
		IsVerified  bool   `json:"is_verified"`
		Role        string `json:"role"`
		Package     string `json:"package"`
		Discipline  string `json:"discipline"`
	}

	UserNonAdminDetailResponse struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		Email        string  `json:"email"`
		Initial      string  `json:"initial"`
		Institution  string  `json:"institution"`
		PhotoProfile *string `json:"photo_profile"`
		Package      string  `json:"package"`
		Discipline   string  `json:"discipline"`
		Role         string  `json:"role"`
	}

	UpdateUserRequest struct {
		Name         string  `json:"name" binding:"required"`
		Email        string  `json:"email" binding:"required,email"`
		Password     string  `json:"password" binding:"required"`
		Initial      string  `json:"initial" binding:"required"`
		Institution  string  `json:"institution" binding:"required"`
		PackageID    string  `json:"package_id" binding:"required"`
		DisciplineID *string `json:"discipline_id" binding:""`
	}
)
