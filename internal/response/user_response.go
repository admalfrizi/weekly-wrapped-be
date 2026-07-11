package response

import ("github.com/admalfrizi/weekly-wrapped-be/internal/model")

type UserResponse struct {
	ID       string    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func MapToUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}
}