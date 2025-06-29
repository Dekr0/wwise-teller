package async

import (
	"context"
	"errors"
	"slices"
)

const InitialSycnTasks = 32
const MaxAsyncTasks = 8

var ExceedMaxAsyncTask = errors.New("Exceeded maximum background tasks.")

type EventLoop struct {
	SyncTaskCounter  uint64
	// For handling some user interactions that need to be run after rendering 
	// phase
	SyncTasks        []*Task 
	AsyncTaskCounter uint64
	AsyncTasks       []*Task
	idle             chan struct{}
}

// Asynchronous task should use synchronous primitives to modify share resources
type Task struct {
	// loop *EventLoop (chaining)
	Id uint64
	OnProcMsg string
	OnDoneMsg string
	Pending bool
	TaskFuncWithCancel func()
	Ctx context.Context
	Cancel context.CancelFunc
}

func NewEventLoop() EventLoop {
	return EventLoop{
		SyncTaskCounter : 0,
		SyncTasks       : make([]*Task, 0, InitialSycnTasks),
		AsyncTaskCounter: 0,
		AsyncTasks      : make([]*Task, MaxAsyncTasks),
		idle            : make(chan struct{}, MaxAsyncTasks),
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

// Delay execution at the end of render phase
func (e *EventLoop) MustRun(taskFunc func(context.Context)) {
	e.SyncTasks = append(e.SyncTasks, &Task{
		Id: e.SyncTaskCounter,
		OnProcMsg         : "",
		OnDoneMsg         : "",
		Pending           : false,
		TaskFuncWithCancel: e.taskFuncWithCancel(nil, nil, taskFunc),
		Ctx               : nil,
		Cancel            : nil,
	})
	e.SyncTaskCounter += 1
}

// Event loop will call cancel function
func (e *EventLoop) QTask(
	ctx       context.Context,
	cancel    context.CancelFunc,
	onProcMsg string,
	onDoneMsg string,
	taskFunc  func(context.Context),
) error {
	if ctx != nil && ctx.Err() != nil {
		panic("Asynchronous task is canceled before queuing.")
	}
	for i := range e.AsyncTasks {
		if e.AsyncTasks[i] == nil {
			e.AsyncTasks[i] = &Task{
				Id                : e.AsyncTaskCounter,
				OnProcMsg         : onProcMsg,
				OnDoneMsg         : onDoneMsg,
				Pending           : true,
				TaskFuncWithCancel: e.taskFuncWithCancel(ctx, cancel, taskFunc),
				Ctx               : ctx,
				Cancel            : cancel,
			}
			e.AsyncTaskCounter += 1
			return nil
		}
	}
	return ExceedMaxAsyncTask
}

func (e *EventLoop) Update() []string {
	for _, t := range e.SyncTasks {
		t.TaskFuncWithCancel()
	}
	e.SyncTasks = slices.Delete(e.SyncTasks, 0, len(e.SyncTasks))

	onDones := []string{}
	for i := range e.AsyncTasks {
		if e.AsyncTasks[i] == nil {
			continue
		}
		a := e.AsyncTasks[i]
		if !a.Pending {
			// Task is either canceled or finished.
			if a.Ctx != nil && a.Ctx.Err() != nil {
				onDones = append(onDones, e.AsyncTasks[i].OnDoneMsg)
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

type TaskStat struct  {
	TotalNumAsyncTask   uint8
	TotalNumSyncTask    uint8
	NumRunningAsyncTask uint8
	NumPendingAsyncTask uint8
}

func (e *EventLoop) TaskStatus() TaskStat {
	stat := TaskStat{}
	for i := range e.AsyncTasks {
		if e.AsyncTasks[i] != nil {
			stat.TotalNumAsyncTask += 1
			if e.AsyncTasks[i].Pending {
				stat.NumPendingAsyncTask += 1
			} else {
				stat.NumRunningAsyncTask += 1
			}
		}
	}
	stat.TotalNumSyncTask = uint8(len(e.SyncTasks))
	return stat
}
