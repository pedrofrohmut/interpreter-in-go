// monkey/parser/parser.go

package parser

import (
    "fmt"
    "strconv"
    "monkey/ast"
    "monkey/lexer"
    "monkey/token"
    "monkey/utils"
)

const (
    // iota gives the constants a ascending numbers
    // _ skips the 0 value
    _ int = iota
    LOWEST
    EQUALS      // ==
    LESSGREATER // > or <
    SUM         // + -
    PRODUCT     // * /
    PREFIX      // -X or !X
    CALL        // myFunction(X)
    HIGHEST
)

var precedences = map[string] int {
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
}

type Parser struct {
    lex *lexer.Lexer
    curr token.Token
    peek token.Token
    errors []string
}

func NewParser(lexer *lexer.Lexer) *Parser {
    parser := &Parser { lex: lexer }

    parser.curr = lexer.GetNextToken()
    parser.peek = lexer.GetNextToken()

    parser.errors = []string {}

    return parser
}

func (this *Parser) isCurr(tokenType string) bool {
    return this.curr.Type == tokenType
}

func (this *Parser) isPeek(tokenType string) bool {
    return this.peek.Type == tokenType
}

func (this *Parser) next() {
    this.curr = this.peek
    this.peek = this.lex.GetNextToken()
}

func (this *Parser) hasNext() bool {
    return this.curr.Type != token.EOF
}

func (this *Parser) currPrecedence() int {
    return precedences[this.curr.Type]
}

func (this *Parser) peekPrecedence() int {
    return precedences[this.peek.Type]
}

func (this *Parser) addTokenError(tokenType string) {
    err := fmt.Sprintf("Expected token to be %s but got %s instead", tokenType, this.curr.Type)
    this.errors = append(this.errors, err)
}

func (this *Parser) parseLetStatement() *ast.LetStatement {
    if !this.isCurr(token.LET) { return nil } // BEFORE

    stm := ast.NewLetStatement()
    hasError := false

    // Check identifier
    this.next()
    if !this.isCurr(token.IDENT) {
        this.addTokenError(token.IDENT)
        hasError = true
    }
    if !hasError {
        stm.Identifier = this.curr.Literal
        this.next()
    }

    // Check for asign symbol
    if !this.isCurr(token.ASSIGN) {
        this.addTokenError(token.ASSIGN)
        hasError = true
    }
    this.next()

    // TODO: parse the expression later
    for !this.isCurr(token.SEMICOLON) { this.next() }

    if hasError { return nil } // AFTER

    return stm // Parse should end with curr == token.SEMICOLON
}

func (this *Parser) parseReturnStatement() *ast.ReturnStatement {
    if !this.isCurr(token.RETURN) { return nil } // BEFORE

    stm := ast.NewReturnStatement()

    this.next()

    // TODO: parse the expression later
    for !this.isCurr(token.SEMICOLON) { this.next() }

    return stm
}

func (this *Parser) parseIdentifierExpression() ast.Expression {
    return ast.NewIdentifier(this.curr.Literal)
}

func (this *Parser) parseIntegerLiteralExpression() ast.Expression {
    intValue, err := strconv.ParseInt(this.curr.Literal, 10, 64)
    if err != nil {
        this.errors = append(this.errors, "Could not convert current token literal to int64")
    }
    return ast.NewIntegerLiteral(intValue)
}

func (this *Parser) parsePrefixExpression() ast.Expression {
    value, err := strconv.ParseInt(this.peek.Literal, 10, 64)
    if err != nil {
        this.errors = append(this.errors, "Could not convert peek token literal to int64")
    }
    operator := this.curr.Literal
    this.next()
    return ast.NewPrefixExpression(operator, value)
}

func (this *Parser) parsePrefixOrSymbol() ast.Expression {
    switch this.curr.Type {
    case token.BANG, token.MINUS:
        value, err := strconv.ParseInt(this.peek.Literal, 10, 64)
        if err != nil {
            this.errors = append(this.errors, "Could not convert peek token literal to int64")
        }
        operator := this.curr.Literal
        this.next()
        return ast.NewPrefixExpression(operator, value)
    case token.IDENT:
        return ast.NewIdentifier(this.curr.Literal)
    case token.INT:
        intValue, err := strconv.ParseInt(this.curr.Literal, 10, 64)
        if err != nil {
            this.errors = append(this.errors, "Could not convert current token literal to int64")
        }
        return ast.NewIntegerLiteral(intValue)
    default:
        this.errors = append(this.errors, "Invalid symbol or prefix to parse: " + this.curr.Type)
        return nil
    }
}

// func (this *Parser) parseInfix(left ast.Expression) ast.Expression {
//     exp := &ast.InfixExpression { Left: left, Operator: this.curr.Literal }
//     this.next()
//     exp.Right = this.parseExpression()
//     return exp
// }

// func (this *Parser) parseInfix(left ast.Expression) ast.Expression {
//         exp := &ast.InfixExpression {}
//         exp.Left = left
//         exp.Operator = this.curr.Literal
//
//         currOperatorPrecedence := this.currPrecedence()
//         this.next()
//
//         if currOperatorPrecedence <= this.peekPrecedence() {
//             exp.Right = this.parseLeft()
//             left = exp
//         }
// }



