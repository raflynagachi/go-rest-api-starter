package postgres

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/raflynagachi/go-rest-api-starter/internal/util/testutil"
)

func TestPostgresRepo_TxBegin(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		wantTx  bool
		wantErr bool
	}{
		{
			name: "success begin transaction",
			fields: fields{
				DB: sqlxDB,
			},
			setup: func() {
				mockSql.ExpectBegin()
			},
			wantTx:  true,
			wantErr: false,
		},
		{
			name: "failed due to begin transaction error",
			fields: fields{
				DB: sqlxDB,
			},
			setup: func() {
				mockSql.ExpectBegin().WillReturnError(testutil.MockErr)
			},
			wantTx:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostgresRepo{
				DB: tt.fields.DB,
			}

			tt.setup()

			got, err := r.TxBegin()
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepo.TxBegin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantTx {
				t.Errorf("PostgresRepo.TxBegin() = %v, want %v", got, tt.wantTx)
			}
		})
	}
}

func TestPostgresRepo_TxEnd(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name: "success commit",
			fields: fields{
				DB: sqlxDB,
			},
			args: args{
				err: nil,
			},
			setup: func() {
				mockSql.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "failed due to commit error",
			fields: fields{
				DB: sqlxDB,
			},
			args: args{
				err: nil,
			},
			setup: func() {
				mockSql.ExpectCommit().WillReturnError(testutil.MockErr)
			},
			wantErr: true,
		},
		{
			name: "failed due to transaction rollback",
			fields: fields{
				DB: sqlxDB,
			},
			args: args{
				err: testutil.MockErr,
			},
			setup: func() {
				mockSql.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "failed due to transaction rollback error",
			fields: fields{
				DB: sqlxDB,
			},
			args: args{
				err: testutil.MockErr,
			},
			setup: func() {
				mockSql.ExpectRollback().WillReturnError(testutil.MockErr)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostgresRepo{
				DB: tt.fields.DB,
			}

			mockSql.ExpectBegin()
			mockTx, err := testutil.InitBeginx(sqlxDB)
			if err != nil {
				t.Fatalf("failed to mock Tx: %v", err)
			}

			tt.setup()

			if err := r.TxEnd(mockTx, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepo.TxEnd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
