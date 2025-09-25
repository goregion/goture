// Package goture provides a Future pattern implementation for Go.
// It allows executing tasks asynchronously and waiting for their completion.
package goture

import (
	"context"
)

// SuccessResult represents successful completion of a task
type SuccessResult struct{}

// Error implements the error interface for SuccessResult.
// It always returns an empty string since SuccessResult indicates success.
//
// Returns:
//   - Empty string indicating no error occurred
func (s SuccessResult) Error() string {
	return ""
}

// Task represents a function that can be executed asynchronously
type Task func(ctx context.Context) error

// Goture represents a future that will complete when the associated task finishes
type Goture struct {
	ctx context.Context
}

// Wait blocks until the associated task completes and returns any error that occurred.
// This method provides synchronous waiting for asynchronous task completion.
//
// Behavior:
//   - Blocks the calling goroutine until the task finishes
//   - Returns nil if the task completed successfully
//   - Returns the actual error if the task failed
//   - Handles panic recovery from the executed task
//
// Returns:
//   - error: nil on success, actual error on failure
//
// Example:
//
//	future := NewGoture(ctx, someTask)
//	if err := future.Wait(); err != nil {
//	    log.Printf("Task failed: %v", err)
//	}
func (f Goture) Wait() error {
	<-f.ctx.Done()
	cause := context.Cause(f.ctx)
	if _, ok := cause.(SuccessResult); ok {
		return nil
	}
	return cause
}

// NewGoture creates a new Goture that executes the given task asynchronously.
// The task begins execution immediately in a separate goroutine upon creation.
//
// This function provides the foundation for asynchronous task execution in Go,
// allowing you to start a task and continue with other work while waiting
// for completion when needed.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - fn: Task function to execute asynchronously
//
// Returns:
//   - Goture: A future object that can be used to wait for task completion
//
// Behavior:
//   - Starts task execution immediately in a new goroutine
//   - Handles panic recovery within the task
//   - Propagates cancellation from the parent context
//   - Returns success result when task completes without error
//
// Example:
//
//	future := NewGoture(ctx, func(ctx context.Context) error {
//	    // Your async work here
//	    return doSomeWork(ctx)
//	})
//	// Continue other work...
//	err := future.Wait() // Wait for completion
func NewGoture(ctx context.Context, fn Task) Goture {
	var localCtx, cancel = context.WithCancelCause(ctx)
	go func() {
		defer recoverCancel(cancel)
		if err := fn(localCtx); err != nil {
			cancel(err)
		} else {
			cancel(SuccessResult{})
		}
	}()
	return Goture{ctx: localCtx}
}

// NewParallelGoture creates a new Goture that executes all given tasks concurrently.
// This function enables efficient parallel processing by running multiple tasks
// simultaneously and waiting for all of them to complete.
//
// The function implements a fail-fast strategy: if any task fails, the error
// is captured and returned. If multiple tasks fail, the first encountered error
// is returned.
//
// Parameters:
//   - parentCtx: Parent context for cancellation and timeout control
//   - tasks: Variable number of Task functions to execute in parallel
//
// Returns:
//   - Goture: A future object representing the completion of all parallel tasks
//
// Behavior:
//   - Returns immediately completed future for empty task list
//   - Executes all tasks concurrently in separate goroutines
//   - Waits for ALL tasks to complete before signaling completion
//   - Handles panic recovery for each individual task
//   - Propagates the first error encountered if any task fails
//   - Returns success only if ALL tasks complete successfully
//
// Example:
//
//	future := NewParallelGoture(ctx,
//	    func(ctx context.Context) error { return task1(ctx) },
//	    func(ctx context.Context) error { return task2(ctx) },
//	    func(ctx context.Context) error { return task3(ctx) },
//	)
//	if err := future.Wait(); err != nil {
//	    log.Printf("One or more parallel tasks failed: %v", err)
//	}
func NewParallelGoture(parentCtx context.Context, tasks ...Task) Goture {
	if len(tasks) == 0 {
		// Return completed future for empty task list
		localCtx, cancel := context.WithCancelCause(parentCtx)
		cancel(SuccessResult{})
		return Goture{ctx: localCtx}
	}

	var localCtx, cancel = context.WithCancelCause(parentCtx)

	// Use sync mechanism to wait for all tasks
	completed := make(chan error, len(tasks))

	for _, fn := range tasks {
		go func(task Task) {
			defer recoverCancelForParallel(completed)
			completed <- task(localCtx)
		}(fn)
	}

	// Goroutine to wait for all tasks completion
	go func() {
		var firstError error
		for i := 0; i < len(tasks); i++ {
			if err := <-completed; err != nil && firstError == nil {
				firstError = err
			}
		}
		if firstError != nil {
			cancel(firstError)
		} else {
			cancel(SuccessResult{})
		}
	}()

	return Goture{ctx: localCtx}
}
