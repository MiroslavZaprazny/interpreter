package ast

import "interpreter/token"

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string {
    if len(p.Statements) < 0 {
        return ""
    }

    return p.Statements[0].TokenLiteral()
}

type Indentifier struct {
    Token token.Token
    Value string
}

func (i *Indentifier) expressionNode()
func (i *Indentifier) TokenLiteral() string {return i.Token.Literal}

type LetStatemet struct {
    Token token.Token
    Name *Indentifier
    Value Expression
}

func (ls *LetStatemet) statementNode()
func (ls *LetStatemet) TokenLiteral() string {return ls.Token.Literal}

