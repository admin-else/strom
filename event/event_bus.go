// Package event is the event handler
package event

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"time"
)

const (
	LoopsPerSecond = 10000 // this is just so the cpu doesn't burn itself in an unlimited infinite loop
	TimePerLoop    = time.Second / LoopsPerSecond
)

type Loop struct {
	HandlerFunctions map[reflect.Type][]reflect.Value

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
	Tick struct{}
)

// Close is fired when the loop shuts down.
type Close struct {
	Reason error
}

// HandlerDoneErr will stop the event loop and return the error if provided.
type HandlerDoneErr struct {
	Return error
}

func (e HandlerDoneErr) Error() string {
	return "HandlerDoneErr"
}

func (e HandlerDoneErr) Unwrap() error {
	return e.Return
}

// MainHandlerErr is a special error used to tell the event loop that the main loop has exited
type MainHandlerErr struct {
	Return error
}

func (e MainHandlerErr) Error() string {
	return "MainHandlerErr"
}

func (e MainHandlerErr) Unwrap() error {
	return e.Return
}

// DontForwardErr will stop the event from being forwarded to other handlers, including anything handlers.
var DontForwardErr = errors.New("dont forward")
var WhileClosingErr = errors.New("error while closing")

// CloseNonCriticalErr will close even when the event is non-critical
var CloseNonCriticalErr = errors.New("close non-critical")
var ContextDoneErr = errors.New("context done")

func ValidateHandler(h any) (eventType reflect.Type, hv reflect.Value) {
	hv = reflect.ValueOf(h)
	if hv.Kind() != reflect.Func {
		panic("expected method")
	}
	if hv.Type().NumIn() != 1 {
		panic("expected method with one argument")
	}
	if hv.Type().NumOut() != 1 || hv.Type().Out(0) != reflect.TypeFor[error]() {
		panic("expected method that returns an error")
	}
	eventType = hv.Type().In(0)
	return
}

// Register registers a handler that will not close the loop if it returns an error.
// Instead, it will log it to slog.Error
// Warning: If a bad handler is passed in the program will panic.
func (l *Loop) Register(h any) {
	l.RegisterCustomType(ValidateHandler(h))
}

// RegisterCustomType registers a handler that will not close the loop if it returns an error.
// Instead, it will log it to slog.Error
// Warning: If a bad handler is passed in the program will panic.
func (l *Loop) RegisterCustomType(eventType reflect.Type, hv reflect.Value) {
	handleFunc := func(packet any) error {
		v := hv.Call([]reflect.Value{reflect.ValueOf(packet)})[0]
		if !v.IsNil() {
			err := v.Interface().(error)
			if errors.Is(err, CloseNonCriticalErr) {
				return err
			}
			slog.Error("handler failed", "err", err)
		}
		return nil
	}
	l.RegisterDirect(eventType, reflect.ValueOf(handleFunc))
}

// RegisterCritical registers a handler that will close the loop if it returns an error.
// Warning: If a bad handler is passed in the program will panic.
func (l *Loop) RegisterCritical(h any) {
	eventType, hv := ValidateHandler(h)
	l.RegisterDirect(eventType, hv)
}

func (l *Loop) RegisterDirect(k reflect.Type, h reflect.Value) {
	l.HandlerFunctions[k] = append(l.HandlerFunctions[k], h)
}

func (l *Loop) RegisterIgnore(t any) {
	l.RegisterDirect(reflect.TypeOf(t), reflect.ValueOf(func(any) error { return nil }))
}

// FireFound triggers all handlers registered for the specified event type and returns if handlers were found or an error occurred.
func (l *Loop) FireFound(event any) (found bool, err error) {
	var handlers []reflect.Value
	if handlers, found = l.HandlerFunctions[reflect.TypeOf(event)]; found {
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
		if errors.Is(err, DontForwardErr) {
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

func (l *Loop) OnTick(_ Tick) (err error) {
	select {
	case <-l.Ctx.Done():
		err = ContextDoneErr
	default:
	}
	return
}

func (l *Loop) startLoop() (err error) {
	_ = *l // exit early on nil loop
	for {
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

func (l *Loop) StartLoop() (err error) {
	go func() {
		l.ErrChan <- MainHandlerErr{l.startLoop()}
	}()
	err = <-l.ErrChan
	l.cancel()
	errMaybe := MainHandlerErr{}
	if errors.As(err, &errMaybe) {
		err = errors.Unwrap(err)
	} else {
		for !errors.Is(<-l.ErrChan, MainHandlerErr{}) {
		}
	}
	return
}

func NewLoop() (l *Loop) {
	l = &Loop{}
	l.ErrChan = make(chan error)
	l.Ctx, l.cancel = context.WithCancel(context.Background())
	l.HandlerFunctions = make(map[reflect.Type][]reflect.Value)
	return
}
