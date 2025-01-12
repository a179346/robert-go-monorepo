package post_use_case

import (
	"context"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/post_provider"
	"github.com/google/uuid"
)

type postCommands struct {
	postProvider post_provider.PostProvider
}

func newPostCommands(postProvider post_provider.PostProvider) postCommands {
	return postCommands{
		postProvider: postProvider,
	}
}

func (postCommands postCommands) createPost(ctx context.Context, authorId string, content string) error {
	authorUUID, err := uuid.Parse(authorId)
	if err != nil {
		return err
	}

	return postCommands.postProvider.CreatePost(ctx, authorUUID, content)
}
