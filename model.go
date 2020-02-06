package rest

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

type TweakFunc func(r *Rest, c *gin.Context)
type TransactionFunc func(i interface{})

type Model struct {
	name              string
	instance          interface{}
	GetModelFunc      TweakFunc
	GetModelIDFunc    TweakFunc
	PostModelFunc     TweakFunc
	DeleteModelIDFunc TweakFunc
	PutModelIDFunc    TweakFunc
	// Initial a pool to reduce the frequency of making instances by reflection and gc
	InstancePool      chan interface{}
	InstanceSlicePool chan interface{}
}

func NewModel(instance interface{}) *Model {
	t := reflect.TypeOf(instance)
	m :=  &Model{
		name:              strings.ToLower(t.Name()),
		instance:          instance,
	}
	m.SetPoolSize(20)
	return m
}

func (m *Model) SetPoolSize(size int) {
	instancePool := make(chan interface{}, size)
	instanceSlicePool := make(chan interface{}, size)
	for i := 0; i < size; i++ {
		instancePool<- makeStruct(m.instance)
		instanceSlicePool<- makeSlice(m.instance)
	}
	m.InstancePool = instancePool
	m.InstanceSlicePool = instanceSlicePool
}

func (m *Model) OperateInstance(f TransactionFunc) {
	select {
	case i := <-m.InstancePool:
		f(i)
		m.InstancePool<- i
	default:
		f(makeStruct(m.instance))
	}
}

func (m *Model) OperateInstanceSlice(f TransactionFunc) {
	select {
	case i := <-m.InstanceSlicePool:
		f(i)
		m.InstanceSlicePool<- i
	default:
		f(makeSlice(m.instance))
	}
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