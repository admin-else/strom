// Package event is the event handler
package event

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Loop struct {
	HandlerFunctions map[reflect.Type][]reflect.Value
	Handlers         []any
}

type Unhandled struct {
	Val any
}

type Anything struct {
	Val any
}

type (
	OnStart     struct{}
	OnLoopCycle struct{}
)

// ErrHandlerDone will stop the event loop.
var ErrHandlerDone = errors.New("handler done")

// ErrDontForward will stop the event from being forwarded to other handlers, including anything handlers.
var ErrDontForward = errors.New("dont forward")

func (l *Loop) FindHandlers(insts []any) {
	l.HandlerFunctions = make(map[reflect.Type][]reflect.Value)
	for _, inst := range insts {
		t := reflect.TypeOf(inst)
		v := reflect.ValueOf(inst)
		for i := 0; i < t.NumMethod(); i++ {
			method := t.Method(i)
			if strings.HasPrefix(method.Name, "On") {
				if method.Type.NumIn() != 2 || method.Type.NumOut() != 1 || method.Type.Out(0) != reflect.TypeFor[error]() {
					fmt.Println("Invalid event handler signature:", t.String(), ".", method.Name)
					continue
				}
			} else {
				continue
			}
			eventType := method.Type.In(1)
			l.HandlerFunctions[eventType] = append(l.HandlerFunctions[eventType], v.Method(i))
		}
	}
}

func (l *Loop) FireFound(event any) (found bool, err error) {
	var handlers []reflect.Value
	if handlers, found = l.HandlerFunctions[reflect.TypeOf(event)]; found {
		for _, handler := range handlers {
			errV := handler.Call([]reflect.Value{reflect.ValueOf(event)})[0]
			if !errV.IsNil() {
				err = errV.Interface().(error)
				err = errors.Join(err, fmt.Errorf("event handler failed: at "))
				return
			}
		}
	}
	return
}

func (l *Loop) Fire(event any) (err error) {
	found, err := l.FireFound(event)
	if err != nil {
		if errors.Is(err, ErrDontForward) {
			err = nil
		}
		return
	}
	if !found {
		_, err = l.FireFound(Unhandled{event})
		if err != nil {
			return
		}
	}
	_, err = l.FireFound(Anything{event})
	return
}
