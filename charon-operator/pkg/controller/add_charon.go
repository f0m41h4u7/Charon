package controller

import (
	"charon-operator/pkg/controller/charon"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, charon.Add)
}
