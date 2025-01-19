package post_use_case

import (
	"context"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/post_provider"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

type postQueries struct {
	postProvider post_provider.PostProvider
}

func newPostQueries(postProvider post_provider.PostProvider) postQueries {
	return postQueries{
		postProvider: postProvider,
	}
}

func (postQueries postQueries) find(ctx context.Context, authorId string) ([]post_provider.PostResult, error) {
	if authorId == "" {
		result, err := postQueries.postProvider.FindAll(ctx)
		if err != nil {
			return result, tracerr.Errorf("find posts error: %w", err)
		}
		return result, nil
	}

	authorUUID, err := uuid.Parse(authorId)
	if err != nil {
		return nil, tracerr.Errorf("uuid parse error: %w", err)
	}

	result, err := postQueries.postProvider.FindByAuthor(ctx, authorUUID)
	if err != nil {
		return result, tracerr.Errorf("find posts error: %w", err)
	}
	return result, nil
}
