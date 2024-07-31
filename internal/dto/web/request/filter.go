package request

type UserFilter struct {
	Email string `json:"email"`
	Pagination
}
