package parser

import (
	"log"
	"strings"

	"github.com/mjudeikis/sample-ini-parser/pkg/lexer"
)

func isEOF(token lexer.Token) bool {
	return token.Type == lexer.TOKEN_EOF
}

func Parse(fileName, input string) IniFile {
	output := IniFile{
		FileName: fileName,
		Sections: make([]IniSection, 0),
	}

	var token lexer.Token
	var tokenValue string

	/* State variables */
	section := IniSection{}
	key := ""

	log.Println("Starting lexer and parser for file", fileName, "...")

	l := lexer.BeginLexing(fileName, input)

	for {
		token = l.NextToken()

		if token.Type != lexer.TOKEN_VALUE {
			tokenValue = strings.TrimSpace(token.Value)
		} else {
			tokenValue = token.Value
		}

		if isEOF(token) {
			output.Sections = append(output.Sections, section)
			break
		}

		switch token.Type {
		case lexer.TOKEN_SECTION:
			// Reset tracking variables
			if len(section.KeyValuePairs) > 0 {
				output.Sections = append(output.Sections, section)
			}

			key = ""

			section.Name = tokenValue
			section.KeyValuePairs = make([]IniKeyValue, 0)

		case lexer.TOKEN_KEY:
			key = tokenValue

		case lexer.TOKEN_VALUE:
			section.KeyValuePairs = append(section.KeyValuePairs, IniKeyValue{Key: key, Value: tokenValue})
			key = ""
		}
	}

	log.Println("Parser has been shutdown")
	return output
}
