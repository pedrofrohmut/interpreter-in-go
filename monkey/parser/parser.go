// monkey/parser/parser.go
/*
     A parser is a software that takes an input and builds up an AST (Abstract
   Syntac Tree) that is structural representation of the input. The parser
   often uses tokens created from a lexer.
*/

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
    SUM         // +
    PRODUCT     // *
    PREFIX      // -X or !X
    CALL        // myFunction(X)
)

// TODO: On book the all strings here are ast.Expression
type (
    PrefixParseFn func() string
    InfixParseFn func(string) string
)

type Parser struct {
    lex *lexer.Lexer
    errors []string
    prefixParseFns map[token.TokenType]PrefixParseFn
    infixParseFns map[token.TokenType]InfixParseFn

    // My custom variables to check tokens
    // TODO: check if it can be replaced to a single reference token like
    // 'currToken' instead of an array for token (lower memory footprint)
    tokens []token.Token
}

func NewParser(lex *lexer.Lexer) *Parser {
    par := &Parser {
        lex: lex,
        errors: []string {},
        tokens: []token.Token {},
    }
    par.prefixParseFns = make(map[token.TokenType]PrefixParseFn)
    par.addPrefixFn(token.IDENT, par.parseIdentifier)
    par.addPrefixFn(token.INT, par.parseIntegerLiteral)
    return par
}

func (par *Parser) parseIdentifier() string {
    return par.GetCurrToken().Literal
}

func (par *Parser) parseIntegerLiteral() string {
    curr := par.GetCurrToken()
    _, err := strconv.ParseInt(curr.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("Could not parse %q as integer", curr.Literal)
        par.errors = append(par.errors, msg)
        return ""
    }
    return curr.Literal
}

func (par *Parser) GetCurrToken() token.Token {
    if len(par.tokens) == 0 {
        return token.Token {}
    }
    return par.tokens[len(par.tokens) - 1]
}

func (par *Parser) GetNextToken() token.Token {
    tok := par.lex.GetNextToken()

    // Dont add extra EOF tokens at the end
    if tok.Type == token.EOF && par.GetCurrToken().Type == token.EOF {
        return tok
    }

    par.tokens = append(par.tokens, tok)
    return tok
}

// TODO: check if this fn is necessary
func (par *Parser) addPrefixFn(typ token.TokenType, fn PrefixParseFn) {
    par.prefixParseFns[typ] = fn
}

// TODO: check if this fn is necessary
func (par *Parser) addInfixFn(typ token.TokenType, fn InfixParseFn) {
    par.infixParseFns[typ] = fn
}

func (par *Parser) addTokenError(expected token.TokenType, tok token.Token) {
    err := fmt.Sprintf("Expected token type to be %s but got %s instead.", expected, tok.Type)
    par.errors = append(par.errors, err)
}

func (par *Parser) parseLetStatement() *ast.LetStatement {
    stm := ast.NewLetStatement()
    hasErr := false

    // Check for identifier
    tok := par.GetNextToken()
    if tok.Type != token.IDENT {
        par.addTokenError(token.IDENT, tok)
        hasErr = true
    }
    stm.Identifier = tok.Literal

    // Check for assign item
    if !hasErr {
        tok = par.GetNextToken()
    }
    if tok.Type != token.ASSIGN {
        par.addTokenError(token.ASSIGN, tok)
        hasErr = true
    }

    // TODO: Check out how to parse Expression from letStm

    // Advance tokens until find a semicolon
    for tok.Type != token.SEMICOLON {
        tok = par.GetNextToken()
    }

    if hasErr { return nil }

    return stm
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
    stm := ast.NewReturnStatement()
    var tok token.Token

    // TODO: Check out how to parse Expression from ReturnStatement

    // Advance tokens until find a semicolon
    for tok.Type != token.SEMICOLON {
        tok = par.GetNextToken()
    }

    return stm
}

func (par *Parser) parseExpression(precedence int) string {
    curr := par.GetCurrToken()
    prefix := par.prefixParseFns[curr.Type]
    if utils.IsNill(prefix) { return "" }
    return prefix()
}

func (par *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    parsedExp := par.parseExpression(LOWEST)
    if parsedExp == "" { return nil }

    stm := ast.NewExpressionStatement(par.GetCurrToken())
    stm.Expression = parsedExp

    // To skip semicolons so you can use no semicolon expressions on the REPL
    if par.GetCurrToken().Type == token.SEMICOLON {
        par.GetNextToken()
    }

    return stm
}

func (par *Parser) parseStatement() ast.Statement {
    switch par.GetCurrToken().Type {
    case token.LET:
        return par.parseLetStatement()
    case token.RETURN:
        return par.parseReturnStatement()
    default:
        return par.parseExpressionStatement()
    }
}

func (par *Parser) ParseProgram() *ast.Program {
    pro := ast.NewProgram()

    tok := par.GetNextToken()
    for tok.Type != token.EOF {
        stm := par.parseStatement()
        if !utils.IsNill(stm) {
            pro.Statements = append(pro.Statements, stm)
        }
        tok = par.GetNextToken()
    }

    return pro
}

func (par *Parser) PrintParserErrors() {
    for i, x := range par.errors {
        fmt.Printf("[%d] Parser error: %s\n", i, x)
    }
}
