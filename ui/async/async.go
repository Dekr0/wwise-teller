package async

import (
	"context"
	"errors"
)

// Can be fine tuned 
const InitialEvents = 32
const MaxAsyncTasks = 8

var ExceedMaxAsyncTask = errors.New("Exceeded maximum background tasks.")

type EventLoop struct {
	AsyncTaskCounter uint64 // Monotonic ID generator for async tasks
	AsyncTasks []*AsyncTask
	idle chan struct{}
}

// AsyncTask does not provide callback when an asynchronous task finishes. It 
// expects the received asynchronous task to modify shared resource with mutex, 
// and it should handle error on its own.
type AsyncTask struct {
	Id uint64
	OnProc string
	OnDone string
	Pending bool
	TaskFuncWithCancel func()
	Ctx context.Context

	// Event loop will create a closure that defers the cancellation function.
	// Do not defer cancellation upon creation
	Cancel context.CancelFunc
}

func NewEventLoop() *EventLoop {
	return &EventLoop{
		AsyncTaskCounter: 0,
		AsyncTasks: make([]*AsyncTask, MaxAsyncTasks),
		idle: make(chan struct{}, MaxAsyncTasks),
	}
}

func (e *EventLoop) taskFuncWithCancel(
	ctx context.Context,
	cancel context.CancelFunc,
	taskFunc func(context.Context),
) func() {
	return func() {
		taskFunc(ctx)
		if cancel != nil {
			cancel()
		}
		<- e.idle
	}
}

// Queue in a asynchronous task in the form of a function. DO NOT CALL CANCEL 
// FUNCTION UPON CREATION! QTask will handle it.
func (e *EventLoop) QTask(
	ctx context.Context,
	cancel context.CancelFunc,
	onProc string,
	onDone string,
	taskFunc func(context.Context),
) error {
	if ctx != nil && ctx.Err() != nil {
		panic("Asynchronous task is canceled before queuing.")
	}
	for i := range e.AsyncTasks {
		if e.AsyncTasks[i] == nil {
			e.AsyncTasks[i] = &AsyncTask{
				Id: e.AsyncTaskCounter,
				OnProc: onProc,
				OnDone: onDone,
				Pending: true,
				TaskFuncWithCancel: e.taskFuncWithCancel(ctx, cancel, taskFunc),
				Ctx: ctx,
				Cancel: cancel,
			}
			e.AsyncTaskCounter += 1
			return nil
		}
	}
	return ExceedMaxAsyncTask
}

// Check status of each asynchronous task. If it's finished, mark its occupied 
// slot as nil. Otherwise, try to schedule its execution if there are more spare 
// workers.
func (e *EventLoop) Update() []string {
	onDones := []string{}

	for i := range e.AsyncTasks {
		if e.AsyncTasks[i] == nil {
			continue
		}
		a := e.AsyncTasks[i]
		if !a.Pending {
			// Task is either canceled or finished.
			if a.Ctx != nil && a.Ctx.Err() != nil {
				onDones = append(onDones, e.AsyncTasks[i].OnDone)
				e.AsyncTasks[i] = nil
			}
			continue
		}
		select {
		case e.idle <- struct{}{}:
			a.Pending = false
			go a.TaskFuncWithCancel()
		default:
		}
	}
	
	return onDones
}
