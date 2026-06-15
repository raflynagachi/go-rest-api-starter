package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/raflynagachi/go-rest-api-starter/pkg/auth"
	"github.com/raflynagachi/go-rest-api-starter/pkg/http/response"
	appjwt "github.com/raflynagachi/go-rest-api-starter/pkg/jwt"
	"github.com/raflynagachi/go-rest-api-starter/pkg/logger"
)

func JWTAuth(secretKey string, log *logger.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.WriteFromError(w, r, response.WrapErrUnauthorized(errors.New("missing authorization header")), log)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.WriteFromError(w, r, response.WrapErrUnauthorized(errors.New("invalid authorization header format")), log)
			return
		}

		claims, err := appjwt.ParseToken(parts[1], secretKey)
		if err != nil {
			response.WriteFromError(w, r, response.WrapErrUnauthorized(err), log)
			return
		}

		ctx := auth.SetEmail(r.Context(), claims.Email)
		next(w, r.WithContext(ctx), ps)
	}
}
