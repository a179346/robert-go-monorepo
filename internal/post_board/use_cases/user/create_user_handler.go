package user_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

type createUserRequestBody struct {
	Name     string `json:"name" validate:"required,gte=8,lte=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=32"`
}

func (u UserUseCase) createUserHandler(c *gohf.Context) gohf.Response {
	body, err := jsonvalidator.Validate[createUserRequestBody](c.Req.GetBody())
	if err != nil {
		return response.Error(
			http.StatusBadRequest,
			err,
		)
	}

	err = u.userCommands.createUser(
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
