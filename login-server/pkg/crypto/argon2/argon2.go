package argon2

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type (
	// OnSaltFunc is a callback function that is called with the generated salt.
	OnSaltFunc func([]byte)
	// Options represents the Argon2 options.
	Options struct {
		// Time is the number of iterations to use. Defaults to 1.
		Time uint32
		// Memory is the amount of memory to use (in kibibytes).
		// Memory must be a power of 2 greater than 1.
		// Defaults to 64 MB.
		Memory uint32
		// Threads is the number of threads to use. Defaults to 4.
		Threads uint8
		// KeyLen is the length of the generated key. Defaults to 32 bytes.
		KeyLen uint32
		// Salt is the salt to use. If nil, a random salt will be generated.
		Salt []byte
		// SaltLen is the length of the generated salt. This argument is used only
		// if the salt argument is nil. Defaults to 32 bytes.
		SaltLen uint32
		// OnSalt is a callback function that is called with the generated salt.
		OnSalt OnSaltFunc
	}
	// OptionFunc represents the Argon2 option.
	OptionFunc func(*Options)
	// Argon2 is a byte slice that holds the Argon2 hash of a password.
	Argon2 []byte
)

// WithSalt sets the salt to use.
func WithSalt(salt []byte) OptionFunc {
	return func(o *Options) { o.Salt = salt }
}

// WithOnSalt sets the OnSalt callback function. The callback function is called
// with the generated salt.
func WithOnSalt(f OnSaltFunc) OptionFunc {
	return func(o *Options) { o.OnSalt = f }
}

// WitchChainedOnSalt sets the OnSalt callback function. The callback function is
// called with the generated salt. The callback function is chained with the
// existing callback function.
func WitchChainedOnSalt(f OnSaltFunc) OptionFunc {
	return func(o *Options) {
		old := o.OnSalt
		o.OnSalt = func(salt []byte) {
			if old != nil {
				old(salt)
			}
			f(salt)
		}
	}
}

func (o *Options) onSalt(salt []byte) {
	if o.OnSalt != nil {
		o.OnSalt(salt)
	}
}

// FillOptions fills the Argon2 options.
func FillOptions(out *Options, opts []OptionFunc) {
	*out = Options{Time: 1, Memory: 64 * 1024, Threads: 4, KeyLen: 32, SaltLen: 32}
	for _, opt := range opts {
		opt(out)
	}
}

// NewArgon2 returns a new Argon2 hash of the password.
func NewArgon2(password []byte, opts ...OptionFunc) (Argon2, error) {
	var o Options
	FillOptions(&o, opts)
	return NewArgon2WithOptions(password, &o)
}

func a2WithOptions(p []byte, o *Options) ([]byte, error) {
	salt := o.Salt
	if len(salt) == 0 {
		salt = make([]byte, o.SaltLen)
		if _, err := rand.Read(o.Salt); err != nil {
			return nil, err
		}
	}
	o.onSalt(salt)
	return argon2.IDKey(p, salt, o.Time, o.Memory, o.Threads, o.KeyLen), nil
}

// NewArgon2WithOptions returns a new Argon2 hash of the password with the given salt.
func NewArgon2WithOptions(password []byte, o *Options) (Argon2, error) {
	a, err := a2WithOptions(password, o)
	return Argon2(a), err
}
