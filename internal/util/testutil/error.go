package testutil

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/pkg/database"
)

var (
	MockErr          = errors.New("mock error")
	MockErrDuplicate = &pq.Error{Code: database.ERR_PQ_CODE_DUPLICATE}
)
