package post_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/authed_context"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/go-playground/validator/v10"
	"github.com/gohf-http/gohf/v4"
	"github.com/gohf-http/gohf/v4/gohf_responses"
)

type createPostRequestBody struct {
	Content string `json:"content" validate:"required"`
}

func (u PostUseCase) createPostHandler(c *gohf.Context) gohf.Response {
	authorId, ok := authed_context.Value(c.Req.Context())
	if !ok {
		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	var body createPostRequestBody

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

	err := u.postCommands.createPost(
		context.Background(),
		authorId,
		body.Content,
	)
	if err != nil {
		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse[interface{}](http.StatusOK, nil)
}
