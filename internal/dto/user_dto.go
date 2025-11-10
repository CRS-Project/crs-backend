package dto

type (
	CreateUserRequest struct {
		Name             string  `json:"name" binding:"required"`
		Email            string  `json:"email" binding:"required,email"`
		Password         string  `json:"password" binding:"required"`
		Initial          string  `json:"initial" binding:"required"`
		Institution      string  `json:"institution" binding:"required"`
		PhotoProfile     *string `json:"photo_profile" binding:""`
		Role             string  `json:"role" binding:"required,oneof=CONTRACTOR REVIEWER"`
		DisciplineNumber int     `json:"discipline_number" binding:"required"`
		PackageID        string  `json:"package_id" binding:"required"`
		DisciplineID     *string `json:"discipline_id" binding:""`
	}

	CreateUserResponse struct {
		ID               string  `json:"id"`
		Name             string  `json:"name"`
		Email            string  `json:"email"`
		Initial          string  `json:"initial"`
		Institution      string  `json:"institution"`
		PhotoProfile     *string `json:"photo_profile"`
		IsVerified       bool    `json:"is_verified"`
		Role             string  `json:"role"`
		DisciplineNumber int     `json:"discipline_number"`
		Package          string  `json:"package"`
		Discipline       string  `json:"discipline"`
		PackageID        *string `json:"package_id"`
		DisciplineID     string  `json:"discipline_id"`
	}

	UserNonAdminDetailResponse struct {
		ID               string  `json:"id"`
		Name             string  `json:"name"`
		Email            string  `json:"email"`
		Initial          string  `json:"initial"`
		Institution      string  `json:"institution"`
		PhotoProfile     *string `json:"photo_profile"`
		Role             string  `json:"role"`
		DisciplineNumber int     `json:"discipline_number"`
		Package          string  `json:"package"`
		Discipline       string  `json:"discipline"`
		PackageID        *string `json:"package_id"`
		DisciplineID     string  `json:"discipline_id"`
	}

	UpdateUserRequest struct {
		Name             string  `json:"name" binding:"required"`
		Email            string  `json:"email" binding:"required,email"`
		Password         *string `json:"password" binding:""`
		Initial          string  `json:"initial" binding:"required"`
		Institution      string  `json:"institution" binding:"required"`
		PhotoProfile     *string `json:"photo_profile" binding:""`
		DisciplineNumber int     `json:"discipline_number" binding:"required"`
		DisciplineID     *string `json:"discipline_id" binding:""`
	}

	UserComment struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		PhotoProfile *string `json:"photo_profile"`
		Role         string  `json:"role"`
	}
)
