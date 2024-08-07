// Package inmem implements the in-memory storage for the user data.
package inmem

import (
	"context"
	"fmt"
	u "login-server/pkg/user"
	t "login-server/types"
	"sync"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
)

// Inmem is the in-memory storage for the user data. It is safe for concurrent use.
// The storage is not persistent and is lost when the application restarts.
// Zero values of Inmem are ready to use.
type Inmem struct {
	sync.Map
}

// New implements the user.Factory interface.
func (m *Inmem) New(ctx context.Context, n t.Name, e t.Email, h t.SecretHash, s t.SecretSalt, opts ...u.NewOptionFunc) (*u.User, error) {
	if _, ok := m.Load(e); ok {
		return nil, fmt.Errorf("inmem new: %w", u.EmailExists(e))
	}
	id, err := nanoid.New()
	if err != nil {
		return nil, fmt.Errorf("inmem new: %w", err)
	}
	var o u.NewOptions
	u.FillOptions(&o, n, e, h, s, opts)
	user := &u.User{
		UserID:     t.UserID(id),
		Created:    time.Now(),
		Name:       o.Name,
		Email:      o.Email,
		Role:       o.Role,
		SecretHash: o.SecretHash,
		SecretSalt: o.SecretSalt,
	}
	x, loaded := m.Map.LoadOrStore(user.Email, user)
	if loaded { // Something else stored the same email.
		return nil, fmt.Errorf("inmem new: %w", u.EmailExists(e))
	}
	return x.(*u.User), nil
}

// FromID implements the user.IDReader interface.
func (m *Inmem) FromID(ctx context.Context, id t.UserID) (*u.User, error) {
	var x *u.User
	m.Range(func(_, v interface{}) bool {
		if v.(*u.User).UserID == id {
			x = v.(*u.User)
			return false
		}
		return true
	})
	if x == nil {
		return nil, fmt.Errorf("inmem: %w", u.IDNotFound(id))
	}
	return x, nil
}

// FromEmail implements the user.EmailReader interface.
func (m *Inmem) FromEmail(ctx context.Context, e t.Email) (*u.User, error) {
	x, ok := m.Load(e)
	if !ok {
		return nil, fmt.Errorf("inmem: %w", u.EmailNotFound(e))
	}
	return x.(*u.User), nil
}
