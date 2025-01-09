package auth_use_case

import (
	"context"
	"errors"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	"github.com/a179346/robert-go-monorepo/pkg/cryption"
	"github.com/go-jet/jet/qrm"
)

type authCommands struct {
	userProvider user_provider.UserProvider
}

func newAuthCommands(userProvider user_provider.UserProvider) authCommands {
	return authCommands{userProvider: userProvider}
}

var errUserNotFound = errors.New("User not found")
var errWrongPassword = errors.New("Wrong password")

func (authCommands authCommands) login(ctx context.Context, email string, password string) (string, error) {
	user, err := authCommands.userProvider.FindByEmail(ctx, email)
	if err != nil {
		if err.Error() == qrm.ErrNoRows.Error() {
			return "", errUserNotFound
		}
		return "", err
	}

	encryptedPass := cryption.SHA256(password)
	if user.EncryptedPass != encryptedPass {
		return "", errWrongPassword
	}

	return "TODO - jwt", nil
}
