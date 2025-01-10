package jwt_provider_test

import (
	"testing"

	"github.com/a179346/robert-go-monorepo/internal/post_board/providers/jwt_provider"
)

func TestJwtProvider(t *testing.T) {
	secret := "my-test-secret"
	jwtProvider := jwt_provider.New(secret, 10)

	t.Run("test jwt provider", func(t *testing.T) {
		id := "asdjiovcasc1235"

		token, err := jwtProvider.Sign(id)
		if err != nil {
			t.Errorf("sign error: %v", err)
			return
		}

		claims, err := jwtProvider.Parse(token)
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
