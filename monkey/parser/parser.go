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
    Lowest
    Equals      // ==
    LessGreater // > or <
    Sum         // + -
    Product     // * /
    Prefix      // -X or !X
    Call        // myFunction(X)
)

var precedences = map[string] int {
    token.Eq:       Equals,
    token.NotEq:    Equals,
    token.Lt:       LessGreater,
    token.Gt:       LessGreater,
    token.Plus:     Sum,
    token.Minus:    Sum,
    token.Slash:    Product,
    token.Asterisk: Product,
    token.Lparen:   Call,
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
    return this.curr.Type != token.Eof
}

func (this *Parser) currPrecedence() int {
    return precedences[this.curr.Type]
}

func (this *Parser) peekPrecedence() int {
    return precedences[this.peek.Type]
}

func (this *Parser) Errors() []string {
    return this.errors
}

func (this *Parser) PrintErrors() {
    for i, err := range this.errors {
        fmt.Printf("[%d] %s\n", i, err)
    }
}

func (this *Parser) addTokenError(tokenType string) {
    err := fmt.Sprintf("Expected token to be %s but got %s instead", tokenType, this.curr.Type)
    this.errors = append(this.errors, err)
}

func (this *Parser) addError(msg string) {
    this.errors = append(this.errors, msg)
}

func (this *Parser) parseLetStatement() *ast.LetStatement {
    // Start: Curr is token.LET
    var stm = &ast.LetStatement {}
    hasError := false

    this.next() // Jumps to the token.IDENT

    // Check identifier
    if !this.isCurr(token.Ident) {
        this.addTokenError(token.Ident)
        hasError = true
    }
    if !hasError {
        stm.Identifier = this.curr.Literal
        this.next() // Jumps to token.ASSIGN
    }

    // Check for asign symbol
    if !this.isCurr(token.Assign) {
        this.addTokenError(token.Assign)
        hasError = true
    }
    this.next() // Jumps to the first token of the expression

    stm.Expression = this.parseExpression(Lowest)
    this.next() // Jumps to the token.SEMICOLON

    if hasError { return nil }

    return stm
}

func (this *Parser) parseReturnStatement() *ast.ReturnStatement {
    // Start: Curr is token.RETURN
    var stm = &ast.ReturnStatement {}

    this.next()

    stm.Expression = this.parseExpression(Lowest)
    this.next() // Jumps to the token.SEMICOLON

    return stm
}

func (this *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
    var indexExpr = &ast.IndexExpression {}
    indexExpr.Left = left
    this.next() // Jumps to the token.Lbracket
    this.next() // Jumps inside the brackets so the expr is not viewed as an array
    indexExpr.Index = this.parseExpression(Lowest)
    this.next() // Jumps to the token.Rbracket
    return indexExpr
}

func (this *Parser) parseMethodExpression(expr ast.Expression) ast.Expression {
    var methExpr = &ast.MethodExpression {}

    // Left value
    methExpr.Expression = expr

    // Separator
    this.next() // Jumps to the token.Dot

    // Right value
    this.next() // Jumps to the first token of the call expression
    var right = this.parseExpression(Lowest)
    var callExpr, ok = right.(*ast.CallExpression)
    if !ok {
        this.addError("The right value of the method expression is not a call expression")
        return nil
    }
    methExpr.Call = callExpr

    return methExpr
}

