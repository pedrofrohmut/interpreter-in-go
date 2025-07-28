// monkey/parser/parser.go

package parser

import (
    "monkey/token"
    "monkey/lexer"
    "monkey/ast"
    "fmt"
)

const (
    // iota gives the constants a ascending numbers
    // _ skips the 0 value
    _ int = iota
    LOWEST
    EQUALS      // ==
    LESSGREATER // > or <
    SUM         // +
    PRODUCT     // *
    PREFIX      // -X or !X
    CALL        // myFunction(X)
)

type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
    lx *lexer.Lexer
    currToken token.Token
    nextToken token.Token
    errors []string
    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns map[token.TokenType]infixParseFn
}

func NewParser(lx *lexer.Lexer) *Parser {
    par := &Parser {
        lx: lx,
        errors: []string {},
    }
    par.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    par.registerPrefix(token.IDENT, par.parseIdentifier)
    // Initialize currToken and nextToken with tokens
    par.getNextToken()
    par.getNextToken()
    return par
}

func (par *Parser) GetErrors() []string {
    return par.errors
}

func (par *Parser) nextErrorMsg(tk token.TokenType) {
    err := fmt.Sprintf("Expected next token to be %s, but got %s", tk, par.nextToken.Type)
    par.errors = append(par.errors, err)
}

func (par *Parser) getNextToken() {
    par.currToken = par.nextToken
    par.nextToken = par.lx.GetNextToken()
}

func (par *Parser) currTokenIs(tk token.TokenType) bool {
    if tk != par.currToken.Type {
        return false
    }
    return true
}

func (par *Parser) nextTokenIs(tk token.TokenType) bool {
    if tk != par.nextToken.Type {
        return false
    }
    return true
}

// Random code from book. Clean this shit later
// Fix this hack random method later
func (par *Parser) expectPeek(tk token.TokenType) bool {
    if par.nextToken.Type != tk {
        par.nextErrorMsg(tk)
        return false
    } else {
        par.getNextToken()
        return true
    }
}

func (par *Parser) parseLetStatement() *ast.LetStatement {
    stm := ast.NewLetStatement(par.currToken)
    // TODO: Remove magic nextToken inside this method that should only check
    if !par.expectPeek(token.IDENT) { return nil }
    stm.Name = &ast.Identifier { Token: par.currToken, Value: par.currToken.Literal }
    // TODO: Remove magic nextToken inside this method that should only check
    if !par.expectPeek(token.ASSIGN) { return nil }
    // TODO: We're skipping the expressions until we find a semicolon
    for !par.currTokenIs(token.SEMICOLON) {
        par.getNextToken()
    }
    return stm
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
    stm := ast.NewReturnStatement(par.currToken)
    par.getNextToken()
    for !par.currTokenIs(token.SEMICOLON) {
        par.getNextToken()
    }

    return stm
}

func (par *Parser) parseExpression(precedence int) ast.Expression {
    prefix := par.prefixParseFns[par.currToken.Type]
    if prefix == nil { return nil }
    leftExp := prefix()
    return leftExp
}

func (par *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stm := ast.NewExpressionStatement(par.currToken)
    stm.Expression = par.parseExpression(LOWEST)
    if par.nextTokenIs(token.SEMICOLON) {
        par.getNextToken()
    }
    return stm
}

func (par *Parser) parseCurrStatement() ast.Statement {
    switch par.currToken.Type {
    // Statements
    case token.LET:
        return par.parseLetStatement()
    case token.RETURN:
        return par.parseReturnStatement()
    // Expressions
    default:
        return par.parseExpressionStatement()
    }
}

func (par *Parser) parseIdentifier() ast.Expression {
    return ast.NewIdentifier(par.currToken)
}

func (par *Parser) ParseProgram() *ast.Program {
    program := ast.NewProgram()

    for par.currToken.Type != token.EOF {
        stm := par.parseCurrStatement()
        if stm != nil {
            program.Statements = append(program.Statements, stm)
        }
        par.getNextToken()
    }

    return program
}

func (par *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    par.prefixParseFns[tokenType] = fn
}

func (par *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    par.infixParseFns[tokenType] = fn
}
