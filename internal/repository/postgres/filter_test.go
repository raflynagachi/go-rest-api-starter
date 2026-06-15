package postgres

import (
	"testing"
	"time"

	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/pkg/random"
	"github.com/stretchr/testify/assert"
)

func TestFilterUser(t *testing.T) {
	mockEmail := random.RandomEmail()
	mockTime := time.Now()

	tests := []struct {
		name       string
		filter     req.UserFilter
		wantClause string
		wantArgs   []interface{}
	}{
		{
			name: "success with filter",
			filter: req.UserFilter{
				Email: mockEmail,
			},
			wantClause: " WHERE deleted_at IS NULL AND email LIKE '%'||?||'%'",
			wantArgs:   []interface{}{mockEmail},
		},
		{
			name:       "success without filter",
			filter:     req.UserFilter{},
			wantClause: " WHERE deleted_at IS NULL",
			wantArgs:   nil,
		},
		{
			name: "success multiple filter",
			filter: req.UserFilter{
				Email:     mockEmail,
				CreatedAt: mockTime,
			},
			wantClause: " WHERE deleted_at IS NULL AND email LIKE '%'||?||'%' AND created_at >= ?",
			wantArgs:   []interface{}{mockEmail, mockTime},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClause, gotArgs := filterUser(tt.filter)
			assert.Equal(t, tt.wantClause, gotClause)
			assert.ElementsMatch(t, tt.wantArgs, gotArgs)
		})
	}
}