/*
    before => curr == symbol or prefix + symbol
    1. convert symbol or prefix + symbol to expression => x
    2. recursion base case => return x if peek is token.SEMICOLON
    3. next to operator

    currVal := left
    currOpe := this.curr.Literal
    currPre := this.currPrecedence()

    exp := &ast.InfixExpression {}

    this.next() // curr -> B

    // Break Recursion on semicolon

    if currPre <= this.peekPrecedence() {
        // Left  <- Infix A + B
        exp.Left = this.parseInfix()
        // Right <- C
        exp.Right = currVal
    } else {
        // Left  <- A
        exp.Left = currVal
        // Right <- B * C
        exp.Right = this.parseInfix()
    }
*/
// func (this *Parser) parseExpression(value ast.Expression, operator string) ast.Expression {
//     x := this.parsePrefixOrSymbol()
//
//     // Recursion base case
//     if this.isPeek(token.SEMICOLON) { return x }
//
//     this.next() // After this: Curr should point to an operator token
//
//     // Debugging aliases
//     currVal := x
//     currPre := this.currPrecedence()
//     currOpe := this.curr.Literal
//
//     this.next() // After this: Curr should point to symbol next to the operator
//
//     // Recursion base case to the last Infix
//     if this.isPeek(token.SEMICOLON) {
//         exp := &ast.InfixExpression {}
//         exp.Left = x
//         exp.Operator = currOpe
//         exp.Right = this.parsePrefixOrSymbol()
//         return exp
//     }
//
//     currInf := &ast.InfixExpression{}
//     currInf.Operator = currOpe
//
//     if currPre <= this.peekPrecedence() {
//         exp := &ast.InfixExpression {}
//
//         exp.Left = this.parseExpression()
//
//         exp.Operator = this.curr.Literal
//         this.next()
//         exp.Right = this.parsePrefixOrSymbol()
//     } else {
//
//     }
// }

// func (this *Parser) newInfix(a ast.Expression, ope string, b ast.Expression) ast.Expression {
//     return &ast.InfixExpression { Left: a, Operator: ope, Right: b }
// }

/*
    With the same precedence for every operator:
    return (left, c) => ((left, c), d) => (((left, c), d), e)

    TODO: use precedence to ajust the expression tree
*/
func (this *Parser) parseInfix(acc ast.Expression) ast.Expression {
    currOpe := this.curr.Literal

    this.next()

    right := this.parsePrefixOrSymbol()

    inf := &ast.InfixExpression {}
    inf.Left = acc
    inf.Operator = currOpe
    inf.Right = right

    // Base case to break recursion
    if this.isPeek(token.SEMICOLON) { return inf }

    this.next()

    return this.parseInfix(inf)
}

func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stm := &ast.ExpressionStatement {}

    first := this.parsePrefixOrSymbol()

    if this.isPeek(token.SEMICOLON) {
        stm.Expression = first
    } else {
        this.next()
        stm.Expression = this.parseInfix(first)
    }

    this.next()
    return stm
}


// func (this *Parser) parseExpression(precedence int) ast.Expression {
//     leftParseFn := this.leftParseFns[this.curr.Type]
//     if utils.IsNill(leftParseFn) {
//         this.errors = append(this.errors, "Left parse function not found for: " + this.curr.Type)
//         return nil
//     }
//     left := leftParseFn()
//
//     if this.isPeek(token.SEMICOLON) { return left }
//
//     if precedence <= this.peekPrecedence() { // a + b + c
//     } else { // a + b * c
//     }
//
//     // this.next()
//     //
//     // for !this.isCurr(token.SEMICOLON) {
//     //     infixParseFn := this.infixParseFns[this.curr.Type]
//     //     if utils.IsNill(infixParseFn) {
//     //         this.errors = append(this.errors, "ERROR: Infix parse function not found for: " + this.curr.Type)
//     //         return nil
//     //     }
//     //     left = infixParseFn(left)
//     //     this.next()
//     // }
//     //
//     // return left
// }

// func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
//     stm := ast.NewExpressionStatement()
//     stm.Expression = this.parseExpression(LOWEST)
//     if this.isPeek(token.SEMICOLON) { this.next() }
//     return stm // Parse should end with curr == token.SEMICOLON
// }

func (this *Parser) parseStatement() ast.Statement {
    switch this.curr.Type {
    case token.LET:
        return this.parseLetStatement()
    case token.RETURN:
        return this.parseReturnStatement()
    case token.ILLEGAL:
        return nil
    default:
        return this.parseExpressionStatement()
    }
}

func (this *Parser) ParseProgram() *ast.Program {
    program := ast.NewProgram()
    for this.hasNext() {
        stm := this.parseStatement()

        if utils.IsNill(stm) {
            this.errors = append(this.errors, "Statement is nill. Something went wrong!!!")
            return nil
        }

        if !this.isCurr(token.SEMICOLON) {
            this.errors = append(this.errors, "The statement did not end with a semicolon")
            return nil
        }

        program.Statements = append(program.Statements, stm)
        this.next() // Jumps the semicolon
    }
    return program
}
