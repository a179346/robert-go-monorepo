package auth_use_case

import (
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

type loginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=32"`
}

func (u AuthUseCase) loginHandler(c *gohf.Context) gohf.Response {
	body, err := jsonvalidator.Validate[loginRequestBody](c.Req.GetBody())
	if err != nil {
		return response.Error(
			http.StatusBadRequest,
			err,
		)
	}

	token, err := u.authCommands.login(c.Req.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, errUserNotFound) || errors.Is(err, errWrongPassword) {
			return response.Error(
				http.StatusNotFound,
				errors.New("User not found"),
			)
		}
		return response.Error(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
