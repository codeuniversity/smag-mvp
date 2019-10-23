package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/codeuniversity/smag-mvp/service"
)

// Worker abstracts away all the executor lifecycle hooks, exposing a much more high level api
type Worker struct {
	executor      *service.Executor
	name          string
	step          func() error
	shutdownHooks []shutdownHook

	stopTimeout time.Duration

	shutdownOnce sync.Once
}

// Start tells the worker to go and do work in another goroutine.
// Returns a wait func that blocks until the worker is closed or encountered an error it can't handle
func (w *Worker) Start() (wait func()) {
	go w.work()

	return w.executor.WaitUntilClosed
}

// Close Worker work
func (w *Worker) Close() {
	w.shutdownOnce.Do(w.shutdown)
}

func (w *Worker) shutdown() {
	log.Println("stopping", w.name)
	w.executor.Stop()
	log.Println("waiting for work to stop")
	w.executor.WaitUntilStopped(w.stopTimeout)

	log.Println("calling shutdown hooks")
	for _, hook := range w.shutdownHooks {
		log.Println("shutting down", hook.name)
		err := w.callHookWithRecover(hook.f)
		if err != nil {
			log.Println("encountered error on shutdown: ", err)
		}
	}

	log.Println(w.name, "is shut down")
	w.executor.MarkAsClosed()
}

func (w *Worker) work() {
	defer w.Close()

	defer func() {
		log.Println(w.name, "is done here")
		w.executor.MarkAsStopped()
	}()

	log.Println("starting", w.name)

	for w.executor.IsRunning() {
		err := w.callHookWithRecover(w.step)
		if err != nil {
			log.Println("encountered error while working: ", err)
			return
		}
	}
}

func (w *Worker) callHookWithRecover(hook func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %s", r)
		}
	}()

	err = hook()
	return
}
