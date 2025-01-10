package auth_use_case

import (
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	"github.com/gohf-http/gohf"
)

type AuthUseCase struct {
	authCommands authCommands
}

func New(userProvider user_provider.UserProvider) AuthUseCase {
	return AuthUseCase{
		authCommands: newAuthCommands(userProvider),
	}
}

func (u AuthUseCase) AppendHandler(router *gohf.Router) {
	router.Handle("POST /login", u.loginHandler)
}
