package worker

import (
	"errors"
	"time"

	"github.com/codeuniversity/smag-mvp/service"
)

const (
	defaultStopTimeout = 5 * time.Second
)

// Builder is for configuring a worker
type Builder struct {
	name          string
	step          func() error
	shutdownHooks []shutdownHook
	stopTimeout   time.Duration
}

// WithName sets the worker name (required)
func (b Builder) WithName(n string) Builder {
	b.name = n
	return b
}

// WithWorkStep sets the step func a worker should repeatedly call (required)
func (b Builder) WithWorkStep(s func() error) Builder {
	b.step = s
	return b
}

// AddShutdownHook registers the hook to be called on shutdown of the worker
func (b Builder) AddShutdownHook(hookName string, hook func() error) Builder {
	b.shutdownHooks = append(b.shutdownHooks, shutdownHook{
		f:    hook,
		name: hookName,
	})

	return b
}

// WithStopTimeout changes the timeout used when stopping the worker loop. If not set, uses defaultStopTimeout
func (b Builder) WithStopTimeout(t time.Duration) Builder {
	b.stopTimeout = t
	return b
}

// Build a Worker with the given configuration
func (b Builder) Build() (*Worker, error) {
	if !b.valid() {
		return nil, errors.New("could not build worker: both name and work step have to be set")
	}
	w := &Worker{
		executor:      service.New(),
		name:          b.name,
		step:          b.step,
		shutdownHooks: b.shutdownHooks,
		stopTimeout:   b.stopTimeout,
	}

	if w.stopTimeout == 0 {
		w.stopTimeout = defaultStopTimeout
	}

	return w, nil
}

// MustBuild a Worker with the given configuration. Panics if not all required config is given
func (b Builder) MustBuild() *Worker {
	w, err := b.Build()

	if err != nil {
		panic(err)
	}

	return w
}

func (b Builder) valid() bool {
	return b.name != "" && b.step != nil
}

type shutdownHook struct {
	f    func() error
	name string
}
