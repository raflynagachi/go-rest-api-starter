package request

import "time"

type UserFilter struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Pagination
}
