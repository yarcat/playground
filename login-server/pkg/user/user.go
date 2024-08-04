package user

import (
	"login-server/pkg/t"
	"time"
)

type (
	// ID is the user identifier.
	ID string
	// Name is the user name.
	Name string
	// SecretHash is the user password hash.
	SecretHash string
	// SecretSalt is the user password salt.
	SecretSalt string
	// Role is the user role.
	Role string

	// User is the user entity.
	User struct {
		t.Email // User email.

		ID         // User identifier.
		Name       // User name.
		Role       // User role.
		SecretHash // User password hash.
		SecretSalt // User password salt.

		Created time.Time // User creation time.
		// TODO: Store an algorithm identifier.
	}
)
