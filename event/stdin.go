package event

import (
	"bufio"
	"context"
	"os"
)

type Stdin struct {
	Val string
}

type StdinHandler struct {
	*Loop
	stdInChan chan Stdin
	ctx       context.Context
	cancel    context.CancelFunc
}

func (s *StdinHandler) OnTick(_ Tick) (err error) {
	select {
	case stdin := <-s.stdInChan:
		err = s.Fire(Stdin{stdin.Val})
	default:

	}
	return
}

func (s *StdinHandler) ScanJob() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.stdInChan <- Stdin{scanner.Text()}
		}
	}
}
func (s *StdinHandler) OnStart(_ Start) (err error) {
	s.stdInChan = make(chan Stdin)
	s.ctx, s.cancel = context.WithCancel(context.Background())
	go s.ScanJob()
	return
}

func (s *StdinHandler) OnClose(_ Close) (err error) {
	s.cancel()
	return
}
