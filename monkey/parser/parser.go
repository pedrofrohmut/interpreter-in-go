// monkey/parser/parser.go

package parser

import (
    "monkey/token"
    "monkey/lexer"
    "monkey/ast"
)

type Parser struct {
    lx *lexer.Lexer
    currToken token.Token
    nextToken token.Token
}

func NewParser(lx *lexer.Lexer) *Parser {
    par := &Parser { lx: lx }

    // Initialize currToken and nextToken with tokens
    par.getNextToken()
    par.getNextToken()
    // t := lx.GetNextToken()
    // p.next = lx.GetNextToken()
    // p.curr = t

    return par
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

// Random code from book. Clean this shit later
// Fix this hack random method later
func (par *Parser) expectPeek(tk token.TokenType) bool {
    if par.nextToken.Type != tk {
        return false
    }
    par.getNextToken()
    return true
}

func (par *Parser) parseLetStatement() *ast.LetStatement {
    stm := ast.NewLetStatement(par.currToken)

    if !par.expectPeek(token.IDENT) {
        return nil
    }

    stm.Name = &ast.Identifier { Token: par.currToken, Value: par.currToken.Literal }

    if !par.expectPeek(token.ASSIGN) {
        return nil
    }

    // TODO: We're skipping the expressions until we find a semicolon
    for !par.currTokenIs(token.SEMICOLON) {
        par.getNextToken()
    }

    return stm
}

func (par *Parser) parseCurrStatement() ast.Statement {
    switch par.currToken.Type {
    case token.LET:
        return par.parseLetStatement()
    default:
        return nil
    }
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
