package user_use_case

import (
	"errors"
	"net/http"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v6"
	"github.com/gohf-http/gohf/v6/response"
)

type getAllUsersElement struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type getAllUsersResponseBody []getAllUsersElement

func (u UserUseCase) getAllUsersHandler(c *gohf.Context) gohf.Response {
	users, err := u.userQueries.findAllUsers(c.Req.Context())
	if err != nil {
		return response.Error(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	data := make(getAllUsersResponseBody, 0)
	for _, user := range users {
		data = append(data, getAllUsersElement{
			ID:        user.ID.String(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, data)
}
