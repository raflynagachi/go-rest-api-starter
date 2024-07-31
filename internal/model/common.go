package model

import (
	"time"

	"github.com/guregu/null/v5"
)

type Created struct {
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
}

type Updated struct {
	UpdatedAt null.Time   `db:"updated_at"`
	UpdatedBy null.String `db:"updated_by"`
}

type Deleted struct {
	DeletedAt null.Time   `db:"deleted_at"`
	DeletedBy null.String `db:"deleted_by"`
}
