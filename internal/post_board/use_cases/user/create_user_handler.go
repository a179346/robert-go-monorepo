package user_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	"github.com/gohf-http/gohf/v6"
)

type createUserRequestBody struct {
	Name     string `json:"name" validate:"required,gte=8,lte=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=32"`
}

func (u UserUseCase) createUserHandler(c *gohf.Context) gohf.Response {
	bytes, ok := gohf_extended.BodyValue(c.Req.Context())
	if !ok {
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	body, err := jsonvalidator.FromBytes[createUserRequestBody](bytes)
	if err != nil {
		return gohf_extended.NewErrorResponse(
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
			return gohf_extended.NewErrorResponse(
				http.StatusConflict,
				err,
			)
		}

		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse[interface{}](http.StatusOK, nil)
}
