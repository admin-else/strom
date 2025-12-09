package strom

import (
	"errors"
	"reflect"
)

type HandlerInst interface {
	Default(event any) (err error)
}

type OnStart struct{}
type OnLoopCycle struct{}

var HandlerDone = errors.New("handler done")

func FindHandlers(inst HandlerInst) (handlers map[reflect.Type]reflect.Value) {
	handlers = make(map[reflect.Type]reflect.Value)
	t := reflect.TypeOf(inst)
	v := reflect.ValueOf(inst)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.Type.NumIn() != 2 || method.Type.NumOut() != 1 {
			continue
		}
		if method.Type.Out(0) != reflect.TypeFor[error]() {
			continue
		}
		eventType := method.Type.In(1)
		handlers[eventType] = v.Method(i)
	}
	return
}

func FireEvent(inst HandlerInst, event any, handlers map[reflect.Type]reflect.Value) (err error) {
	if handler, ok := handlers[reflect.TypeOf(event)]; ok {
		errV := handler.Call([]reflect.Value{reflect.ValueOf(event)})[0]
		if !errV.IsNil() {
			err = errV.Interface().(error)
		}
	} else {
		err = inst.Default(event)
	}
	return
}

func (c *Conn) Start(inst HandlerInst) (err error) {
	_ = *c // exit early on nil connection
	handlers := FindHandlers(inst)
	err = FireEvent(inst, OnStart{}, handlers)
	for err == nil {
		var event any
		event, err = c.Receive()
		if err != nil {
			return
		}
		err = FireEvent(inst, event, handlers)
		if err != nil {
			break
		}
		err = FireEvent(inst, OnLoopCycle{}, handlers)
	}
	if errors.Is(err, HandlerDone) {
		err = nil
	}
	return
}
