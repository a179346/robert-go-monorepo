package post_use_case

import (
	"context"
	"net/http"

	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/authed_context"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/a179346/robert-go-monorepo/pkg/jsonvalidator"
	"github.com/gohf-http/gohf/v6"
	"github.com/ztrue/tracerr"
)

type createPostRequestBody struct {
	Content string `json:"content" validate:"required"`
}

func (u PostUseCase) createPostHandler(c *gohf.Context) gohf.Response {
	authorId, ok := authed_context.Value(c.Req.Context())
	if !ok {
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			tracerr.New("failed to get user id"),
		)
	}

	bytes, ok := gohf_extended.BodyValue(c.Req.Context())
	if !ok {
		return gohf_extended.NewErrorResponse(
			http.StatusInternalServerError,
			"Something went wrong",
			tracerr.New("failed to get body value"),
		)
	}

	body, err := jsonvalidator.FromBytes[createPostRequestBody](bytes)
	if err != nil {
		return gohf_extended.NewErrorResponse(
			http.StatusBadRequest,
			err.Error(),
			tracerr.Errorf("body valdiation error: %w", err),
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
			"Something went wrong",
			err,
		)
	}

	return gohf_extended.NewCustomJsonResponse[interface{}](http.StatusOK, nil)
}
