package response

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	CreatedResponse
	UpdatedResponse
}
