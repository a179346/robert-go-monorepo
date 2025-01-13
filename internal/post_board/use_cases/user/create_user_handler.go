package user_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/go-playground/validator/v10"
	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
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
		return response.Error(
			http.StatusBadRequest,
			err,
		)
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return response.Error(
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
			return response.Error(
				http.StatusConflict,
				err,
			)
		}

		return response.Error(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse[interface{}](http.StatusOK, nil)
}
