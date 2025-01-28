package auth_jwt_test

import (
	"testing"

	"github.com/a179346/robert-go-monorepo/services/post_board/shared/auth_jwt"
)

func TestJwtProvider(t *testing.T) {
	t.Run("test jwt provider", func(t *testing.T) {
		id := "asdjiovcasc1235"

		token, err := auth_jwt.Sign(id)
		if err != nil {
			t.Errorf("sign error: %v", err)
			return
		}

		claims, err := auth_jwt.Parse(token)
		if err != nil {
			t.Errorf("parse error: %v", err)
			return
		}
		if claims == nil {
			t.Errorf("parsed claims is nil")
			return
		}
		if got, want := claims.ID, id; got != want {
			t.Errorf("id should be identical -> got:%s want:%s", got, want)
			return
		}
	})
}
