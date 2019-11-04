package service

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//Service is a closeable service that usually includes an Executor
type Service interface {
	Close()
}

// CloseOnSignal calls the closeFunc on the os signals SIGINT and SIGTERM
func CloseOnSignal(s Service) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		sig := <-signals
		log.Println("Received Signal:", sig)

		s.Close()
	}()
}

//Executor handles gracefully closing execution for services
//Closing an executor goes through the states running, stopping, stopped and closed.
type Executor struct {
	stopChan    chan struct{}
	stoppedChan chan struct{}
	closedChan  chan struct{}
}

//New returns an Executor ready for use.
func New() *Executor {
	return &Executor{
		stopChan:    make(chan struct{}, 1),
		stoppedChan: make(chan struct{}, 1),
		closedChan:  make(chan struct{}, 1),
	}
}

//IsRunning should be used to determine if the execution should be halted.
// Will be false until Close() is called
func (e *Executor) IsRunning() bool {
	return len(e.stopChan) == 0
}

// Stop is to be called when the executor should stop
// This is only safe to call once
func (e *Executor) Stop() {
	e.stopChan <- struct{}{}
}

//MarkAsStopped is to be called when the execution was stopped.
//All used resources can be cracefully closed and dispossed of now
func (e *Executor) MarkAsStopped() {
	e.stoppedChan <- struct{}{}
}

// WaitUntilStopped blocks until either MarkAsStopped() was called in any goroutine
// or the timeout has passed
func (e *Executor) WaitUntilStopped(timeout time.Duration) {
	t := time.NewTimer(timeout)
	select {
	case <-t.C:
		break
	case <-e.stoppedChan:
		t.Stop()
		break
	}
}

//MarkAsClosed is to be called when the service using the executor is closed.
//The process can be stopped or killed now
func (e *Executor) MarkAsClosed() {
	e.closedChan <- struct{}{}
}

// WaitUntilClosed waits until the Close func call of the executor is finished
func (e *Executor) WaitUntilClosed() {
	<-e.closedChan
}
