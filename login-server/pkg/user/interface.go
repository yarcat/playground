package user

import (
	"context"
	"login-server/pkg/t"
)

type (
	// Hasher is the interface that wraps the basic password hashing functionality.
	Hasher interface {
		Hash(t.Password) (SecretHash, SecretSalt, error)
	}
	// SaltHasher is a factory that creates a new Hasher with a salt.
	SaltHasher interface {
		New(SecretSalt) Hasher
	}

	// IDReader is the interface that wraps the basic FromID method.
	IDReader interface {
		FromID(context.Context, ID) (*User, error)
	}
	// EmailReader is the interface that wraps the basic FromEmail method.
	EmailReader interface {
		FromEmail(context.Context, t.Email) (*User, error)
	}

	// NewOptions represents the user creation options.
	NewOptions struct {
		User       // User entity.
		SecretHash // User password.
		SecretSalt // User password salt.
	}
	// NewOptionFunc represents the user creation option.
	NewOptionFunc func(*NewOptions)

	// Factory is the interface that wraps the basic create functionality.
	Factory interface {
		New(context.Context, Name, t.Email, SecretHash, SecretSalt, ...NewOptionFunc) (*User, error)
	}
)

// WithRole sets the user role.
func WithRole(r Role) NewOptionFunc {
	return func(o *NewOptions) { o.Role = r }
}

// FillOptions fills the user creation options.
func FillOptions(out *NewOptions, n Name, e t.Email, sh SecretHash, ss SecretSalt, opts []NewOptionFunc) {
	out.Name, out.Email = n, e
	for _, opt := range opts {
		opt(out)
	}
}
