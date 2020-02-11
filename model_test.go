package rest

import (
	"testing"
)

type Person struct {
	name string
	age int
}

func TestNewModel(t *testing.T) {
	m := makeStruct(&Person{})
	n := makeSlice(&Person{})
	var p Person
	var q []Person
	t.Logf("%T\n", &p)
	t.Logf("%T\n", m)
	t.Logf("%T\n", &q)
	t.Logf("%T\n", n)
}
