package event

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Default struct {
	Val any
}

type OnStart struct{}
type OnLoopCycle struct{}

var HandlerDone = errors.New("handler done")

func FindHandlers(insts []any) (handlers map[reflect.Type][]reflect.Value) {
	handlers = make(map[reflect.Type][]reflect.Value)
	for _, inst := range insts {
		t := reflect.TypeOf(inst)
		v := reflect.ValueOf(inst)
		for i := 0; i < t.NumMethod(); i++ {
			method := t.Method(i)
			if strings.HasPrefix(method.Name, "On") {
				if method.Type.NumIn() != 2 || method.Type.NumOut() != 1 || method.Type.Out(0) != reflect.TypeFor[error]() {
					fmt.Println("Invalid event handler signature:", method.Name)
					continue
				}
			} else {
				continue
			}
			eventType := method.Type.In(1)
			handlers[eventType] = append(handlers[eventType], v.Method(i))
		}
	}
	return
}

func FireFound(event any, handlersMap map[reflect.Type][]reflect.Value) (found bool, err error) {
	var handlers []reflect.Value
	if handlers, found = handlersMap[reflect.TypeOf(event)]; found {
		for _, handler := range handlers {
			errV := handler.Call([]reflect.Value{reflect.ValueOf(event)})[0]
			if !errV.IsNil() {
				err = errV.Interface().(error)
			}
		}
	}
	return
}

func Fire(event any, handlersMap map[reflect.Type][]reflect.Value) (err error) {
	found, err := FireFound(event, handlersMap)
	if err != nil {
		return
	}
	if !found {
		_, err = FireFound(Default{event}, handlersMap)
	}
	return
}
