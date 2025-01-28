package user_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
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
			"Something went wrong",
			tracerr.New("failed to get body value"),
			true,
		)
	}

	body, err := jsonvalidator.FromBytes[createUserRequestBody](bytes)
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			err.Error(),
			tracerr.Errorf("body validation error: %w", err),
			false,
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
				"email has been taken",
				err,
				false,
			)
		}

		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			err,
			true,
		)
	}

	return gohf_extended.NewCustomJsonResponse[interface{}](http.StatusOK, nil)
}
