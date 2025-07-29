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
)

type Parser struct {
    lex *lexer.Lexer
    errors []string

    // My custom variables to check tokens
    tokens []token.Token
}

func NewParser(lex *lexer.Lexer) *Parser {
    return &Parser {
        lex: lex,
        errors: []string {},
        tokens: []token.Token {},
    }
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

func (par *Parser) addTokenError(expected token.TokenType, tok token.Token) {
    err := fmt.Sprintf("Expected token type to be %s but got %s instead.", expected, tok.Type)
    par.errors = append(par.errors, err)
}

func (par *Parser) parseLetStatement() *ast.LetStatement {
    stm := ast.NewLetStatement()

    // Check for identifier
    tok := par.GetNextToken()
    if tok.Type != token.IDENT {
        par.addTokenError(token.IDENT, tok)
        return nil
    }
    stm.Identifier = tok.Literal

    // Check for assign item
    tok = par.GetNextToken()
    if tok.Type != token.ASSIGN {
        par.addTokenError(token.ASSIGN, tok)
        return nil
    }

    // TODO: Check out how to parse Expression from letStm

    // Advance tokens until find a semicolon
    for tok.Type != token.SEMICOLON {
        tok = par.GetNextToken()
    }

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

func (par *Parser) parseStatement(tok token.Token) ast.Statement {
    switch tok.Type {
    case token.LET:
        return par.parseLetStatement()
    case token.RETURN:
        return par.parseReturnStatement()
    default:
        return nil
    }
}

func (par *Parser) ParseProgram() *ast.Program {
    pro := ast.NewProgram()

    tok := par.GetNextToken()
    for tok.Type != token.EOF {
        stm := par.parseStatement(tok)
        if stm != nil {
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
