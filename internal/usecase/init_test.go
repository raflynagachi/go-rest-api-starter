package usecase

import (
	"reflect"
	"testing"

	"github.com/raflynagachi/go-rest-api-starter/config"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/repository/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/repository/definition/mocks"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
)

var (
	mockRepo = new(mocks.SQLRepo)
	mockCfg  = &config.Config{}
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestNew(t *testing.T) {
	mockRepo := new(mocks.SQLRepo)

	type args struct {
		cfg     *config.Config
		sqlRepo repo.SQLRepo
	}
	tests := []struct {
		name string
		args args
		want uc.APIUsecase
	}{
		{
			name: "success",
			args: args{
				cfg:     &config.Config{},
				sqlRepo: mockRepo,
			},
			want: &APIUsecaseImpl{
				cfg:  &config.Config{},
				repo: mockRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg, tt.args.sqlRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
