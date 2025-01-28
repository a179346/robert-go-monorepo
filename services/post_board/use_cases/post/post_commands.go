package post_use_case

import (
	"context"

	"github.com/a179346/robert-go-monorepo/services/post_board/providers/post_provider"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
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
		return tracerr.Errorf("uuid parse error: %w", err)
	}

	err = postCommands.postProvider.CreatePost(ctx, authorUUID, content)
	if err != nil {
		return tracerr.Errorf("create post error: %w", err)
	}

	return nil
}
