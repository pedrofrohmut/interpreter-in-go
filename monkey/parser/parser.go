// monkey/parser/parser.go

package parser

import (
    "monkey/token"
    "monkey/lexer"
    "monkey/ast"
    "fmt"
)

type Parser struct {
    lx *lexer.Lexer
    currToken token.Token
    nextToken token.Token
    errors []string
}

func NewParser(lx *lexer.Lexer) *Parser {
    par := &Parser {
        lx: lx,
        errors: []string {},
    }

    // Initialize currToken and nextToken with tokens
    par.getNextToken()
    par.getNextToken()
    // tmp := lx.GetNextToken()
    // par.nextToken = lx.GetNextToken()
    // par.currToken = tmp

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
    if !par.expectPeek(token.IDENT) {
        return nil
    }

    stm.Name = &ast.Identifier { Token: par.currToken, Value: par.currToken.Literal }

    // TODO: Remove magic nextToken inside this method that should only check
    if !par.expectPeek(token.ASSIGN) {
        return nil
    }

    // TODO: We're skipping the expressions until we find a semicolon
    for !par.currTokenIs(token.SEMICOLON) {
        par.getNextToken()
    }

    return stm
}

func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
    stm := ast.NewReturnStatement(par.currToken)

    // if !par.currTokenIs(token.RETURN) {
    //     fmt.Printf("First token is not Return\n")
    // } else {
    //     fmt.Printf("It begins with return\n")
    // }

    par.getNextToken()
    for !par.currTokenIs(token.SEMICOLON) {
        par.getNextToken()
    }

    return stm
}

func (par *Parser) parseCurrStatement() ast.Statement {
    switch par.currToken.Type {
    case token.LET:
        return par.parseLetStatement()
    case token.RETURN:
        return par.parseReturnStatement()
    default:
        fmt.Printf("# WARNING: Not an expected type of statement\n")
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
