package dto

type (
	LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	LoginResponse struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}

	LoginWithGoogleResponse struct {
		NeedRegistration bool   `json:"need_registration"`
		Token            string `json:"token"`
		RegisterToken    string `json:"register_token"`
		Role             string `json:"role"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	ChangePasswordRequest struct {
		NewPassword string `json:"new_password"`
	}

	GetMe struct {
		PersonalInfo       PersonalInfo       `json:"personal_info"`
		UserDisciplineInfo UserDisciplineInfo `json:"user_discipline_info"`
	}

	PersonalInfo struct {
		ID           string  `json:"id"`
		Name         string  `json:"name"`
		Email        string  `json:"email"`
		Initial      string  `json:"initial"`
		Institution  string  `json:"institution"`
		PhotoProfile *string `json:"photo_profile"`
		Role         string  `json:"role"`
	}
)
