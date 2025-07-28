// monkey/parser/parser.go
/*
      A parser is a software that takes an input and builds up an AST (Abstract
    Syntac Tree) that is structural representation of the input. The parser
    often uses tokens created from a lexer.
*/

package parser

import (
    "monkey/lexer"
    "monkey/token"
    "monkey/ast"
)

type Parser struct {
    lex *lexer.Lexer

    // TODO: check if is necessary this 2 fields
    curr token.Token
    next token.Token
}

func NewParser(lex *lexer.Lexer) *Parser {
    par := &Parser { lex: lex }
    return par
}

func (par *Parser) readFirstTokens() {
    tmp := par.lex.GetNextToken()
    par.next = par.lex.GetNextToken()
    par.curr = tmp
}

func (par *Parser) readNextToken() {
    par.curr = par.next
    par.next = par.lex.GetNextToken()
}

func (par *Parser) parseLetStatement() *ast.LetStatement {
    stm := ast.NewLetStatement()

    // Identifier
    if par.next.Type != token.IDENT { return nil }
    stm.Identifier = ast.NewIdentifier(par.next.Literal)

    par.readNextToken()

    // Checks for the '=' in the middle
    if par.next.Type != token.ASSIGN { return nil }

    par.readNextToken()

    // TODO: Check out how to parse Expression from letStm

    // Check for empty expression
    // if par.curr.Type == token.SEMICOLON { return nil }

    // Advance to semicolon
    for par.curr.Type != token.SEMICOLON {
        par.readNextToken()
    }

    return stm
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
    return nil
}

func (par *Parser) parseStatement() ast.Statement {
    switch par.curr.Type {
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

    if par.curr == (token.Token{}) {
        par.readFirstTokens()
    }

    for par.curr.Type != token.EOF {
        stm := par.parseStatement()
        if stm != nil {
            pro.Statements = append(pro.Statements, stm)
        }
        par.readNextToken()
    }

    return pro
}
