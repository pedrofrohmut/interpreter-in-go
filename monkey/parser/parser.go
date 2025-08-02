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
    val, err := strconv.ParseInt(this.currToken.Literal, 10, 64)
    if err != nil {
        msg := fmt.Sprintf("Could not parse %q as integer", this.currToken.Literal)
        this.errors = append(this.errors, msg)
        return nil
    }
    return ast.NewIntegerLiteral(this.currToken, val)
}

func (this *Parser) parsePrefixExpression() ast.Expression {
    exp := ast.NewPrefixExpression(this.currToken, this.currToken.Literal)
    this.nextToken()
    exp.Right = this.parseExpression(PREFIX)
    return exp
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

    // TODO: Skipping the expression until find a semicolon (ParseExpression)
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

    fmt.Printf("Parse Expression curr token: %s\n", this.currToken.Type)

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

//
// // TODO: On book the all strings here are ast.Expression
// type (
//     PrefixParseFn func() ast.Expression
//     InfixParseFn func(ast.Expression) ast.Expression
// )
//
// type Parser struct {
//     lex *lexer.Lexer
//     errors []string
//     prefixParseFns map[token.TokenType]PrefixParseFn
//     infixParseFns map[token.TokenType]InfixParseFn
//
//     // My custom variables to check tokens
//     // TODO: check if it can be replaced to a single reference token like
//     // 'currToken' instead of an array for token (lower memory footprint)
//     tokens []token.Token
// }
//
// func NewParser(lex *lexer.Lexer) *Parser {
//     par := &Parser {
//         lex: lex,
//         errors: []string {},
//         tokens: []token.Token {},
//     }
//     par.prefixParseFns = make(map[token.TokenType]PrefixParseFn)
//     par.addPrefixFn(token.IDENT, par.parseIdentifier)
//     par.addPrefixFn(token.INT, par.parseIntegerLiteral)
//     par.addPrefixFn(token.BANG, par.parsePrefixExpression)
//     par.addPrefixFn(token.MINUS, par.parsePrefixExpression)
//     return par
// }
//
// func (par *Parser) parseIdentifier() ast.Expression {
//     return par.GetCurrToken().Literal
// }
//

// func (par *Parser) parseIntegerLiteral() ast.Expression {
//     curr := par.GetCurrToken()
//     _, err := strconv.ParseInt(curr.Literal, 0, 64)
//     if err != nil {
//         msg := fmt.Sprintf("Could not parse %q as integer", curr.Literal)
//         par.errors = append(par.errors, msg)
//         return nil
//     }
//     // return curr.Literal
// }
//
// // TODO: This is wrong somehow
// func (par *Parser) parsePrefixExpression() ast.Expression {
//     curr := par.GetCurrToken()
//     exp := &ast.PrefixExpression { Token: curr, Operator: curr.Literal }
//     exp.Right = par.parseExpression(PREFIX)
//     return exp
// }
//
// func (par *Parser) GetCurrToken() token.Token {
//     if len(par.tokens) == 0 {
//         return token.Token {}
//     }
//     return par.tokens[len(par.tokens) - 1]
// }
//
// func (par *Parser) GetNextToken() token.Token {
//     tok := par.lex.GetNextToken()
//
//     // Dont add extra EOF tokens at the end
//     if tok.Type == token.EOF && par.GetCurrToken().Type == token.EOF {
//         return tok
//     }
//
//     par.tokens = append(par.tokens, tok)
//     return tok
// }
//
// // TODO: check if this fn is necessary
// func (par *Parser) addPrefixFn(typ token.TokenType, fn PrefixParseFn) {
//     par.prefixParseFns[typ] = fn
// }
//
// // TODO: check if this fn is necessary
// func (par *Parser) addInfixFn(typ token.TokenType, fn InfixParseFn) {
//     par.infixParseFns[typ] = fn
// }
//
// func (par *Parser) addTokenError(expected token.TokenType, tok token.Token) {
//     err := fmt.Sprintf("Expected token type to be %s but got %s instead.", expected, tok.Type)
//     par.errors = append(par.errors, err)
// }
//
// func (par *Parser) parseLetStatement() *ast.LetStatement {
//     stm := ast.NewLetStatement()
//     hasErr := false
//
//     // Check for identifier
//     tok := par.GetNextToken()
//     if tok.Type != token.IDENT {
//         par.addTokenError(token.IDENT, tok)
//         hasErr = true
//     }
//     stm.Identifier = tok.Literal
//
//     // Check for assign item
//     if !hasErr {
//         tok = par.GetNextToken()
//     }
//     if tok.Type != token.ASSIGN {
//         par.addTokenError(token.ASSIGN, tok)
//         hasErr = true
//     }
//
//     // TODO: Check out how to parse Expression from letStm
//
//     // Advance tokens until find a semicolon
//     for tok.Type != token.SEMICOLON {
//         tok = par.GetNextToken()
//     }
//
//     if hasErr { return nil }
//
//     return stm
// }
//
// func (par *Parser) parseReturnStatement() *ast.ReturnStatement {
//     stm := ast.NewReturnStatement()
//     var tok token.Token
//
//     // TODO: Check out how to parse Expression from ReturnStatement
//
//     // Advance tokens until find a semicolon
//     for tok.Type != token.SEMICOLON {
//         tok = par.GetNextToken()
//     }
//
//     return stm
// }
//
// func (par *Parser) parseExpression(precedence int) string {
//     curr := par.GetCurrToken()
//     prefix := par.prefixParseFns[curr.Type]
//
//     fmt.Printf("Parse Expression Curr: %s, \tPrefix? %t\n", curr.Type, utils.IsNill(prefix))
//
//     if utils.IsNill(prefix) {
//         msg := fmt.Sprintf("No prefix parse function for %s found", curr.Type)
//         par.errors = append(par.errors, msg)
//         return ""
//     }
//
//     return prefix()
// }
//
// func (par *Parser) parseExpressionStatement() *ast.ExpressionStatement {
//     parsedExp := par.parseExpression(LOWEST)
//     if parsedExp == "" { return nil }
//
//     stm := ast.NewExpressionStatement(par.GetCurrToken())
//     stm.Expression = parsedExp
//
//     // To skip semicolons so you can use no semicolon expressions on the REPL
//     if par.GetCurrToken().Type == token.SEMICOLON {
//         par.GetNextToken()
//     }
//
//     return stm
// }
//
// func (par *Parser) parseStatement() ast.Statement {
//     switch par.GetCurrToken().Type {
//     case token.LET:
//         return par.parseLetStatement()
//     case token.RETURN:
//         return par.parseReturnStatement()
//     default:
//         return par.parseExpressionStatement()
//     }
// }
//
// func (par *Parser) ParseProgram() *ast.Program {
//     pro := ast.NewProgram()
//
//     tok := par.GetNextToken()
//     for tok.Type != token.EOF {
//         stm := par.parseStatement()
//         if !utils.IsNill(stm) {
//             pro.Statements = append(pro.Statements, stm)
//         }
//         tok = par.GetNextToken()
//     }
//
//     return pro
// }
//
// func (par *Parser) PrintParserErrors() {
//     for i, x := range par.errors {
//         fmt.Printf("[%d] Parser error: %s\n", i, x)
//     }
// }
