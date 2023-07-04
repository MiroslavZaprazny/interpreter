package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

const (
    _ int = iota
    LOWEST
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

type Parser struct {
    l *lexer.Lexer
    curToken token.Token
    peekToken token.Token
    errors []string
    prefixParserFns map[token.TokenType]prefixParserFn
    infixParserFns map[token.TokenType]infixParserFn
}

type (
    prefixParserFn func() ast.Expression
    infixParserFn func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l: l,
        errors: []string{},
    }

    p.prefixParserFns = make(map[token.TokenType]prefixParserFn)
    p.registerPrefix(token.IDENT, p.parseIdenfier)

    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) parseIdenfier() ast.Expression {
    return &ast.Indentifier{Token: p.curToken, Value: p.curToken.Literal}
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
        return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: p.curToken} 

    stmt.Expression = p.parseExpression(LOWEST)

    if p.peekToken.Type == token.SEMICOLON {
        p.nextToken()
    }

    return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
    prefix := p.prefixParserFns[p.curToken.Type]

    if prefix == nil {
        return nil
    }

    leftExp := prefix()

    return leftExp
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParserFn) {
    p.prefixParserFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParserFn) {
    p.infixParserFns[tokenType] = fn
}
