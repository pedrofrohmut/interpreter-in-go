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

type (
    LeftParseFn func() ast.Expression
    InfixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
    lex *lexer.Lexer
    curr token.Token
    peek token.Token
    errors []string
    leftParseFns map[string]LeftParseFn
    InfixParseFns map[string]InfixParseFn
}

func NewParser(lexer *lexer.Lexer) *Parser {
    parser := &Parser { lex: lexer }
    parser.curr = lexer.GetNextToken()
    parser.peek = lexer.GetNextToken()

    parser.errors = []string {}

    parser.leftParseFns = make(map[string]LeftParseFn)
    parser.leftParseFns[token.IDENT] = parser.parseIdentifierExpression
    parser.leftParseFns[token.INT]   = parser.parseIntegerLiteralExpression
    parser.leftParseFns[token.BANG]  = parser.parsePrefixExpression
    parser.leftParseFns[token.MINUS] = parser.parsePrefixExpression

    return parser
}

func (this *Parser) isCurr(tokenType string) bool {
    return this.curr.Type == tokenType
}

func (this *Parser) next() {
    this.curr = this.peek
    this.peek = this.lex.GetNextToken()
}

func (this *Parser) hasNext() bool {
    return this.curr.Type != token.EOF
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

    if !this.isCurr(token.SEMICOLON) || hasError { return nil } // AFTER

    return stm // Parse should end with curr == token.SEMICOLON
}

func (this *Parser) parseReturnStatement() *ast.ReturnStatement {
    if !this.isCurr(token.RETURN) { return nil } // BEFORE

    stm := ast.NewReturnStatement()

    this.next()

    // TODO: parse the expression later
    for !this.isCurr(token.SEMICOLON) { this.next() }

    if !this.isCurr(token.SEMICOLON) { return nil } // AFTER

    return stm
}

func (this *Parser) parseExpression() ast.Expression {
    leftParseFn := this.leftParseFns[this.curr.Type]
    if utils.IsNill(leftParseFn) {
        this.errors = append(this.errors, "Left parse function not found for: " + this.curr.Type)
        return nil
    }
    return leftParseFn()
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

func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stm := ast.NewExpressionStatement()
    stm.Expression = this.parseExpression()
    this.next()

    if !this.isCurr(token.SEMICOLON) { return nil } // AFTER

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
        if !utils.IsNill(stm) {
            program.Statements = append(program.Statements, stm)
        }
        this.next()
    }
    return program
}
