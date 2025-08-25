// Have to use this weird name because golang do some auto stuff with 'tests'
// And cannot just put on utils because of cycling imports

package tezts

import (
    "monkey/utils"
    "monkey/parser"
    "monkey/ast"
    "testing"
    "fmt"
)

func TryGetInput[T any](t *testing.T, tests []T) string {
    var input, err = utils.GetInput(tests)
    if err != nil {
        t.Fatalf("ERROR: Could not get input from tests")
    }
    return input
}

func CheckForParserErrors(t *testing.T, parser *parser.Parser) {
    if len(parser.Errors()) > 0 {
        for i, err := range parser.Errors() {
            fmt.Printf("# [%d] - ERROR: %s\n", i, err)
        }
        t.Fatalf("Has parser errors")
    }
}

func CheckProgram(t *testing.T, program *ast.Program, expLenStms int) {
    if program == nil {
        t.Fatalf("Program is nil")
    }
    if len(program.Statements) != expLenStms {
        t.Fatalf("Expected program to have length of %d but got %d instead", expLenStms, len(program.Statements))
    }
}
