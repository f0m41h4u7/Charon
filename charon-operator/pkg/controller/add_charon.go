package controller

import (
	"charon-operator/pkg/controller/charon"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, charon.Add)
}
