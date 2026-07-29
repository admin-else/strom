// Package event is the event handler
package event

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
)

type Loop struct {
	HandlerFunctions map[reflect.Type][]reflect.Value

	Ctx          context.Context
	Cancel       context.CancelFunc
	eventSources []reflect.SelectCase
}

type Unhandled struct {
	Val any
}

type Anything struct {
	Val any
}

// Close is fired when the loop shuts down.
type Close struct {
	Reason error
}

// HandlerDoneErr will stop the event loop and return the error if provided.
type HandlerDoneErr struct {
	Return error
}

// CallFunc is a special value that when fired with Loop.Fire it calls a function on the main loop
// Its usefull for event event sources to call functions on the main loop
// because event sources dont run on the main loop
type CallFunc struct {
	F    reflect.Value
	Args []reflect.Value
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
var ContextDoneErr = errors.New("context done")
var RecvNotOkErr = errors.New("recv not ok")

// ValidateHandler validates that h is a function taking a single argument and returning an error.
// It panics if the handler does not match the expected signature.
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

// RegisterEventSource registers an event source an event source is a channel this channel along with all other event sources will be reflect.Select 'ed until the loop closes
// here is an example of an event source function
//
//	func() (chSending <-chan any) {
//			ch := make(chan any)
//			go func() {
//				ctx := conn.Ctx
//				for {
//					select {
//					case <-ctx.Done():
//						return
//					case eventt := <-otherSource:
//						err = c.DoCommand(eventt)
//					}
//				}
//			}()
//			chSending = ch
//			return
//		}
func (l *Loop) RegisterEventSource(ch <-chan any) {
	l.eventSources = append(l.eventSources, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
}

// Register registers a handler that will not close the loop if it returns an error.
// Instead, it will log it to slog.Error
// Warning: If a bad handler is passed in the program will panic.
func (l *Loop) Register(h any) {
	eventType, hv := ValidateHandler(h)
	l.RegisterLow(eventType, hv, false)
}

// RegisterLow registers a handler that will not close the loop if it returns an error.
// Instead, it will log it to slog.Error
// This does not hold true for signaling errors like HandlerDoneErr or DontForwardErr
// Warning: If a bad handler is passed in the program will panic.
func (l *Loop) RegisterLow(eventType reflect.Type, hv reflect.Value, critical bool) {
	handleFunc := hv
	if !critical {
		handleFunc = reflect.ValueOf(func(packet any) error {
			v := hv.Call([]reflect.Value{reflect.ValueOf(packet)})[0]
			if !v.IsNil() {
				err := v.Interface().(error)
				if errors.As(err, &HandlerDoneErr{}) || errors.Is(err, DontForwardErr) {
					return err
				}
				slog.Error("handler failed", "err", err)
			}
			return nil
		})
	}
	l.HandlerFunctions[eventType] = append(l.HandlerFunctions[eventType], handleFunc)
}

// RegisterCritical registers a handler that will close the loop if it returns an error.
// Warning: If a bad handler is passed in the program will panic.
func (l *Loop) RegisterCritical(hs ...any) {
	for _, h := range hs {
		eventType, hv := ValidateHandler(h)
		l.HandlerFunctions[eventType] = append(l.HandlerFunctions[eventType], hv)
	}
}

// FireFound triggers all handlers registered for the specified event type and returns if handlers were found or an error occurred.
func (l *Loop) FireFound(event any) (found bool, err error) {
	var handlers []reflect.Value
	if handlers, found = l.HandlerFunctions[reflect.TypeOf(event)]; found {
		for _, handler := range handlers {
			errV := handler.Call([]reflect.Value{reflect.ValueOf(event)})[0]
			if !errV.IsNil() {
				err = errV.Interface().(error)
				return
			}
		}
	}
	return
}

// Fire dispatches the provided event to all appropriate handlers and invokes fallback handlers if the event is unhandled.
func (l *Loop) Fire(event any) (err error) {
	err, isErr := event.(error)
	if isErr {
		return
	}
	err = nil

	f, isCallFunc := event.(CallFunc)
	if isCallFunc {
		errV := f.F.Call(f.Args)[0].Interface()
		if errV != nil {
			err = errV.(error)
		}
		return
	}

	found, err := l.FireFound(event)
	if err != nil {
		if errors.Is(err, DontForwardErr) {
			err = nil
		}
		return
	}
	if !found {
		found, err = l.FireFound(Unhandled{event})
		if err != nil {
			return
		}
	}
	found, err = l.FireFound(Anything{event})
	return
}

func (l *Loop) ContextDoneListener() (chSending <-chan any) {
	ch := make(chan any)
	go func() {
		<-l.Ctx.Done()
		ch <- ContextDoneErr
	}()
	return ch
}

func (l *Loop) IsTypeRegistered(t reflect.Type) bool {
	_, ok := l.HandlerFunctions[t]
	return ok
}

func (l *Loop) StartLoop() (err error) {
	_ = *l // exit early on nil loop
	l.RegisterEventSource(l.ContextDoneListener())
	for {
		_, v, ok := reflect.Select(l.eventSources)
		if !ok {
			err = RecvNotOkErr
			break
		}
		err = l.Fire(v.Interface())
		if err != nil {
			break
		}
	}
	var contextDoneErr HandlerDoneErr
	if errors.As(err, &contextDoneErr) {
		err = contextDoneErr.Return
	}
	closeErr := l.Fire(Close{err})
	if closeErr != nil {
		err = errors.Join(err, errors.Join(WhileClosingErr, closeErr))
	}
	l.Cancel()
	return
}

// NewLoop creates a new Loop with a background context.
func NewLoop() (l *Loop) {
	return NewLoopCtx(context.Background())
}

// NewLoopCtx creates a new Loop instance with a cancellable context, using the provided context or context.Background if nil.
func NewLoopCtx(ctx context.Context) (l *Loop) {
	if ctx == nil {
		ctx = context.Background()
	}
	l = &Loop{}
	l.Ctx, l.Cancel = context.WithCancel(ctx)
	l.HandlerFunctions = make(map[reflect.Type][]reflect.Value)
	return
}
