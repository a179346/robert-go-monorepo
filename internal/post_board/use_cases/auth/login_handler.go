package auth_use_case

import (
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

type loginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=32"`
}

func (u AuthUseCase) loginHandler(c *gohf.Context) gohf.Response {
	bytes, ok := gohf_extended.BodyValue(c.Req.Context())
	if !ok {
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			tracerr.New("Failed to get body value"),
		)
	}

	body, err := jsonvalidator.FromBytes[loginRequestBody](bytes)
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			err.Error(),
			tracerr.Errorf("body valdiation error: %w", err),
		)
	}

	token, err := u.authCommands.login(c.Req.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, errUserNotFound) || errors.Is(err, errWrongPassword) {
			return gohf_extended.NewErrorResponse(
				http.StatusNotFound,
				"User not found",
				err,
			)
		}
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			err,
		)
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
