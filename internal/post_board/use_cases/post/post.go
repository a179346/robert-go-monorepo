package post_use_case

import (
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/post_provider"
	"github.com/gohf-http/gohf/v6"
)

type PostUseCase struct {
	postQueries  postQueries
	postCommands postCommands
}

func New(postProvider post_provider.PostProvider) PostUseCase {
	return PostUseCase{
		postQueries:  newPostQueries(postProvider),
		postCommands: newPostCommands(postProvider),
	}
}

func (u PostUseCase) AppendHandler(router *gohf.Router) {
	router.Handle("GET /", u.getPostsHandler)
	router.Handle("POST /", u.createPostHandler)
}
