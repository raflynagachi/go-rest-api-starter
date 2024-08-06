package postgres

import (
	"database/sql"
	"io"
	"log"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/raflynagachi/go-rest-api-starter/internal/util/testutil"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

var (
	mockSql    sqlmock.Sqlmock
	mockLogger = logger.NewLogger()
	db         *sql.DB
	sqlxDB     *sqlx.DB
)

func TestMain(m *testing.M) {
	var err error
	db, sqlxDB, mockSql, err = testutil.InitMockDB()
	if err != nil {
		log.Fatalf("failed to mock DB: %v", err)
	}
	defer db.Close()

	log.SetOutput(io.Discard)
	m.Run()
}

func TestNew(t *testing.T) {
	type args struct {
		db        *sqlx.DB
		appLogger *logger.Logger
	}
	tests := []struct {
		name string
		args args
		want *PostgresRepo
	}{
		{
			name: "success",
			args: args{
				db:        sqlxDB,
				appLogger: mockLogger,
			},
			want: &PostgresRepo{
				DB:        sqlxDB,
				appLogger: mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.db, tt.args.appLogger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
