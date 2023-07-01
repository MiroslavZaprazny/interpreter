package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

type Parser struct {
    l *lexer.Lexer
    curToken token.Token
    peekToken token.Token
    errors []string
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l: l,
        errors: []string{},
    }

    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram () *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for p.curToken.Type != token.EOF {
        stmt := p.parseStatement()

        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.nextToken()
    }

    return program
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.curToken.Type {
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return nil
    }
}

func (p *Parser) parseLetStatement() *ast.LetStatemet {
    stmt := &ast.LetStatemet{Token: p.curToken}
    
    if !p.expectPeek(token.IDENT) {
        return nil
    }
    
    stmt.Name = &ast.Indentifier{Token: p.curToken, Value: p.curToken.Literal}

    if !p.expectPeek(token.ASSIGN) {
        return nil 
    }

    for p.curToken.Type != token.SEMICOLON {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: p.curToken}
    p.nextToken()

    for p.curToken.Type != token.SEMICOLON {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekToken.Type != t {
        p.peekError(t)
        return false
    }
    p.nextToken()
    
    return true
}

func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf(
        "expected next token to be %s got %s instead", 
        t,
        p.peekToken.Type,
    ) 

    p.errors = append(p.errors, msg)
}
