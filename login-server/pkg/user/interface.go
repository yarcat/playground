package user

import (
	"context"
	t "login-server/types"
)

type (
	// Hasher is the interface that wraps the basic password hashing functionality.
	Hasher interface {
		Hash(t.Password) (t.SecretHash, t.SecretSalt, error)
	}
	// SaltHasher is a factory that creates a new Hasher with a salt.
	SaltHasher interface {
		New(t.SecretSalt) Hasher
	}

	// IDReader is the interface that wraps the basic FromID method.
	IDReader interface {
		FromID(context.Context, t.UserID) (*User, error)
	}
	// EmailReader is the interface that wraps the basic FromEmail method.
	EmailReader interface {
		FromEmail(context.Context, t.Email) (*User, error)
	}

	// NewOptions represents the user creation options.
	NewOptions struct {
		User         // User entity.
		t.SecretHash // User password.
		t.SecretSalt // User password salt.
	}
	// NewOptionFunc represents the user creation option.
	NewOptionFunc func(*NewOptions)

	// Factory is the interface that wraps the basic create functionality.
	Factory interface {
		New(context.Context, t.Name, t.Email, t.SecretHash, t.SecretSalt, ...NewOptionFunc) (*User, error)
	}
)

// WithRole sets the user role.
func WithRole(r t.Role) NewOptionFunc {
	return func(o *NewOptions) { o.Role = r }
}

// FillOptions fills the user creation options.
func FillOptions(out *NewOptions, n t.Name, e t.Email, sh t.SecretHash, ss t.SecretSalt, opts []NewOptionFunc) {
	out.Name, out.Email = n, e
	for _, opt := range opts {
		opt(out)
	}
}
