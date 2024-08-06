package usecase

import (
	"io"
	"log"
	"reflect"
	"testing"

	"github.com/raflynagachi/go-rest-api-starter/config"
	repo "github.com/raflynagachi/go-rest-api-starter/internal/repository/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/repository/definition/mocks"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

var (
	mockRepo   = new(mocks.SQLRepo)
	mockCfg    = &config.Config{}
	mockLogger = logger.NewLogger()
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	m.Run()
}

func TestNew(t *testing.T) {
	mockRepo := new(mocks.SQLRepo)

	type args struct {
		cfg       *config.Config
		sqlRepo   repo.SQLRepo
		appLogger *logger.Logger
	}
	tests := []struct {
		name string
		args args
		want uc.APIUsecase
	}{
		{
			name: "success",
			args: args{
				cfg:       mockCfg,
				sqlRepo:   mockRepo,
				appLogger: mockLogger,
			},
			want: &APIUsecaseImpl{
				cfg:       mockCfg,
				repo:      mockRepo,
				appLogger: mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg, tt.args.appLogger, tt.args.sqlRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
