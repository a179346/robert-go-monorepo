package user_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/go-playground/validator/v10"
	"github.com/gohf-http/gohf/v5"
	"github.com/gohf-http/gohf/v5/gohf_responses"
)

type createUserRequestBody struct {
	Name     string `json:"name" validate:"required,gte=8,lte=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=32"`
}

func (u UserUseCase) createUserHandler(c *gohf.Context) gohf.Response {
	var body createUserRequestBody

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

	err := u.userCommands.createUser(
		context.Background(),
		body.Email,
		body.Name,
		body.Password,
	)
	if err != nil {
		if errors.Is(err, errDuplicatedEmail) {
			return gohf_responses.NewErrorResponse(
				http.StatusConflict,
				err,
			)
		}

		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse[interface{}](http.StatusOK, nil)
}
