package user_use_case

import (
	"github.com/a179346/robert-go-monorepo/services/post_board/providers/user_provider"
	"github.com/gohf-http/gohf/v6"
)

type UserUseCase struct {
	userQueries  userQueries
	userCommands userCommands
}

func New(userProvider user_provider.UserProvider) UserUseCase {
	return UserUseCase{
		userQueries:  newUserQueries(userProvider),
		userCommands: newUserCommands(userProvider),
	}
}

func (u UserUseCase) AppendHandler(router *gohf.Router) {
	router.GET("/me", u.getMeHandler)
	router.GET("/", u.getAllUsersHandler)
	router.POST("/", u.createUserHandler)
}
