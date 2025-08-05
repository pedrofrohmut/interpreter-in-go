// monkey/parser/parser.go
/*
     A parser is a software that takes an input and builds up an AST (Abstract
   Syntac Tree) that is structural representation of the input. The parser
   often uses tokens created from a lexer.
*/

package parser

import (
    "fmt"
    "monkey/ast"
    "monkey/lexer"
    "monkey/token"
    "monkey/utils"
    _ "monkey/utils"
    "strconv"
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
)

var precedences = map[token.TokenType] int {
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
}

type (
    PrefixParseFn func() ast.Expression
    InfixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
    lex *lexer.Lexer
    errors []string
    currToken token.Token
    peekToken token.Token
    prefixParseFns map[token.TokenType]PrefixParseFn
    infixParseFns map[token.TokenType]InfixParseFn
}

func NewParser(lex *lexer.Lexer) *Parser {
    par := &Parser {
        lex: lex,
        errors: []string {},
    }

    // Initialize tokens
    par.currToken = lex.GetNextToken()
    par.peekToken = lex.GetNextToken()

    // Register Prefix Functions
    par.prefixParseFns = make(map[token.TokenType]PrefixParseFn)
    par.registerPrefix(token.IDENT, par.parseIdentifierPrefix)
    par.registerPrefix(token.INT,   par.parseIntegerLiteral)
    par.registerPrefix(token.BANG,  par.parsePrefixExpression)
    par.registerPrefix(token.MINUS, par.parsePrefixExpression)
    par.registerPrefix(token.TRUE,  par.parseBoolean)
    par.registerPrefix(token.FALSE, par.parseBoolean)

    // Register Infix Functions
    par.infixParseFns = make(map[token.TokenType]InfixParseFn)
    par.registerInfix(token.PLUS,     par.parseInfixExpression)
    par.registerInfix(token.MINUS,    par.parseInfixExpression)
    par.registerInfix(token.SLASH,    par.parseInfixExpression)
    par.registerInfix(token.ASTERISK, par.parseInfixExpression)
    par.registerInfix(token.EQ,       par.parseInfixExpression)
    par.registerInfix(token.NOT_EQ,   par.parseInfixExpression)
    par.registerInfix(token.LT,       par.parseInfixExpression)
    par.registerInfix(token.GT,       par.parseInfixExpression)

    return par
}

func (this *Parser) Errors() []string {
    return this.errors
}

func (this *Parser) addTokenError(check token.TokenType, expected token.TokenType) {
    msg := fmt.Sprintf("Expected token type to be '%s' but got '%s' instead", expected, check)
    this.errors = append(this.errors, msg)
}

func (this *Parser) nextToken() {
    this.currToken = this.peekToken
    this.peekToken = this.lex.GetNextToken()
}

func (this *Parser) peekPrecedence() int {
    pre, ok := precedences[this.peekToken.Type]
    if !ok {
        return LOWEST
    }
    return pre
}

func (this *Parser) currPrecedence() int {
    pre, ok := precedences[this.currToken.Type]
    if !ok {
        return LOWEST
    }
    return pre
}

func (this *Parser) registerPrefix(tokenType token.TokenType, fn PrefixParseFn) {
    this.prefixParseFns[tokenType] = fn
}

func (this *Parser) parseIdentifierPrefix() ast.Expression {
    return ast.NewIdentifier(this.currToken, this.currToken.Literal)
}

func (this *Parser) parseIntegerLiteral() ast.Expression {
    intVal, err := strconv.ParseInt(this.currToken.Literal, 10, 64)
    if err != nil {
        msg := fmt.Sprintf("Could not parse %q as integer", this.currToken.Literal)
        this.errors = append(this.errors, msg)
        return nil
    }
    return ast.NewIntegerLiteral(this.currToken, intVal)
}

func (this *Parser) parsePrefixExpression() ast.Expression {
    exp := ast.NewPrefixExpression(this.currToken, this.currToken.Literal)
    this.nextToken()
    exp.Right = this.parseExpression(PREFIX)
    return exp
}

func (this *Parser) parseBoolean() ast.Expression {
    t := this.currToken.Type
    if t != token.TRUE && t != token.FALSE {
        msg := fmt.Sprintf("Expected current token to be boolean but got %T instead", t)
        this.errors = append(this.errors, msg)
        return nil
    }
    val := t == token.TRUE
    return ast.NewBoolean(this.currToken, val)
}

func (this *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    exp := ast.NewInfixExpression(this.currToken, left)
    precedence := this.currPrecedence()
    this.nextToken()
    exp.Right = this.parseExpression(precedence)
    return exp
}

func (this *Parser) registerInfix(tokenType token.TokenType, fn InfixParseFn) {
    this.infixParseFns[tokenType] = fn
}

func (this *Parser) parseLetStatement() *ast.LetStatement {
    stm := ast.NewLetStatement()
    hasError := false

    // Check for indentifier
    this.nextToken()
    if this.currToken.Type != token.IDENT {
        this.addTokenError(this.currToken.Type, token.IDENT)
        hasError = true
    }
    stm.Identifier = ast.NewIdentifier(this.currToken, this.currToken.Literal)

    // Check for the assign operator
    if !hasError {
        this.nextToken()
    }
    if this.currToken.Type != token.ASSIGN {
        this.addTokenError(this.currToken.Type, token.ASSIGN)
        hasError = true
    }

    // Parse the expression
    if !hasError {
        this.nextToken()
        stm.Expression = this.parseExpression(LOWEST)
    }

    for this.currToken.Type != token.SEMICOLON { this.nextToken() }

    if hasError { return nil }

    return stm
}

func (this *Parser) parseReturnStatement() *ast.ReturnStatement {
    stm := ast.NewReturnStatement()

    // TODO: Skipping the expression until find a semicolon (ParseExpression)
    for this.currToken.Type != token.SEMICOLON { this.nextToken() }

    return stm
}

func (this *Parser) parseExpression(precedence int) ast.Expression {
    prefixFn := this.prefixParseFns[this.currToken.Type]
    if utils.IsNill(prefixFn) {
        msg := fmt.Sprintf("No prefix parse function for %s found", this.currToken.Type)
        this.errors = append(this.errors, msg)
        return nil
    }
    left := prefixFn()

    for this.peekToken.Type != token.SEMICOLON {
        if precedence >= this.peekPrecedence() { break }

        infixFn := this.infixParseFns[this.peekToken.Type]
        if utils.IsNill(infixFn) {
            return left
        }
        this.nextToken()
        left = infixFn(left)
    }

    return left
}

func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stm := ast.NewExpressionStatement(this.currToken)
    stm.Expression = this.parseExpression(LOWEST)
    if this.currToken.Type != token.SEMICOLON {
        this.nextToken()
    }
    return stm
}

func (this *Parser) parseStatement() ast.Statement {
    switch this.currToken.Type {
    case token.LET:
        return this.parseLetStatement()
    case token.RETURN:
        return this.parseReturnStatement()
    default:
        return this.parseExpressionStatement()
    }
}

func (this *Parser) ParseProgram() *ast.Program {
    pro := ast.NewProgram()
    for this.currToken.Type != token.EOF {
        stm := this.parseStatement()
        if !utils.IsNill(stm) {
            pro.Statements = append(pro.Statements, stm)
        }
        this.nextToken()
    }
    return pro
}
