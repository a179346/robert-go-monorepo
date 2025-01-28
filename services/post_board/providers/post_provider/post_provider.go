package post_provider

import (
	"context"
	"database/sql"

	. "github.com/a179346/robert-go-monorepo/services/post_board/database/.jet_gen/post-board/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/a179346/robert-go-monorepo/services/post_board/database/.jet_gen/post-board/public/model"
)

type PostProvider struct {
	db *sql.DB
}

func New(db *sql.DB) PostProvider {
	return PostProvider{db: db}
}

type PostResult struct {
	model.Post
	Author model.User
}

func (postProvider PostProvider) FindAll(ctx context.Context) ([]PostResult, error) {
	stmt := SELECT(
		Post.AllColumns,
		User.AllColumns.Except(User.EncryptedPass),
	).FROM(
		Post.
			INNER_JOIN(User, User.ID.EQ(Post.AuthorID)),
	).ORDER_BY(
		Post.CreatedAt.DESC(),
	)

	var dest []PostResult
	err := stmt.QueryContext(ctx, postProvider.db, &dest)
	return dest, err
}

func (postProvider PostProvider) FindByAuthor(ctx context.Context, authorId uuid.UUID) ([]PostResult, error) {
	stmt := SELECT(
		Post.AllColumns,
		User.AllColumns.Except(User.EncryptedPass),
	).FROM(
		Post.
			INNER_JOIN(User, User.ID.EQ(Post.AuthorID)),
	).WHERE(
		Post.AuthorID.EQ(UUID(authorId)),
	).ORDER_BY(
		Post.CreatedAt.DESC(),
	)

	var dest []PostResult
	err := stmt.QueryContext(ctx, postProvider.db, &dest)
	return dest, err
}

func (postProvider PostProvider) CreatePost(
	ctx context.Context,
	authorId uuid.UUID,
	content string,
) error {
	newPost := model.Post{
		ID:       uuid.New(),
		AuthorID: authorId,
		Content:  content,
	}

	columnList := ColumnList{Post.ID, Post.AuthorID, Post.Content}
	stmt := Post.INSERT(columnList).MODEL(newPost)

	_, err := stmt.ExecContext(ctx, postProvider.db)
	return err
}
