package object

type Enviroment struct {
    store map[string]Object
    outer *Enviroment
}

func NewEnviroment() *Enviroment {
    return &Enviroment{store: make(map[string]Object), outer: nil}
}

func NewEnclosedEnviorment(outer *Enviroment) *Enviroment {
    env := NewEnviroment()
    env.outer = outer

    return env
}

func (e *Enviroment) Get(name string) (Object, bool) {
    obj, ok := e.store[name]

    if !ok && e.outer != nil {
        obj, ok = e.outer.Get(name)
    }

    return obj, ok
}

func (e *Enviroment) Set(name string, obj Object) {
    e.store[name] = obj
}
