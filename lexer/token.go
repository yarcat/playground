package tokenizer

import "fmt"

type (
	// TokenType is a type for token type enumerator. Custom token types should
	// start with TokenCustom+1.
	TokenType int

	// Token represents a group of characters having collective meaning (e.g.
	// numbers or identifiers) separated by a lexical scanner.
	Token struct {
		// Type represents current token type.
		Type TokenType
		// Value is a token value. For some tokens (e.g. TokenEOF) the value may
		// be unset. If Type is TokenError, the value contains an error message.
		Value string
	}
)

const (
	// TokenUndefined represents an unknown token type. Its meant to be assigned
	// as a zero-value, and it shouldn't be used.
	TokenUndefined TokenType = iota
	// TokenError represents a token that contains an error message.
	TokenError
	// TokenEOF is generated on attempt to read past the input string.
	TokenEOF

	// TokenCustom is a meta token type, which is meant to be used for custom
	// token types, which should start from TokenCustom+1.
	TokenCustom TokenType = 1000
)

var tokenTypeNames = map[TokenType]string{
	TokenUndefined: "UNDEFINED",
	TokenError:     "ERROR",
	TokenEOF:       "EOF",
}

// String returns a human-readable token type.
func (tt TokenType) String() string {
	name, ok := tokenTypeNames[tt]
	if !ok {
		return "???"
	}
	return name
}

// String returns a human-readable token representation.
func (tok Token) String() string {
	if len(tok.Value) == 0 {
		return fmt.Sprint(tok.Type)
	}
	return fmt.Sprintf("%s %q", tok.Type, tok.Value)
}

// AddTokenTypeNames adds custom token type mapping for producing human-readable
// token type strings.
func AddTokenTypeNames(names map[TokenType]string) {
	for tt, name := range names {
		tokenTypeNames[tt] = name
	}
}
