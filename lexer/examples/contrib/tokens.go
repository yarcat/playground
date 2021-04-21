package contrib

import "github.com/yarcat/playground/tokenizer"

const (
	TokenContribNumber tokenizer.TokenType = tokenizer.TokenCustom + 1 + iota
	TokenContribOperator
	TokenContribOpenParen
	TokenContribClosingParen

	TokenContribLast
)

func init() {
	tokenizer.AddTokenTypeNames(map[tokenizer.TokenType]string{
		TokenContribNumber:       "NUMBER",
		TokenContribOperator:     "OPERATOR",
		TokenContribOpenParen:    "OPEN_PAREN",
		TokenContribClosingParen: "CLOSING_PAREN",
	})
}
