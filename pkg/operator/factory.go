package operator

import (
	"reflect"

	"github.com/go-jarvis/hertzer/pkg/reflectx"
)

// OperatorFactory to create new operator
type OperatorFactory struct {
	Type     reflect.Type
	Operator Operator
}

func NewOperatorFactory(op Operator) *OperatorFactory {
	fact := &OperatorFactory{}

	// get real operator refelct type
	fact.Type = reflectx.Deref(reflect.TypeOf(op))
	fact.Operator = op

	return fact
}

// New create a new operator
func (fact *OperatorFactory) New() Operator {

	opc := reflect.New(fact.Type).Interface().(Operator)
	return opc
}

type DeepCopier interface {
	DeepCopy() Operator
}

func DeepCopy(op Operator) Operator {
	if copier, ok := op.(DeepCopier); ok {
		return copier.DeepCopy()
	}
	return NewOperatorFactory(op).New()
}
