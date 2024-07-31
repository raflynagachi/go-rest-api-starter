package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomUser(t *testing.T) {
	user := RandomUser()

	require.NotNil(t, user)
	assert.NotEmpty(t, user.ID, "ID should not be empty")
	assert.NotEmpty(t, user.Email, "Email should not be empty")

	assert.Equal(t, "SYSTEM", user.Created.CreatedBy, "CreatedBy should be 'SYSTEM'")
	assert.Equal(t, "SYSTEM", user.Updated.UpdatedBy.String, "UpdatedBy should be 'SYSTEM'")
}