func (this *Parser) parsePrefixOrSymbol() ast.Expression {
    switch this.curr.Type {
    case token.Bang, token.Minus:
        var pre = &ast.PrefixExpression {}
        pre.Operator = this.curr.Literal
        this.next()
        pre.Value = this.parseExpression(Prefix)
        return pre
    case token.True, token.False:
        return &ast.Boolean { Value: this.isCurr(token.True) } // Easy convert to bool trick :D
    case token.Ident:
        var identifier = &ast.Identifier { Value: this.curr.Literal }

        if this.isPeek(token.Lbracket) {
            return this.parseIndexExpression(identifier)
        }

        if this.isPeek(token.Dot) {
            return this.parseMethodExpression(identifier)
        }

        return identifier

    case token.Int:
        var intValue, err = strconv.ParseInt(this.curr.Literal, 10, 64)
        if err != nil {
            this.addError("Could not convert current token literal to int64")
        }
        return &ast.IntegerLiteral { Value: intValue }
    case token.String:
        return &ast.StringLiteral { Value: this.curr.Literal }
    case token.Lparen:
        this.next() // Jumps the token.LPAREN
        var exp = this.parseExpression(Lowest)
        if !this.isPeek(token.Rparen) {
            this.addError("Grouped expression did not end with and token.RPAREN")
            return nil
        }
        this.next() // Jumps the token.RPAREN
        return exp
    case token.Lbracket:
        var array = &ast.ArrayLiteral {}
        array.Elements = []ast.Expression {}

        this.next()
        for !this.isCurr(token.Rbracket) {
            var elem = this.parseExpression(Lowest)
            array.Elements = append(array.Elements, elem)
            this.next()
            if this.isCurr(token.Comma) { this.next() }
        }

        if !this.isPeek(token.Lbracket) {
            return array
        }

        return this.parseIndexExpression(array)

    case token.Lbrace:
        var hash = &ast.HashLiteral {}
        hash.Pairs = make(map[ast.Expression] ast.Expression)

        this.next() // Jumps to the first token inside the hash or the token.Rbrace

        for !this.isCurr(token.Rbrace) {
            var key = this.parseExpression(Lowest)

            if !this.isPeek(token.Colon) {
                this.addError(fmt.Sprintf("Expected ':' after hash key but found '%s' instead", this.curr.Literal))
                return nil
            }

            this.next() // Jumps to the token.Colon
            this.next() // Jumps to the first token of the pair value side

            var value = this.parseExpression(Lowest)

            hash.Pairs[key] = value

            this.next() // Normal iteration
            if this.isCurr(token.Comma) { this.next() } // If has next pair
        }

        return hash

    case token.If:
        return this.parseIfExpression()

    case token.Function:
        return this.parseFunctionLiteral()

    default:
        this.addError("Invalid or not covered symbol or prefix to parse: " + this.curr.Type)
        return nil
    }
}

func (this *Parser) parseIfExpression() ast.Expression {
    // Start: Curr is token.IF
    if !this.isPeek(token.Lparen) {
        this.addError("Expected token.LPAREN but got " + this.peek.Type + " instead")
        return nil
    }
    this.next() // Jumps to token.LPAREN

    this.next() // Jumps to first token in the condition

    var exp = &ast.IfExpression {}
    exp.Condition = this.parseExpression(Lowest)

    if !this.isPeek(token.Rparen) {
        this.addError("Expected token.RPAREN but got " + this.peek.Type + " instead")
        return nil
    }
    this.next() // Jumps to token.RPAREN

    if !this.isPeek(token.Lbrace) {
        this.addError("Expected token.LBRACE but got " + this.peek.Type + " instead")
        return nil
    }
    this.next() // Jumps to token.LBRACE

    this.next() // Jumps to the first token in the consequence block

    var consequences = []ast.Statement {}
    for !this.isCurr(token.Rbrace) && !this.isCurr(token.Eof) {
        var stm = this.parseStatement()
        consequences = append(consequences, stm)
        if this.isCurr(token.Semicolon) { this.next() } // Jumps the semicolon
    }
    exp.ConsequenceBlock = &ast.StatementsBlock { Statements: consequences }

    if !this.isPeek(token.Else) { return exp }

    this.next() // Jumps to token.ELSE
    this.next() // Jumps to token.LBRACE
    this.next() // Jumps to the first token in the alternative block

    var alternatives = []ast.Statement {}
    for !this.isCurr(token.Rbrace) && !this.isCurr(token.Eof) {
        var stm = this.parseStatement()
        alternatives = append(alternatives, stm)
        if this.isCurr(token.Semicolon) { this.next() } // Jumps the semicolon
    }
    exp.AlternativeBlock = &ast.StatementsBlock { Statements: alternatives }

    return exp
}

