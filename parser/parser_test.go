package parser

import (
	"fmt"
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

func TestParsingPrefixExpressions(t *testing.T) {
    prefixTests := []struct {
        intput string
        operator string
        integerValue int64
    }{
        {"!5", "!", 5},
        {"-15", "-", 15},
    }

    for _, tt := range prefixTests {
        l := lexer.New(tt.intput)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf(
                "program.statements does not cotain %d statements, got=%d",
                1, len(program.Statements),
            )
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "program.statements[0] is not ast.ExpressionStatement got=%T",
                program.Statements[0],
            )
        }

        exp, ok := stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("stmt is not ast.PrefixExpression, got=%T", stmt.Expression)
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s', got=%s", tt.operator, exp.Operator)
        }

        if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
            return
        }
    }

}

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input string
        leftValue int64
        operator string
        rightValue int64
    } {
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},
    }

    for _, tt := range infixTests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf(
                "program.Statements does not cotain %d statements, got=%d", 
                1, 
                len(program.Statements),
            )
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf(
                "program.Statements[0] is not ast.ExpressionStatement got=%T", 
                program.Statements[0],
            )
        }

        exp, ok := stmt.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("exp is not ast.InfixExpression, got=%T", stmt.Expression)
        }

        if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
            return
        }

        if exp.Operator != tt.operator {
            t.Fatalf(
                "exp.Operator os not %s got=%s",
                tt.operator,
                exp.Operator,
            )
        }

        if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
            return
        }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {"-a * b", "((-a) * b)"},
        {"!-a", "(!(-a))"},
        {"a + b + c", "((a + b) + c)"},
        {"a + b - c", "((a + b) - c)"},
        {"a * b * c","((a * b) * c)"},
        {"a * b / c", "((a * b) / c)"},
        {"a + b / c", "(a + (b / c))"},
        {"a + b * c + d / e - f","(((a + (b * c)) + (d / e)) - f)"},
        {"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
        {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
        {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
        {"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
        {"3 + 4 * 5 == 3 * 1 + 4 * 5","((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
    }

    for _, tt := range tests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        actual := program.String()
        if actual != tt.expected {
            t.Errorf("expected=%q, got=%q", tt.expected, actual)
        }
    }

}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
    integer, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il not ast.IntegerLiteral, got=%T", il)
        return false
    }

    if integer.Value != value {
        t.Errorf("integer.value not %d, got=%d", value, integer.Value)
        return false
    }

    if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
        t.Errorf("integer.TokenLiteral not %d, got=%s", value, integer.TokenLiteral())
        return false
    }

    return true
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
