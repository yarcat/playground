package user

import (
	t "login-server/types"
	"time"
)

// User is the user entity.
type User struct {
	t.Email      // User email.
	t.UserID     // User identifier.
	t.Name       // User name.
	t.Role       // User role.
	t.SecretHash // User password hash.
	t.SecretSalt // User password salt.

	Created time.Time // User creation time.
	// TODO: Store an algorithm identifier.
}
