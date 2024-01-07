package reflectx

import "reflect"

// Deref 返回对象类型。 如 typ 是指针则返回指针指向的类型
func Deref(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

// Indirect 返回最底层的value数据结构
func Indirect(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv
}
