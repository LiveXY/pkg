package argon2id

import (
	"etms/pkg/logx"

	argon "github.com/alexedwards/argon2id"
	"go.uber.org/zap"
)

func Hash(password string) string {
	hash, err := argon.CreateHash(password, argon.DefaultParams)
	if err != nil {
		logx.Error.Error("argon2id CreateHash error", zap.Error(err))
	}
	return hash
}

func Verify(password, hash string) bool {
	match, err := argon.ComparePasswordAndHash(password, hash)
	return err == nil && match
}
