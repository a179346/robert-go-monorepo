package post_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/authed_context"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	"github.com/gohf-http/gohf/v6"
)

type createPostRequestBody struct {
	Content string `json:"content" validate:"required"`
}

func (u PostUseCase) createPostHandler(c *gohf.Context) gohf.Response {
	authorId, ok := authed_context.Value(c.Req.Context())
	if !ok {
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	body, err := jsonvalidator.Validate[createPostRequestBody](c.Req.GetBody())
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			err,
		)
	}

	err = u.postCommands.createPost(
		context.Background(),
		authorId,
		body.Content,
	)
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse[interface{}](http.StatusOK, nil)
}
