package user_use_case

import (
	"errors"
	"net/http"
	"time"

	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/authed_context"
	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf"
	"github.com/gohf-http/gohf/gohf_responses"
)

type getMeResponseBody struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (u UserUseCase) getMeHandler(c *gohf.Context) gohf.Response {
	userId, ok := authed_context.Value(c.Req.Context())
	if !ok {
		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	user, err := u.userQueries.findUserById(c.Req.Context(), userId)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			return gohf_responses.NewErrorResponse(
				http.StatusUnauthorized,
				errors.New("User not found"),
			)
		}
		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, getMeResponseBody{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}
