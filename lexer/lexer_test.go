package lexer

import (
	"interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
    input := `=+(){},;`

    tests := []struct {
        expectedType token.TokenType
        expectedLiteral string
    }{
        {token.ASSIGN, "="},
        {token.PLUS, "+"},
        {token.LPAREN, "("},
        {token.RPAREN, ")"},
        {token.LBRACE, "{"},
        {token.RBRACE, "}"},
        {token.COMMA, ","},
        {token.SEMICOLON, ";"},
        {token.EOF, ""},
    }

    l:= New(input)

    for i, tt := range tests {
        print(i, tt)
        token := l.nextToken()
        if token.Type != tt.expectedType {
            t.Fatalf("expected %q, got %q", tt.expectedType, token.Type)
        }
    }
}
