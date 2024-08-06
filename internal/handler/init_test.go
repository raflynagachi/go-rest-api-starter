package handler

import (
	"io"
	"log"
	"reflect"
	"testing"

	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/handler/router"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition/mocks"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

var (
	cfg         = &config.Config{}
	mockUc      = new(mocks.APIUsecase)
	mockLogger  = logger.NewLogger()
	mockHandler = router.New(cfg, mockLogger, New(mockUc, mockLogger))
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	m.Run()
}

func TestNew(t *testing.T) {
	mockUc := new(mocks.APIUsecase)

	type args struct {
		usecase uc.APIUsecase
	}
	tests := []struct {
		name string
		args args
		want hn.APIHandler
	}{
		{
			name: "success",
			args: args{
				usecase: mockUc,
			},
			want: &APIHandlerImpl{
				usecase:   mockUc,
				appLogger: mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.usecase, mockLogger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
