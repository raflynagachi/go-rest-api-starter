package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/raflynagachi/go-rest-api-starter/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectDB(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	tests := []struct {
		name        string
		cfg         *config.Database
		setup       func()
		expectedErr bool
	}{
		{
			name: "success connect",
			cfg: &config.Database{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "password",
				Name:     "dbname",
			},
			setup: func() {
				sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) {
					return sqlxDB, nil
				}
			},
			expectedErr: false,
		},
		{
			name: "failed due to connection error",
			cfg: &config.Database{
				Host:     "localhost",
				Port:     5432,
				User:     "user",
				Password: "password",
				Name:     "dbname",
			},
			setup: func() {
				sqlxConnect = func(driverName, dataSourceName string) (*sqlx.DB, error) {
					return sqlxDB, sql.ErrConnDone
				}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpSqlxConnect := sqlxConnect
			defer func() {
				sqlxConnect = tmpSqlxConnect
			}()

			tt.setup()

			db, err := ConnectDB(tt.cfg)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}
