package firemodel

import "github.com/pkg/errors"

var (
	registeredModelers = map[string]Modeler{}
)

func RegisterModeler(name string, m Modeler) {
	if _, ok := registeredModelers[name]; ok {
		panic(errors.Errorf("firemodel: %s modeler already registered", name))
	}
	registeredModelers[name] = m
}

type Language struct {
	Language string
	Output   string
}

func AllModelers() (ret []string) {
	ret = []string{}
	for modelerName, _ := range registeredModelers {
		ret = append(ret, modelerName)
	}
	return ret
}

func (l Language) Modeler() Modeler {
	m, ok := registeredModelers[l.Language]
	if !ok {
		err := errors.Errorf("firemodel: config includes unimplemented language: %s (don't forget to _ import the modeler)", l)
		panic(err)
	}
	return m
}
