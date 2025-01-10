package auth_use_case

import (
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/jwt_provider"
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	"github.com/a179346/robert-go-monorepo/pkg/gohf"
)

type AuthUseCase struct {
	authCommands authCommands
}

func New(userProvider user_provider.UserProvider, jwtProvider jwt_provider.JwtProvider) AuthUseCase {
	return AuthUseCase{
		authCommands: newAuthCommands(userProvider, jwtProvider),
	}
}

func (u AuthUseCase) AppendHandler(router *gohf.Router) {
	router.Handle("POST /login", u.loginHandler)
}
