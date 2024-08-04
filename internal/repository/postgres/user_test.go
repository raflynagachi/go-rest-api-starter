package postgres

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	req "github.com/raflynagachi/go-rest-api-starter/internal/dto/web/request"
	"github.com/raflynagachi/go-rest-api-starter/internal/model"
	randomutil "github.com/raflynagachi/go-rest-api-starter/internal/util/random"
	"github.com/raflynagachi/go-rest-api-starter/internal/util/testutil"
	"github.com/raflynagachi/go-rest-api-starter/pkg/random"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepo_GetUser(t *testing.T) {
	mockUser := randomutil.RandomUser()
	mockUsers := []*model.User{mockUser}

	mockFilter := req.UserFilter{
		Email: mockUser.Email,
		Pagination: req.Pagination{
			Page:  random.RandomInt(1, 10),
			Limit: random.RandomInt(1, 10),
		},
	}

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		filter req.UserFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		want    []*model.User
		wantErr bool
	}{
		{
			name:   "success get user with filter",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), filter: mockFilter},
			setup: func() {
				mockSql.ExpectQuery("SELECT").
					WithArgs(mockUser.Email).
					WillReturnRows(
						mockSql.NewRows([]string{"id", "email", "created_at", "created_by", "updated_at", "updated_by"}).
							AddRow(mockUser.ID, mockUser.Email, mockUser.CreatedAt, mockUser.CreatedBy, mockUser.UpdatedAt, mockUser.UpdatedBy))
			},
			want:    mockUsers,
			wantErr: false,
		},
		{
			name:    "failed due to invalid pagination",
			fields:  fields{DB: sqlxDB},
			args:    args{ctx: context.Background(), filter: req.UserFilter{}},
			setup:   func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "failed due to connection error",
			fields: fields{DB: sqlxDB},
			args: args{ctx: context.Background(), filter: req.UserFilter{
				Pagination: mockFilter.Pagination,
			}},
			setup: func() {
				mockSql.ExpectQuery("SELECT").WithoutArgs().
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostgresRepo{
				DB: tt.fields.DB,
			}

			tt.setup()

			got, err := r.GetUser(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepo.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresRepo.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_CountUser(t *testing.T) {
	mockUser := randomutil.RandomUser()
	mockCount := int64(1)

	mockFilter := req.UserFilter{
		Email: mockUser.Email,
		Pagination: req.Pagination{
			Page:  random.RandomInt(1, 10),
			Limit: random.RandomInt(1, 10),
		},
	}

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		filter req.UserFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		want    int64
		wantErr bool
	}{
		{
			name:   "success count user with filter",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), filter: mockFilter},
			setup: func() {
				mockSql.ExpectQuery("SELECT").
					WithArgs(mockUser.Email).
					WillReturnRows(
						mockSql.NewRows([]string{"count"}).
							AddRow(mockCount))
			},
			want:    mockCount,
			wantErr: false,
		},
		{
			name:   "success without filter",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), filter: req.UserFilter{}},
			setup: func() {
				mockSql.ExpectQuery("SELECT").WithoutArgs().
					WillReturnRows(
						mockSql.NewRows([]string{"count"}).
							AddRow(mockCount))
			},
			want:    mockCount,
			wantErr: false,
		},
		{
			name:   "failed due to connection error",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), filter: req.UserFilter{}},
			setup: func() {
				mockSql.ExpectQuery("SELECT").WithoutArgs().
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostgresRepo{
				DB: tt.fields.DB,
			}

			tt.setup()

			got, err := r.CountUser(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepo.CountUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresRepo.CountUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_GetUserByID(t *testing.T) {
	mockUser := randomutil.RandomUser()

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		want    *model.User
		wantErr bool
	}{
		{
			name:   "success get by userID",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), id: mockUser.ID},
			setup: func() {
				mockSql.ExpectQuery("SELECT").
					WithArgs(mockUser.ID).
					WillReturnRows(
						mockSql.NewRows([]string{"id", "email", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by"}).
							AddRow(mockUser.ID, mockUser.Email, mockUser.CreatedAt, mockUser.CreatedBy, mockUser.UpdatedAt, mockUser.UpdatedBy, mockUser.DeletedAt, mockUser.DeletedBy))
			},
			want:    mockUser,
			wantErr: false,
		},
		{
			name:   "failed due to user not found",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), id: mockUser.ID},
			setup: func() {
				mockSql.ExpectQuery("SELECT").
					WithArgs(mockUser.ID).
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "failed due to connection error",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), id: mockUser.ID},
			setup: func() {
				mockSql.ExpectQuery("SELECT").
					WithArgs(mockUser.ID).
					WillReturnError(sql.ErrConnDone)
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PostgresRepo{
				DB: tt.fields.DB,
			}

			tt.setup()

			got, err := r.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepo.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresRepo.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_InsertUser(t *testing.T) {
	mockUser := randomutil.RandomUser()

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *model.User
		tx   *sqlx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		want    int64
		wantErr bool
	}{
		{
			name:   "success insert user",
			fields: fields{DB: sqlxDB},
			args: args{
				ctx:  context.Background(),
				user: mockUser,
			},
			setup: func() {
				mockSql.ExpectQuery("INSERT").WithArgs(mockUser.Email, mockUser.CreatedAt, mockUser.CreatedBy).
					WillReturnRows(
						mockSql.NewRows([]string{"id"}).
							AddRow(mockUser.ID))
			},
			want:    mockUser.ID,
			wantErr: false,
		},
		{
			name:   "failed due to unique email conflict",
			fields: fields{DB: sqlxDB},
			args: args{
				ctx:  context.Background(),
				user: mockUser,
			},
			setup: func() {
				mockSql.ExpectQuery("INSERT").WithArgs(mockUser.Email, mockUser.CreatedAt, mockUser.CreatedBy).
					WillReturnError(testutil.MockErrDuplicate)
			},
			wantErr: true,
		},
		{
			name:   "failed due to connection error",
			fields: fields{DB: sqlxDB},
			args: args{
				ctx:  context.Background(),
				user: mockUser,
			},
			setup: func() {
				mockSql.ExpectQuery("INSERT").WithArgs(mockUser.Email, mockUser.CreatedAt, mockUser.CreatedBy).
					WillReturnError(sql.ErrConnDone)
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
			tt.setup()

			var err error
			tt.args.tx, err = testutil.InitBeginx(sqlxDB)
			assert.NoError(t, err)

			got, err := r.InsertUser(tt.args.ctx, tt.args.tx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepo.InsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostgresRepo.InsertUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepo_UpdateUser(t *testing.T) {
	mockUser := randomutil.RandomUser()

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user *model.User
		tx   *sqlx.Tx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func()
		wantErr bool
	}{
		{
			name:   "success update user",
			fields: fields{DB: sqlxDB},
			args: args{
				ctx:  context.Background(),
				user: mockUser,
			},
			setup: func() {
				mockSql.ExpectExec("UPDATE").WithArgs(mockUser.Email, mockUser.UpdatedAt, mockUser.UpdatedBy, mockUser.ID).
					WillReturnResult(sqlmock.NewResult(mockUser.ID, 1))
			},
			wantErr: false,
		},
		{
			name:   "failed due to unique constraint error",
			fields: fields{DB: sqlxDB},
			args:   args{ctx: context.Background(), user: mockUser},
			setup: func() {
				mockSql.ExpectExec("UPDATE").WithArgs(mockUser.Email, mockUser.UpdatedAt, mockUser.UpdatedBy, mockUser.ID).
					WillReturnError(testutil.MockErrDuplicate)
			},
			wantErr: true,
		},
		{
			name:   "failed due to connection error",
			fields: fields{DB: sqlxDB},
			args: args{
				ctx:  context.Background(),
				user: mockUser,
			},
			setup: func() {
				mockSql.ExpectExec("UPDATE").WithArgs(mockUser.Email, mockUser.UpdatedAt, mockUser.UpdatedBy, mockUser.ID).
					WillReturnError(sql.ErrConnDone)
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
			tt.setup()

			var err error
			tt.args.tx, err = testutil.InitBeginx(sqlxDB)
			assert.NoError(t, err)

			if err := r.UpdateUser(tt.args.ctx, tt.args.tx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("PostgresRepo.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
