package controller

import (
	"deployer-operator/pkg/controller/deployer"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, deployer.Add)
}
