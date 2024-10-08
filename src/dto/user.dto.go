package dto

type ResendOtpRequest struct {
	Email string `json:"email" validate:"required,email,excludesall=;()="`
}

type VerifyUserRequest struct {
	Email string `json:"email" validate:"required,email,excludesall=;()="`
	Otp   string `json:"otp" validate:"required,len=6,excludesall=;()="`
}

type SignupRequest struct {
	Email     string `json:"email" validate:"required,email,excludesall=;()="`
	Password  string `json:"password" validate:"required,excludesall=;()="`
	FirstName string `json:"firstName" validate:"required,excludesall=;()="`
	LastName  string `json:"lastName" validate:"required,excludesall=;()="`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,excludesall=;()="`
	Password string `json:"password" validate:"required,excludesall=;()="`
}

type UpdateProfileRequest struct {
	FirstName string `json:"firstName" form:"firstName" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"lastName" form:"lastName" validate:"omitempty,min=2,max=50"`
	// Note: Image is handled separately in the form-data, so it's not included here.
}
