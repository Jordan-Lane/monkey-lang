package object

type Environment struct {
	Outer *Environment
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvrionment(outerEnv *Environment) *Environment {
	env := NewEnvironment()
	env.Outer = outerEnv
	return env
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]
	if !ok && env.Outer != nil {
		obj, ok = env.Outer.Get(name) //Checks for the var, if its not in the current env we check if its in the Outer environment (recursively)
	}
	return obj, ok
}

func (env *Environment) Set(name string, obj Object) Object {
	env.store[name] = obj
	return obj
}
