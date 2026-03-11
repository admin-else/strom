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
	Ctx       context.Context
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
		case <-s.Ctx.Done():
			return
		default:
			s.stdInChan <- Stdin{scanner.Text()}
		}
	}
}
func (s *StdinHandler) OnStart() (err error) {
	s.stdInChan = make(chan Stdin)
	go s.ScanJob()
	return
}
