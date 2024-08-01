package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestTxBegin(t *testing.T) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		args    args
		setup   func()
		wantTx  bool
		wantErr bool
	}{
		{
			name: "success create tx",
			args: args{
				db: sqlxDB,
			},
			setup: func() {
				mockSql.ExpectBegin()
			},
			wantTx:  true,
			wantErr: false,
		},
		{
			name: "failed due to beginx error",
			args: args{
				db: sqlxDB,
			},
			setup: func() {
				mockSql.ExpectBegin().WillReturnError(errors.New("some error"))
			},
			wantTx:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := TxBegin(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxBegin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantTx {
				t.Errorf("TxBegin() = %v, wantTx %v", got, tt.wantTx)
			}
		})
	}
}

func TestTxRollback(t *testing.T) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "success rollback tx",
			setup: func() {
				mockSql.ExpectRollback()
			},
			wantErr: false,
		},
		{
			name: "failed due to rollback error",
			setup: func() {
				mockSql.ExpectRollback().WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSql.ExpectBegin()
			tx, err := sqlxDB.Beginx()
			if err != nil {
				t.Fatalf("failed to mock Tx: %v", err)
			}

			tt.setup()

			err = TxRollback(tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxRollback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Errorf("TxRollback() unmet expectations: %v", err)
			}
		})
	}
}

func TestTxCommit(t *testing.T) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "success commit tx",
			setup: func() {
				mockSql.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "failed due to commit error",
			setup: func() {
				mockSql.ExpectCommit().WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSql.ExpectBegin()
			tx, err := sqlxDB.Beginx()
			if err != nil {
				t.Fatalf("failed to mock Tx: %v", err)
			}

			tt.setup()

			err = TxCommit(tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Errorf("TxCommit() unmet expectations: %v", err)
			}
		})
	}
}
