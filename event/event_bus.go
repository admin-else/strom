// Package event is the event handler
package event

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"
)

const (
	LoopsPerSecond = 100 // i dunno seems reasonable
	TimePerLoop    = time.Second / LoopsPerSecond
)

type Loop struct {
	handlerFunctions map[reflect.Type][]reflect.Value
	Handlers         []any

	Ctx     context.Context
	ErrChan chan error
	cancel  context.CancelFunc
}

type Unhandled struct {
	Val any
}

type Anything struct {
	Val any
}

type (
	Start struct{}
	Tick  struct{}
)

// HandlerDoneErr will stop the event loop and return the error if provided.
type HandlerDoneErr struct {
	Return error
}

// Close is fired when the loop shuts down.
type Close struct {
	Reason error
}

func (e HandlerDoneErr) Error() string {
	return "HandlerDoneErr"
}

func (e HandlerDoneErr) Unwrap() error {
	return e.Return
}

// ErrDontForward will stop the event from being forwarded to other handlers, including anything handlers.
var ErrDontForward = errors.New("dont forward")
var WhileClosingErr = errors.New("error while closing")

// UpdateHandlerFunctions initializes and populates the handlerFunctions map by reflecting over methods of registered handlers.
func (l *Loop) UpdateHandlerFunctions() {
	l.handlerFunctions = make(map[reflect.Type][]reflect.Value)
	for _, inst := range l.Handlers {
		t := reflect.TypeOf(inst)
		v := reflect.ValueOf(inst)
		for i := 0; i < t.NumMethod(); i++ {
			method := t.Method(i)
			if strings.HasPrefix(method.Name, "On") {
				if method.Type.NumIn() != 2 || method.Type.NumOut() != 1 || method.Type.Out(0) != reflect.TypeFor[error]() {
					slog.Warn("Invalid event handler signature", "at", t.String()+"."+method.Name)
					continue
				}
			} else {
				continue
			}
			eventType := method.Type.In(1)
			l.handlerFunctions[eventType] = append(l.handlerFunctions[eventType], v.Method(i))
			slog.Debug("Registered Handler", "at", t.String()+"."+method.Name, "for", eventType.String())
		}
	}
}

// FireFound triggers all handlers registered for the specified event type and returns if handlers were found or an error occurred.
func (l *Loop) FireFound(event any) (found bool, err error) {
	var handlers []reflect.Value
	if handlers, found = l.handlerFunctions[reflect.TypeOf(event)]; found {
		for _, handler := range handlers {
			errV := handler.Call([]reflect.Value{reflect.ValueOf(event)})[0]
			if !errV.IsNil() {
				err = errV.Interface().(error)
				err = errors.Join(fmt.Errorf("event handler failed at: %s(%T)", handler.Type().In(0), event), err)
				return
			}
		}
	}
	return
}

// Fire dispatches the provided event to all appropriate handlers and invokes fallback handlers if the event is unhandled.
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

func (l *Loop) startLoop() (err error) {
	_ = *l // exit early on nil loop
	l.UpdateHandlerFunctions()

	err = l.Fire(Start{})
	for err == nil {
		tickStartTime := time.Now()
		err = l.Fire(Tick{})
		if err != nil {
			break
		}
		time.Sleep(TimePerLoop - time.Since(tickStartTime))
	}
	if errors.Is(err, HandlerDoneErr{}) {
		err = errors.Unwrap(err)
	}
	closeErr := l.Fire(Close{err})
	if closeErr != nil {
		err = errors.Join(err, errors.Join(WhileClosingErr, closeErr))
	}
	return
}

func (l *Loop) Start() (err error) {
	l.ErrChan = make(chan error)
	l.Ctx, l.cancel = context.WithCancel(context.Background())
	go func() {
		l.ErrChan <- l.startLoop()
	}()
	err = <-l.ErrChan
	l.cancel()
	return
}