func (this *Parser) parseFunctionLiteral() ast.Expression {
    // Start: Curr is token.FUNCTION
    if !this.isPeek(token.Lparen) {
        this.addError("Expected token.LPAREN but got " + this.peek.Type + " instead")
        return nil
    }
    this.next() // Jumps to the token.LPAREN

    var funLiteral = &ast.FunctionLiteral {}
    funLiteral.Parameters = []ast.Identifier {}

    this.next() // Jumps to the first token of the function arguments or the right paren if none

    for !this.isCurr(token.Rparen) { // Parse function args
        var iden = ast.Identifier { Value: this.curr.Literal }
        funLiteral.Parameters = append(funLiteral.Parameters, iden)
        this.next()
        if this.isCurr(token.Comma) { this.next() }
    }

    if !this.isPeek(token.Lbrace) {
        this.addError("Expected token.LBRACE but got " + this.peek.Type + " instead")
        return nil
    }
    this.next() // Jumps to the token.LBRACE

    this.next() // Jumps to the first token in the function body

    var body = []ast.Statement {}
    for !this.isCurr(token.Rbrace) {
        var stm = this.parseStatement()
        body = append(body, stm)
        if this.isCurr(token.Semicolon) { this.next() } // Jumps the semicolon
    }
    funLiteral.Body = &ast.StatementsBlock { Statements: body }

    return funLiteral
}

func (this *Parser) makeInfix(left ast.Expression) ast.Expression {
    // Start: Curr is operator
    var inf = &ast.InfixExpression {}
    inf.Left = left
    inf.Operator = this.curr.Literal
    var precedence = this.currPrecedence()
    this.next() // Curr to next value
    inf.Right = this.createNewInfixGroup(precedence)
    return inf
}

func (this *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
    // Start: Curr is token.LPAREN
    var callExp = &ast.CallExpression {}
    callExp.Expression = fn

    this.next() // Jumps to token.RPAREN or to first token of parameters

    callExp.Parameters = []ast.Expression {}
    for !this.isCurr(token.Rparen) {
        var exp = this.parseExpression(Lowest)
        callExp.Parameters = append(callExp.Parameters, exp)
        this.next()
        if this.isCurr(token.Comma) { this.next() }
    }

    return callExp
}

func (this *Parser) parseInfix(expression ast.Expression) ast.Expression {
    switch this.curr.Type {
    case token.Plus, token.Minus, token.Slash, token.Asterisk, token.Eq, token.NotEq, token.Lt, token.Gt:
        return this.makeInfix(expression)
    case token.Lparen:
        return this.parseCallExpression(expression)
    default:
        this.addError("Invalid or not covered symbol for infix parse: " + this.curr.Type)
        return nil
    }
}

/// 1. add to left every time precedence is the same or lower
/// 2. create a new group every time it goes up
/// 3. close the group when it goes lower again
func (this *Parser) createNewInfixGroup(ctxPrecedence int) ast.Expression {
    var parsedValue = this.parsePrefixOrSymbol()

    var acc = parsedValue

    for !this.isPeek(token.Semicolon) && this.peekPrecedence() > ctxPrecedence {
        this.next() // Curr to operator
        acc = this.parseInfix(acc)
    }

    return acc
}

func (this *Parser) parseExpression(precedence int) ast.Expression {
    return this.createNewInfixGroup(precedence)
}

func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stm := &ast.ExpressionStatement {}
    stm.Expression = this.parseExpression(Lowest)
    this.next()
    return stm
}

func (this *Parser) parseStatement() ast.Statement {
    switch this.curr.Type {
    case token.Let:
        return this.parseLetStatement()
    case token.Return:
        return this.parseReturnStatement()
    case token.Illegal:
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
            this.addError("Statement is nill. Something went wrong!!!")
            return nil
        }

        if !this.isCurr(token.Semicolon) && !this.isCurr(token.Eof) {
            this.addError("The statement did not end with a semicolon")
            return nil
        }

        program.Statements = append(program.Statements, stm)
        this.next() // Jumps the semicolon
    }
    return program
}
