// monkey/ast/ast.go
/*
     AST - Abstract Syntax Tree. Is the data structure the results for parsing
   the source code of monkey language
*/

package ast

import (
    "bytes"
    "strconv"
    "monkey/token"
)

type Node interface {
    TokenLiteral() string
    String() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

// @Impl
func (this *Program) TokenLiteral() string {
    if len(this.Statements) == 0 { return "" }
    return this.Statements[0].TokenLiteral()
}

// @Impl
func (this *Program) String() string {
    var out bytes.Buffer
    for _, stm := range this.Statements {
        out.WriteString(stm.String())
    }
    return out.String()
}

func NewProgram() *Program {
    return &Program {
        Statements: [] Statement{},
    }
}

// // @Impl
// func (pro *Program) String() string {
//     var out bytes.Buffer
//     for _, stm := range pro.Statements {
//         out.WriteString(stm.String())
//     }
//     return out.String()
// }
//
// func NewProgram() *Program {
//     return &Program { Statements: [] Statement{} }
// }

type Identifier struct {
    Token token.Token
    Value string
}

// @Impl
func (this *Identifier) expressionNode() {}

// @Impl
func (this *Identifier) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *Identifier) String() string { return this.Value }

func NewIdentifier(tok token.Token, val string) *Identifier {
    return &Identifier { Token: tok, Value: val }
}

// // TODO: Keeping checking if Identifier needs to be a separated struct for reuse
// type LetStatement struct {
//     Statement
//     Token token.Token
//     Identifier string
//     Expression string
// }
type LetStatement struct {
    Token token.Token
    Identifier *Identifier
    Expression Expression
}

// @Impl
func (this *LetStatement) statementNode() {}

// @Impl
func (this *LetStatement) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *LetStatement) String() string {
    var out bytes.Buffer
    out.WriteString(this.Token.Literal + " ")
    out.WriteString(this.Identifier.Value + " = ")
    if this.Expression != nil {
        out.WriteString(this.Expression.String())
    }
    out.WriteString(";")
    return out.String()
}

func NewLetStatement() *LetStatement {
    return &LetStatement {
        Token: token.NewTokenStr(token.LET, "let"),
    }
}

//
// // @Impl
// func (stm *LetStatement) String() string {
//     var out bytes.Buffer
//     out.WriteString(stm.Token.Literal + " ")
//     out.WriteString(stm.Identifier + " = ")
//     if stm.Expression != "" {
//         out.WriteString(stm.Expression)
//     }
//     out.WriteString(";")
//     return out.String()
// }
//
// func NewLetStatement() *LetStatement {
//     return &LetStatement {
//         Token: token.NewTokenStr(token.LET, "let"),
//     }
// }

type ReturnStatement struct {
    Token token.Token
    Expression Expression
}

// @Impl
func (this *ReturnStatement) statementNode() {}

// @Impl
func (this *ReturnStatement) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *ReturnStatement) String() string {
    var out bytes.Buffer
    out.WriteString(this.Token.Literal + " ")
    if this.Expression != nil {
        out.WriteString(this.Expression.String())
    }
    out.WriteString(";")
    return out.String()
}

func NewReturnStatement() *ReturnStatement {
    return &ReturnStatement {
        Token: token.NewTokenStr(token.RETURN, "return"),
    }
}

// type ReturnStatement struct {
//     Statement
//     Token token.Token
//     Expression string
// }
//
// // @Impl
// func (stm *ReturnStatement) String() string {
//     var out bytes.Buffer
//     out.WriteString(stm.Token.Literal + " ")
//     if stm.Expression != "" {
//         out.WriteString(stm.Expression)
//     }
//     out.WriteString(";")
//     return out.String()
// }
//
// func NewReturnStatement() *ReturnStatement {
//     return &ReturnStatement {
//         Token: token.NewTokenStr(token.RETURN, "return"),
//     }
// }
//

// We are having Expression Statements because in monkey you can have expressions
// as statements. Exp: 5 * 5 + 3;. So it is needed to have it as statement here
type ExpressionStatement struct {
    Token token.Token // The first token of the expression
    Expression Expression
}

// @Impl
func (this *ExpressionStatement) statementNode() {}

// @Impl
func (this *ExpressionStatement) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *ExpressionStatement) String() string {
    if this.Expression == nil {
        return ""
    }
    return this.Expression.String()
}

func NewExpressionStatement(tok token.Token) *ExpressionStatement {
    return &ExpressionStatement { Token: tok }
}

// // We are having Expression Statements because in monkey you can have expressions
// // as statements. Exp: 5 * 5 + 3;. So it is needed to have it as statement here
// type ExpressionStatement struct {
//     Token token.Token // The first token of the expression
//     Expression Expression
// }
//
// func (this *ExpressionStatement) statementNode() {}
//
// // @Impl
// func (stm *ExpressionStatement) String() string {
//     var out bytes.Buffer
//     if stm.Expression != "" {
//         out.WriteString(stm.Expression)
//     }
//     // TODO: Check if should just return an empty string or ; is fine
//     out.WriteString(";")
//     return out.String()
// }
//
// func NewExpressionStatement(first token.Token) *ExpressionStatement {
//     return &ExpressionStatement { Token: first }
// }
//

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

// @Impl
func (this *IntegerLiteral) expressionNode() {}

// @Impl
func (this *IntegerLiteral) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *IntegerLiteral) String() string { return this.Token.Literal }

func NewIntegerLiteral(val int64) *IntegerLiteral {
    return &IntegerLiteral {
        Token: token.Token { Type: token.INT, Literal: strconv.FormatInt(val, 10) },
        Value: val,
    }
}

// type IntegerLiteral struct {
//     Token token.Token
//     Value int64
// }
//
// func NewIntegerLiteral(val int64) *IntegerLiteral {
//     return &IntegerLiteral {
//         Token: token.Token { Type: token.INT, Literal: strconv.FormatInt(val, 10) },
//         Value: val,
//     }
// }
//
// // @Impl
// func (lit *IntegerLiteral) String() string {
//     return lit.Token.Literal
// }
//
// type PrefixExpression struct {
//     Token token.Token // The prefix token
//     Operator string
//     Right string
// }
//
// func (this *PrefixExpression) expressionNode() {}
//
// func (exp *PrefixExpression) String() string {
//     var out bytes.Buffer
//     out.WriteString("(")
//     out.WriteString(exp.Operator)
//     out.WriteString(exp.Right)
//     out.WriteString(")")
//     return out.String()
// }
//
// func NewPrefixExpression() *PrefixExpression {
//     return &PrefixExpression {}
// }
