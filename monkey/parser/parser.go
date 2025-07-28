// monkey/parser/parser.go
/*
      A parser is a software that takes an input and builds up an AST (Abstract
    Syntac Tree) that is structural representation of the input. The parser
    often uses tokens created from a lexer.
*/

package parser

import (
    "monkey/ast"
    "monkey/lexer"
    "monkey/token"
)

type Parser struct {
    lex *lexer.Lexer
}

func NewParser(lex *lexer.Lexer) *Parser {
    return &Parser { lex: lex }
}

func (par *Parser) parseLetStatement() *ast.LetStatement {
    stm := ast.NewLetStatement()

    // Check for identifier
    tok := par.lex.GetNextToken()
    if tok.Type != token.IDENT { return nil }
    stm.Identifier = tok.Literal

    // Check for assign item
    tok = par.lex.GetNextToken()
    if tok.Type != token.ASSIGN { return nil }

    // TODO: Check out how to parse Expression from letStm

    // Advance tokens until find a semicolon
    for tok.Type != token.SEMICOLON {
        tok = par.lex.GetNextToken()
    }

    return stm
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
    stm := ast.NewReturnStatement()
    var tok token.Token

    // TODO: Check out how to parse Expression from ReturnStatement

    // Advance tokens until find a semicolon
    for tok.Type != token.SEMICOLON {
        tok = par.lex.GetNextToken()
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

    tok := par.lex.GetNextToken()
    for tok.Type != token.EOF {
        stm := par.parseStatement(tok)
        if stm != nil {
            pro.Statements = append(pro.Statements, stm)
        }
        tok = par.lex.GetNextToken()
    }

    return pro
}
