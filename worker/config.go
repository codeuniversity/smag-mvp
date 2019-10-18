package worker

import (
	"errors"

	"github.com/codeuniversity/smag-mvp/service"
)

// Builder is for configuring a worker
type Builder struct {
	name          string
	step          func() error
	shutdownHooks []shutdownHook
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

// Build a Worker with the given configuration
func (b Builder) Build() (*Worker, error) {
	if !b.valid() {
		return nil, errors.New("could not build worker: both name and work step have to be set")
	}
	return &Worker{
		executor:      service.New(),
		name:          b.name,
		step:          b.step,
		shutdownHooks: b.shutdownHooks,
	}, nil
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
