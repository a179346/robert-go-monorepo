package post_use_case

import (
	"context"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/post_provider"
	"github.com/google/uuid"
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
		return postQueries.postProvider.FindAll(ctx)
	}

	authorUUID, err := uuid.Parse(authorId)
	if err != nil {
		return nil, err
	}
	return postQueries.postProvider.FindByAuthor(ctx, authorUUID)
}
