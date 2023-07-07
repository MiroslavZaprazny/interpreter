package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetSatements(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let foobar = 838383;
    `

    l := lexer.New(input)
    p := New(l)
    
    program := p.ParseProgram()
    checkParserErrors(t, p)
    if program == nil {
        t.Fatalf("ParseProgram returned nil")
    }

    if len(program.Statements) != 3 {
        t.Fatalf(
            "program.Statements does not contain 3 statemens, got=%d", 
            len(program.Statements),
        )
    }

    tests := []struct {
        expectedIdentifier string
    }{
        {"x"},
        {"y"},
        {"foobar"},
    }

    for i, tt := range tests {
        stmt := program.Statements[i]
        if !testLetStatement(t, stmt, tt.expectedIdentifier) {
            return
        }
    }
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
   if s.TokenLiteral() != "let" {
       t.Errorf("s.TokenLiteral not 'let' got=%T", s)

       return false
   } 

   letStmt, ok := s.(*ast.LetStatemet)

   if !ok {
        t.Errorf("s not *ast.Statement got=%T", s)

        return false
   }

   if letStmt.Name.Value != name {
       t.Errorf("letStmt.Name.Value not %s got %s", name, letStmt.Name.Value)
   }

   if letStmt.Name.TokenLiteral() != name {
       t.Errorf("s nme not %s got %s", name, letStmt.Name)
   }

   return false
}

func TestReturnStatements(t *testing.T) {
    input := `
        return 5;
        return 10;
        return 92231;
    `

    l:= lexer.New(input)
    p:= New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 3 {
        t.Fatalf(
            "expected to cointain 3 statements got %d", 
            len(program.Statements),
        )
    }

    for _, stmt := range program.Statements {
        returnStmt, ok := stmt.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("stmt not a return statement got %d", stmt)
            continue
        }

        if returnStmt.TokenLiteral() != "return" {
            t.Errorf("stmt literal not 'return' got %s", returnStmt.TokenLiteral())
        }
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statemetnts. got=%d", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

    if !ok {
        t.Fatalf("statemnt is not an expression, got=%T", program.Statements[0])
    }

    ident, ok := stmt.Expression.(*ast.Indentifier)

    if !ok {
        t.Fatalf("expression not an indentifier got=%T", stmt.Expression)
    }

    if ident.Value != "foobar" {
        t.Errorf("expecter identifier foobar got=%s", ident.Value)
    }

    if ident.TokenLiteral() != "foobar" {
        t.Errorf("expectet token literal foobar got=%s", ident.Value)
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statemetnts. got=%d", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

    if !ok {
        t.Fatalf("statemnt is not an expression, got=%T", program.Statements[0])
    }

    literal, ok := stmt.Expression.(*ast.IntegerLiteral)

    if !ok {
        t.Fatalf("expression not an indentifier got=%T", stmt.Expression)
    }

    if literal.Value != 5 {
        t.Errorf("expecter identifier foobar got=%d", literal.Value)
    }

    if literal.TokenLiteral() != "5" {
        t.Errorf("expectet token literal foobar got=%s", literal.TokenLiteral())
    }
}

func checkParserErrors(t *testing.T, p *Parser) {
    errors := p.Errors() 

    if len(errors) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errors))

    for _, msg := range errors {
        t.Errorf("parser errorr %q", msg)
    }

    t.FailNow()
}
