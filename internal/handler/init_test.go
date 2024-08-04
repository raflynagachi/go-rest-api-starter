package handler

import (
	"reflect"
	"testing"

	"github.com/raflynagachi/go-rest-api-starter/config"
	hn "github.com/raflynagachi/go-rest-api-starter/internal/handler/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/handler/router"
	uc "github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition"
	"github.com/raflynagachi/go-rest-api-starter/internal/usecase/definition/mocks"
)

var (
	cfg         = &config.Config{}
	mockUc      = new(mocks.APIUsecase)
	mockHandler = router.New(cfg, New(mockUc))
)

func TestMain(m *testing.M) {
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
				usecase: mockUc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.usecase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
