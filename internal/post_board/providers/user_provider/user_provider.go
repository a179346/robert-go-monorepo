package user_provider

import (
	"context"
	"database/sql"

	. "github.com/a179346/robert-go-monorepo/internal/post_board/database/.jet_gen/post-board/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"

	"github.com/a179346/robert-go-monorepo/internal/post_board/database/.jet_gen/post-board/public/model"
)

type UserProvider struct {
	db *sql.DB
}

func New(db *sql.DB) UserProvider {
	return UserProvider{db: db}
}

func (userProvider UserProvider) FindAll(ctx context.Context) ([]model.User, error) {
	stmt := SELECT(
		User.AllColumns,
	).FROM(
		User,
	)

	var dest []model.User
	err := stmt.QueryContext(ctx, userProvider.db, &dest)
	return dest, err
}

func (userProvider UserProvider) FindById(ctx context.Context, id uuid.UUID) (model.User, error) {
	stmt := SELECT(
		User.AllColumns,
	).FROM(
		User,
	).WHERE(
		User.ID.EQ(UUID(id)),
	).LIMIT(1)

	var dest model.User
	err := stmt.QueryContext(ctx, userProvider.db, &dest)
	return dest, err
}

func (userProvider UserProvider) FindByEmail(ctx context.Context, email string) (model.User, error) {
	stmt := SELECT(
		User.AllColumns,
	).FROM(
		User,
	).WHERE(
		User.Email.EQ(Text(email)),
	).LIMIT(1)

	var dest model.User
	err := stmt.QueryContext(ctx, userProvider.db, &dest)
	return dest, err
}

func (userProvider UserProvider) CreateUser(
	ctx context.Context,
	email string,
	name string,
	encryptedPass string,
) error {
	newUser := model.User{
		ID:            uuid.New(),
		Email:         email,
		Name:          name,
		EncryptedPass: encryptedPass,
	}

	columnList := ColumnList{User.ID, User.Email, User.Name, User.EncryptedPass}
	stmt := User.INSERT(columnList).MODEL(newUser)

	_, err := stmt.ExecContext(ctx, userProvider.db)
	return err
}
