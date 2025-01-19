package user_use_case

import (
	"context"
	"errors"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	"github.com/a179346/robert-go-monorepo/pkg/cryption"
	"github.com/lib/pq"
	"github.com/ztrue/tracerr"
)

type userCommands struct {
	userProvider user_provider.UserProvider
}

func newUserCommands(userProvider user_provider.UserProvider) userCommands {
	return userCommands{
		userProvider: userProvider,
	}
}

var errDuplicatedEmail = errors.New("email has been taken")

func (userCommands userCommands) createUser(ctx context.Context, email string, name string, password string) error {
	encryptedPass := cryption.SHA256(password)
	err := userCommands.userProvider.CreateUser(ctx, email, name, encryptedPass)
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" && err.Constraint == "user_email_key" {
			return tracerr.Wrap(errDuplicatedEmail)
		}
	}
	if err != nil {
		return tracerr.Errorf("create user error: %w", err)
	}
	return nil
}
