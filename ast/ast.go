package ast

import (
	"bytes"
	"interpreter/token"
)

type Node interface {
    TokenLiteral() string
    String() string
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

func (p *Program) String() string {
    var out bytes.Buffer

    for _,s := range p.Statements {
        out.WriteString(s.String())
    }

    return out.String()
}

type Indentifier struct {
    Token token.Token //IDENT
    Value string
}

func (i *Indentifier) expressionNode() {}
func (i *Indentifier) TokenLiteral() string {return i.Token.Literal}
func (i *Indentifier) String() string {return i.Value}

type LetStatemet struct {
    Token token.Token //LET
    Name *Indentifier
    Value Expression
}

func (ls *LetStatemet) statementNode() {}
func (ls *LetStatemet) TokenLiteral() string {return ls.Token.Literal}
func (ls *LetStatemet) String() string {
    var out bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Name.String())
    out.WriteString(" = ")

    if ls.Value != nil {
        out.WriteString(ls.Value.String())
    }

    out.WriteString(";")

    return out.String()
}


type ReturnStatement struct {
    Token token.Token //RETURN
    ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {return rs.Token.Literal}
func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")

    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }

    out.WriteString(";")

    return out.String()
}


type ExpressionStatement struct {
    Token token.Token //first token of the expression
    Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {return es.Token.Literal}
func (es *ExpressionStatement) String() string {
    if es.Expression != nil {
        return es.Expression.String()
    }

    return ""
}
