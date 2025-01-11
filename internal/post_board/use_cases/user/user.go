package user_use_case

import (
	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/user_provider"
	"github.com/gohf-http/gohf/v3"
)

type UserUseCase struct {
	userQueries userQueries
}

func New(userProvider user_provider.UserProvider) UserUseCase {
	return UserUseCase{
		userQueries: newUserQueries(userProvider),
	}
}

func (u UserUseCase) AppendHandler(router *gohf.Router) {
	router.Handle("GET /me", u.getMeHandler)
}
