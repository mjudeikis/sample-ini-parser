package lexer

// TokenType is int defined by iota
type TokenType int

// LexFn is recursive function to return lexer function based on token
type LexFn func(*Lexer) LexFn

// Define TokenTypes for ini file
const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_LEFT_BRACKET
	TOKEN_RIGHT_BRACKET
	TOKEN_EQUAL_SIGN
	TOKEN_NEWLINE

	TOKEN_SECTION
	TOKEN_KEY
	TOKEN_VALUE
)

const EOF rune = 0

// Define identifiers for the file we are parsing
const LEFT_BRACKET string = "["
const RIGHT_BRACKET string = "]"
const EQUAL_SIGN string = "="
const NEWLINE string = "\n"

// Token defines one of the tokens
type Token struct {
	Type  TokenType
	Value string
}

// Lexer defines main lexer object
type Lexer struct {
	Name   string
	Input  string
	Tokens chan Token
	State  LexFn

	Start int
	Pos   int
	Width int
}
