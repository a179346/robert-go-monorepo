package post_use_case

import (
	"errors"
	"net/http"
	"time"

	"github.com/a179346/robert-go-monorepo/pkg/gohf_extended"
	"github.com/gohf-http/gohf/v4"
	"github.com/gohf-http/gohf/v4/gohf_responses"
)

type author struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type getPostElement struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Author    author `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type getPostsResponseBody []getPostElement

func (u PostUseCase) getPostsHandler(c *gohf.Context) gohf.Response {
	posts, err := u.postQueries.find(c.Req.Context(), c.Req.GetQuery("authorId"))
	if err != nil {
		return gohf_responses.NewErrorResponse(
			http.StatusInternalServerError,
			errors.New("Something went wrong"),
		)
	}

	data := make(getPostsResponseBody, 0)
	for _, post := range posts {
		author := author{
			ID:        post.Author.ID.String(),
			Email:     post.Author.Email,
			Name:      post.Author.Name,
			CreatedAt: post.Author.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.Author.UpdatedAt.Format(time.RFC3339),
		}

		data = append(data, getPostElement{
			ID:        post.ID.String(),
			Content:   post.Content,
			Author:    author,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		})
	}

	return gohf_extended.NewCustomJsonResponse(http.StatusOK, data)
}
