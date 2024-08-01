package response

import (
	"time"

	"github.com/guregu/null/v5"
)

type ListResponse struct {
	Data               interface{} `json:"data"`
	PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Page      int   `json:"page"`
	TotalPage int64 `json:"total_page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
}

type CreatedResponse struct {
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}

type UpdatedResponse struct {
	UpdatedAt null.Time   `json:"updated_at"`
	UpdatedBy null.String `json:"updated_by"`
}
