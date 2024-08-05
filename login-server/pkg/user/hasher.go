package user

import (
	"login-server/pkg/crypto/argon2"
	t "login-server/types"
)

type (
	// Argon2Hasher is the Argon2 password hasher.
	Argon2Hasher struct {
		*argon2.Options
	}
)

// Argon2Hash hashes the password using Argon2.
func NewArgon2Hasher(o ...argon2.OptionFunc) Argon2Hasher {
	opts := new(argon2.Options)
	argon2.FillOptions(opts, o)
	return Argon2Hasher{opts}
}

// Hash hashes the password using Argon2. Hash implements the Hasher interface.
func (a Argon2Hasher) Hash(passwd t.Password) (t.SecretHash, t.SecretSalt, error) {
	opts := *a.Options
	var salt []byte
	argon2.WitchChainedOnSalt(func(s []byte) { salt = s })(&opts)
	hash, err := argon2.NewArgon2WithOptions([]byte(passwd), &opts)
	return t.SecretHash(string(hash)), t.SecretSalt(string(salt)), err
}
