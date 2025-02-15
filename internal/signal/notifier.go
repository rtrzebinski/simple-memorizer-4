package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Notifier holds the context and channels to listen to the notifications.
type Notifier struct {
	done chan struct{}
	sig  chan os.Signal
}

// NewNotifier creates a new Notifier, if no signal is supplied, it will use
// the default signals, which are: os.Interrupt and syscall.SIGTERM.
func NewNotifier(ctx context.Context, signals ...os.Signal) *Notifier {
	if signals == nil {
		// default signals
		signals = []os.Signal{
			os.Interrupt,
			syscall.SIGTERM,
		}
	}

	notifier := Notifier{
		done: make(chan struct{}),
		sig:  make(chan os.Signal),
	}

	signal.Notify(notifier.sig, signals...)

	go notifier.listenToSignal(ctx)

	return &notifier
}

// listenToSignal is a blocking statement that listens to two channels:
//
//   - s.sig: is the os.Signal that will the triggered by the signal.Notify once
//     the expected signals are executed by the OS in the service
//   - ctx.Done(): in case of close of context, the service should also shut down.
func (s *Notifier) listenToSignal(ctx context.Context) {
	for {
		select {
		case <-s.sig:
			s.done <- struct{}{}

			return
		case <-ctx.Done():
			s.done <- struct{}{}

			return
		}
	}
}

// Done returns the call of the done channel.
func (s *Notifier) Done() <-chan struct{} { return s.done }
