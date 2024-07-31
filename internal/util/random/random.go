package random

import (
	"time"

	"github.com/guregu/null/v5"
	"github.com/raflynagachi/go-rest-api-starter/internal/model"
	"github.com/raflynagachi/go-rest-api-starter/pkg/random"
)

func RandomUser() *model.User {
	timeNow := time.Now()

	return &model.User{
		ID:    random.RandomID(),
		Email: random.RandomEmail(),
		Created: model.Created{
			CreatedAt: timeNow,
			CreatedBy: "SYSTEM",
		},
		Updated: model.Updated{
			UpdatedAt: null.TimeFrom(timeNow),
			UpdatedBy: null.StringFrom("SYSTEM"),
		},
	}
}
