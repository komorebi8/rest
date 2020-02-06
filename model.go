package rest

import "reflect"

type Model struct {
	name string
	instance interface{}
	GetModel TweakFunc
	GetModelID TweakFunc
	PostModel TweakFunc
	DeleteModelID TweakFunc
	PutModelID TweakFunc
}


// returns *[]instance
// Using make() to generate a slice will cause an unaddressed pointer error.
func makeSlice(instance interface{}) interface{} {
	t := reflect.TypeOf(instance)
	slice := reflect.MakeSlice(reflect.SliceOf(t), 10, 10)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)
	return x.Interface()
}

// returns *instance
func makeStruct(instance interface{}) interface{} {
	st := reflect.TypeOf(instance)
	x := reflect.New(st)
	return x.Interface()
}