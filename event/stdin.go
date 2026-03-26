package event

import (
	"bufio"
	"os"
)

type Stdin struct {
	Val string
}

func StartListingStdin(l *Loop) {
	l.RegisterEventSource(func() (chSending <-chan any) {
		ch := make(chan any)
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			ctx := l.Ctx
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					return
				default:
					ch <- Stdin{scanner.Text()}
				}
			}
			if err := scanner.Err(); err != nil {
				ch <- err
			}
		}()
		chSending = ch
		return
	}())
	return
}
