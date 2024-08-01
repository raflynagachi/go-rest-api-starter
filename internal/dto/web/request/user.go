package request

type CreateUpdateUserReq struct {
	Email string `json:"email" validate:"required,email"`
}
