package auth_use_case

import (
	"context"
	"errors"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	"github.com/a179346/robert-go-monorepo/internal/post_board/shared/auth_jwt"
	"github.com/a179346/robert-go-monorepo/pkg/cryption"
	"github.com/go-jet/jet/qrm"
	"github.com/ztrue/tracerr"
)

type authCommands struct {
	userProvider user_provider.UserProvider
}

func newAuthCommands(userProvider user_provider.UserProvider) authCommands {
	return authCommands{
		userProvider: userProvider,
	}
}

var errUserNotFound = errors.New("User not found")
var errWrongPassword = errors.New("Wrong password")

func (authCommands authCommands) login(ctx context.Context, email string, password string) (string, error) {
	user, err := authCommands.userProvider.FindByEmail(ctx, email)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return "", tracerr.Wrap(errUserNotFound)
		}
		return "", tracerr.Errorf("find user by email error: %w", err)
	}

	encryptedPass := cryption.SHA256(password)
	if user.EncryptedPass != encryptedPass {
		return "", tracerr.Wrap(errWrongPassword)
	}

	token, err := auth_jwt.Sign(user.ID.String())
	if err != nil {
		return "", tracerr.Errorf("jwt sign error: %w", err)
	}
	return token, nil
}
