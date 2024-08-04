package jwt

import "login-server/pkg/crypto/argon2"

// Argon2 is the Argon2 password hasher.
type Argon2 struct {
	argon2.Argon2
}

// NewArgon2 creates a new Argon2 password hasher.
func NewArgon2(password string) (Argon2, error) {
	h, err := argon2.NewArgon2([]byte(password))
	return Argon2{h}, err
}

// Secret implements the SecretProvider interface.
func (a Argon2) Secret() []byte { return a.Argon2 }
