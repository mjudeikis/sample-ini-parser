package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Emit puts a token onto the token channel. The value of this token is
// read from the input based on the current lexer position.
func (l *Lexer) Emit(tokenType TokenType) {
	l.Tokens <- Token{Type: tokenType, Value: l.Input[l.Start:l.Pos]}
	l.Start = l.Pos
}

// Inc increment the position
func (l *Lexer) Inc() {
	l.Pos++
	if l.Pos >= utf8.RuneCountInString(l.Input) {
		l.Emit(TOKEN_EOF)
	}
}

// InputToEnd return a slice of the input from the current lexer position
// to the end of the input string.
func (l *Lexer) InputToEnd() string {
	return l.Input[l.Pos:]
}

// Skips whitespace until we get something meaningful.
func (l *Lexer) SkipWhitespace() {
	for {
		ch := l.Next()

		if !unicode.IsSpace(ch) {
			l.Dec()
			break
		}

		if ch == EOF {
			l.Emit(TOKEN_EOF)
			break
		}
	}
}

// Dec decrements the possition
func (l *Lexer) Dec() {
	l.Pos--
}

// Next Reads the next rune (character) from the input stream
// and advances the lexer position.
func (l *Lexer) Next() rune {
	if l.Pos >= utf8.RuneCountInString(l.Input) {
		l.Width = 0
		return EOF
	}

	result, width := utf8.DecodeRuneInString(l.Input[l.Pos:])

	l.Width = width
	l.Pos += l.Width
	return result
}

// IsEOF returns the true/false if the lexer is at the end of the input stream.
func (l *Lexer) IsEOF() bool {
	return l.Pos >= len(l.Input)
}

// IsWhitespace eturns true/false if then next character is whitespace
func (l *Lexer) IsWhitespace() bool {
	ch, _ := utf8.DecodeRuneInString(l.Input[l.Pos:])
	return unicode.IsSpace(ch)
}

// BeginLexing start a new lexer with a given input string. This returns the
// instance of the lexer and a channel of tokens. Reading this stream
// is the way to parse a given input and perform processing.
func BeginLexing(name, input string) *Lexer {
	l := &Lexer{
		Name:   name,
		Input:  input,
		State:  LexBegin,
		Tokens: make(chan Token, 3),
	}

	return l
}

// LexBegin lexer function starts everything off. It determines if we are
// beginning with a key/value assignment or a section.
func LexBegin(l *Lexer) LexFn {
	l.SkipWhitespace()

	if strings.HasPrefix(l.InputToEnd(), LEFT_BRACKET) {
		return LexLeftBracket
	} else {
		return LexKey
	}
}

// LexLeftBracket lexer function emits a TOKEN_LEFT_BRACKET then returns
// the lexer for a section header.
func LexLeftBracket(l *Lexer) LexFn {
	l.Pos += len(LEFT_BRACKET)
	l.Emit(TOKEN_LEFT_BRACKET)
	return LexSection
}

// LexSection lexer function emits a TOKEN_SECTION with the name of an
// INI file section header.
func LexSection(l *Lexer) LexFn {
	for {
		if l.IsEOF() {
			return l.Errorf(LEXER_ERROR_MISSING_RIGHT_BRACKET)
		}

		if strings.HasPrefix(l.InputToEnd(), RIGHT_BRACKET) {
			l.Emit(TOKEN_SECTION)
			return LexRightBracket
		}

		l.Inc()
	}
}

// LexRightBracket  lexer function emits a TOKEN_RIGHT_BRACKET then returns
// the lexer for a begin.
func LexRightBracket(l *Lexer) LexFn {
	l.Pos += len(RIGHT_BRACKET)
	l.Emit(TOKEN_RIGHT_BRACKET)
	return LexBegin
}

// LexKey lexer function emits a TOKEN_KEY with the name of an
// key that will be assigned a value.
func LexKey(l *Lexer) LexFn {
	for {
		if strings.HasPrefix(l.InputToEnd(), EQUAL_SIGN) {
			l.Emit(TOKEN_KEY)
			return LexEqualSign
		}

		l.Inc()

		if l.IsEOF() {
			return l.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

// LexEqualSign lexer function emits a TOKEN_EQUAL_SIGN then returns
// the lexer for value.
func LexEqualSign(l *Lexer) LexFn {
	l.Pos += len(EQUAL_SIGN)
	l.Emit(TOKEN_EQUAL_SIGN)
	return LexValue
}

// LexValue lexer function emits a TOKEN_VALUE with the value to be assigned
// to a key.
func LexValue(l *Lexer) LexFn {
	for {
		if strings.HasPrefix(l.InputToEnd(), NEWLINE) {
			l.Emit(TOKEN_VALUE)
			return LexBegin
		}

		l.Inc()

		if l.IsEOF() {
			return l.Errorf(LEXER_ERROR_UNEXPECTED_EOF)
		}
	}
}

func (l *Lexer) Errorf(format string, args ...interface{}) LexFn {
	l.Tokens <- Token{
		Type:  TOKEN_ERROR,
		Value: fmt.Sprintf(format, args...),
	}

	return nil
}

// NextToken return the next token from the channel

func (l *Lexer) NextToken() Token {
	for {
		select {
		case token := <-l.Tokens:
			return token
		default:
			l.State = l.State(l)
		}
	}

	panic("Lexer.NextToken reached an invalid state!!")
}
