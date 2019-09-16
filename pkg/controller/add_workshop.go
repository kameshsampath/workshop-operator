package controller

import "github.com/kameshsampath/workshop-operator/pkg/controller/workshop"

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, workshop.Add)
}
