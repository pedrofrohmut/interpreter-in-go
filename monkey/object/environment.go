// monkey/object/environment.go

package object

import (
    "bytes"
    "strings"
)

type Environment struct {
    store map[string] Object
    outer *Environment
}

func NewEnvironment() *Environment {
    return &Environment {
        store: make(map[string]Object),
    }
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
    var env = NewEnvironment()
    env.outer = outer
    return env
}

func (this *Environment) Get(name string) (Object, bool) {
    var val, ok = this.store[name]
    if !ok && this.outer != nil {
        var outerVal, okOuter = this.outer.store[name]
        return outerVal, okOuter
    }
    return val, ok
}

func (this *Environment) Set(name string, val Object) {
    this.store[name] = val
}

func (this *Environment) String() string {
    var out bytes.Buffer

    var items = []string {}
    for key, value := range this.store {
        var item string = key + ": " + value.Inspect()
        items = append(items, item)
    }

    out.WriteString("ENV { ")
    out.WriteString(strings.Join(items, ", "))
    out.WriteString(" }")

    return out.String()
}
