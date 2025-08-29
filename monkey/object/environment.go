// monkey/object/environment.go

package object

type Environment struct {
    store map[string] Object
}

func NewEnvironment() *Environment {
    return &Environment {
        store: make(map[string]Object),
    }
}

func (this *Environment) Get(name string) (Object, bool) {
    var val, ok = this.store[name]
    return val, ok
}

func (this *Environment) Set(name string, val Object) {
    this.store[name] = val
}
