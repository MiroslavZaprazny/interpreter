package token

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}

const (
    ILLEGAL = "ILLEGAL"
    EOF = "EOF"
    IDENT = "IDENT"
    INT = "INT"
    STRING = "STRING"

    ASSIGN = "="
    PLUS = "+"
    MINUS = "-"
    BANG = "!"
    ASTERISK = "*"
    SLASH = "/"
    LT = "<"
    GT = ">"
    COMMA = ","
    SEMICOLON = ";"
    LPAREN = "("
    RPAREN = ")"
    LBRACE  = "{"
    RBRACE  = "}"

    FUNCTION = "FUNCTION"
    LET = "LET"
    IF = "IF"
    ELSE = "ELSE"
    RETURN = "RETURN" 
    TRUE = "TRUE"
    FALSE = "FALSE"
    
    EQ = "=="
    NOT_EQ = "!="
)

var keywords = map[string] TokenType {
    "fn": FUNCTION,
    "let": LET,
    "if": IF,
    "else": ELSE,
    "return": RETURN,
    "true": TRUE,
    "false": FALSE,
}

func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }

    return IDENT
}
