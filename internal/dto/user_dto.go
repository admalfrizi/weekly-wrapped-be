package dto


type UserProfileResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
}