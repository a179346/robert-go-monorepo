package user_use_case

import (
	"context"
	"errors"

	"github.com/a179346/robert-go-monorepo/internal/post_board/database/.jet_gen/post-board/public/model"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	"github.com/go-jet/jet/qrm"
	"github.com/google/uuid"
)

type userQueries struct {
	userProvider user_provider.UserProvider
}

func newUserQueries(userProvider user_provider.UserProvider) userQueries {
	return userQueries{
		userProvider: userProvider,
	}
}

var errUserNotFound = errors.New("User not found")

func (userQueries userQueries) findUserById(ctx context.Context, userId string) (model.User, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return model.User{}, err
	}

	user, err := userQueries.userProvider.FindById(ctx, id)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return model.User{}, errUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (userQueries userQueries) findAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := userQueries.userProvider.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
