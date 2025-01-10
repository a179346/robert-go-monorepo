package auth_use_case

import (
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf"
	"github.com/a179346/robert-go-monorepo/pkg/gohf/gohf_responses"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/go-playground/validator/v10"
)

type loginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=32"`
}

func (u AuthUseCase) loginHandler(c *gohf.Context) gohf.Response {
	var body loginRequestBody

	defer c.Req.GetBody().Close()
	if err := c.Req.GetBody().JsonDecode(&body); err != nil {
		return gohf_responses.NewErrorResponse(
			http.StatusBadRequest,
			err,
		)
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return gohf_responses.NewErrorResponse(
			http.StatusBadRequest,
			err,
		)
	}

	token, err := u.authCommands.login(c.Req.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, errUserNotFound) || errors.Is(err, errWrongPassword) {
			return gohf_responses.NewErrorResponse(
				http.StatusNotFound,
				errors.New("User not found"),
			)
		}
		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}