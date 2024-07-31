package postgres

import (
	"strings"

	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
)

func filterUser(filter req.UserFilter) (string, []interface{}) {
	values := []string{}
	args := []interface{}{}

	if filter.Email != "" {
		values = append(values, `email LIKE '%'||?||'%'`)
		args = append(args, filter.Email)
	}

	if !filter.CreatedAt.IsZero() {
		values = append(values, "created_at >= ?")
		args = append(args, filter.CreatedAt)
	}

	var whereClause string
	if len(values) > 0 {
		whereClause = " WHERE " + strings.Join(values, " AND ")
	}
	return whereClause, args
}
