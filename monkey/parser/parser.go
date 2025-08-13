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

func (this *Parser) addError(msg string) {
    this.errors = append(this.errors, msg)
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
    a + b + c + d
    return (a, b) => ((a, b), c) => (((a, b), c), d)

    TODO: use precedence to ajust the expression tree

    Every time you break precedence you start a new acc and Add it to the right
    until you break precedence lower again
    a + b * c - d
    return ((a + (b * c)) - d)
    if currPre >= peekPre {
        acc = parseInfix(acc)
    } else {
        right = parseInfix(rightValue)
    }
*/
func (this *Parser) OLD_parseInfix(acc ast.Expression) ast.Expression {
    // All parseInfix call should start with a operator token in the curr position
    if !token.IsOperator(this.curr) {
        this.addError("Parse infix is expected to start with an operator token but found '" +
            this.curr.Literal + "' instead")
        return nil
    }

    inf := &ast.InfixExpression {}
    inf.Left = acc
    inf.Operator = this.curr.Literal
    this.next()
    inf.Right = this.parsePrefixOrSymbol()

    if this.isPeek(token.SEMICOLON) {
        return inf
    }

    this.next() // jumps to next operator

    return this.OLD_parseInfix(inf)
}

// TODO: Change this method to have precedence to know when to break back
// func (this *Parser) parseInfix(acc ast.Expression) ast.Expression {
//     // All parseInfix call should start with a operator token in the curr position
//     if !token.IsOperator(this.curr) {
//         this.addError("Parse infix is expected to start with an operator token but found '" +
//             this.curr.Literal + "' instead")
//         return nil
//     }
//
//     currOpe := this.curr.Literal
//     currPre := this.currPrecedence()
//
//     this.next()
//     rightValue := this.parsePrefixOrSymbol()
//
//     if this.isPeek(token.SEMICOLON) {
//         inf := &ast.InfixExpression {}
//         inf.Left = acc
//         inf.Operator = currOpe
//         inf.Right = rightValue
//         return inf
//     }
//
//     fmt.Printf("Precedence Curr: %d, Peek: %d\n", currPre, this.peekPrecedence())
//
//     // If precedence is not going up keep accumulating on Infix.Left
//     if currPre >= this.peekPrecedence() {
//         inf := &ast.InfixExpression {}
//         inf.Left = acc
//         inf.Operator = currOpe
//         inf.Right = rightValue
//
//         // if this.isPeek(token.SEMICOLON) { return inf }
//
//         this.next() // goes to next operator
//
//         return this.parseInfix(inf)
//     } else {
//         newInf := &ast.InfixExpression {}
//         newInf.Left = rightValue
//         this.next()
//         newInf.Operator = this.curr.Literal
//         this.next()
//         newInf.Right = this.parsePrefixOrSymbol()
//
//         inf := &ast.InfixExpression {}
//         inf.Left = acc
//         inf.Operator = currOpe
//         inf.Right = newInf // ADD the right recursion here -> add left until precedence goes down
//         // inf.Right = this.parseInfix(newInf)
//
//         return inf
//     }
// }
        // newInf := &ast.InfixExpression {}
        // newInf.Left = rightValue
        // this.next()
        // newInf.Operator = this.curr.Literal
        // this.next()
        // newInf.Right = this.parsePrefixOrSymbol()
        //
        // inf := &ast.InfixExpression {}
        // inf.Left = acc
        // inf.Operator = currOpe
        // inf.Right = this.parseInfix(newInf)
        //
        // if this.isPeek(token.SEMICOLON) { return inf }
        //
        // this.next() // goes to next operator
        //
        // return this.parseInfix(inf)
    //
    // inf := &ast.InfixExpression {}
    // inf.Left = acc
    // inf.Operator = this.curr.Literal
    // this.next()
    // inf.Right = this.parsePrefixOrSymbol()
    //
    // if currOpe < this.peekPrecedence() || this.isPeek(token.SEMICOLON) {
    //     return inf
    // }
    //
    // this.next() // jumps to next operator
    //
    // return this.parseInfix(inf)

// func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
//     stm := &ast.ExpressionStatement {}
//
//     first := this.parsePrefixOrSymbol()
//
//     if this.isPeek(token.SEMICOLON) {
//         stm.Expression = first
//     } else {
//         this.next()
//         stm.Expression = this.parseInfix(first)
//     }
//
//     this.next()
//     return stm
// }

// func (this *Parser) OLDparseExpression() ast.Expression {
//     acc := this.parsePrefixOrSymbol()
//
//     for !this.isPeek(token.SEMICOLON) {
//         this.next() // Jump to operator
//         inf := &ast.InfixExpression {}
//         inf.Left = acc
//         inf.Operator = this.curr.Literal
//         this.next() // Jump to right symbol or prefix expression
//         inf.Right = this.parsePrefixOrSymbol()
//         acc = inf
//     }
//
//     return acc
// }

func (this *Parser) parseInfix(left ast.Expression, precedence int) ast.Expression {
    inf := &ast.InfixExpression {}

    inf.Left = left

    inf.Operator = this.curr.Literal

    // TODO: Check currPrecedence vs argument Precedence to decide what to do in the inf.Right

    this.next()
    inf.Right = this.parsePrefixOrSymbol()

    return inf
}

func (this *Parser) parseExpression(precedence int) ast.Expression {
    value := this.parsePrefixOrSymbol()

    if this.isPeek(token.SEMICOLON) { return value }

    acc := value

    for {
        this.next() // Move curr to operator

        acc = this.parseInfix(acc, precedence)
    }

    return acc
}

/// From: a + b * c To: (a + (b * c))
/// from: a + b * c * d / e + f + g * h
/// to:   ((a + (((b * c) * d) / e) + f) + (g * h))
/// 1. add to left every time precedence is the same or lower
/// 2. create a new group every time it goes up
/// 3. close the group when it goes lower again
func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stm := &ast.ExpressionStatement {}
    stm.Expression = this.parseExpression(HIGHEST)
    this.next()
    return stm
}

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
