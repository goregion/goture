package goture

import (
	"context"
	"fmt"
)

// recoverCancel is a panic recovery function that converts panics to errors
// and cancels the context with the appropriate error cause.
//
// This function is designed to be used with defer in goroutines to ensure
// that any panic is properly handled and converted to an error that can
// be propagated through the context cancellation mechanism.
//
// Parameters:
//   - cancel: Context cancellation function that accepts an error cause
//
// Behavior:
//   - Recovers from panic if one occurred
//   - Converts error-type panics directly to cancellation cause
//   - Converts non-error panics to formatted error messages
//   - Does nothing if no panic occurred
//
// Usage:
//
//	defer recoverCancel(cancel)
func recoverCancel(cancel context.CancelCauseFunc) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			cancel(err)
			return
		}
		cancel(fmt.Errorf("%v", r))
	}
}

// recoverCancelForParallel is a specialized panic recovery function for parallel task execution.
// It handles panics in parallel goroutines by converting them to errors and sending them
// through a channel for centralized error collection.
//
// This function is specifically designed for use in parallel task execution scenarios
// where multiple goroutines need to report their completion status (success or failure)
// through a shared channel.
//
// Parameters:
//   - ch: Error channel for reporting task completion status
//
// Behavior:
//   - Recovers from panic if one occurred in the goroutine
//   - Converts error-type panics directly and sends to channel
//   - Converts non-error panics to formatted error messages
//   - Ensures that panicked goroutines still report their status
//   - Does nothing if no panic occurred (normal completion)
//
// Usage:
//
//	defer recoverCancelForParallel(errorChannel)
func recoverCancelForParallel(ch chan<- error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			ch <- err
			return
		}
		ch <- fmt.Errorf("%v", r)
	}
}

// makeErrorFromPanic converts a panic value into a proper error type.
// This utility function provides consistent error conversion from panic values,
// ensuring that all panic types are properly transformed into errors.
//
// The function handles the common pattern of panic recovery where the panic
// value might already be an error or might be some other type that needs
// to be converted to an error.
//
// Parameters:
//   - r: The recovered panic value (interface{})
//
// Returns:
//   - error: Properly formatted error representing the panic
//
// Behavior:
//   - Returns the error as-is if panic value is already an error type
//   - Converts non-error panic values to formatted error messages
//   - Ensures consistent error handling across panic scenarios
//
// Example:
//
//	defer func() {
//	    if r := recover(); r != nil {
//	        err := makeErrorFromPanic(r)
//	        log.Printf("Recovered from panic: %v", err)
//	    }
//	}()
func makeErrorFromPanic(r interface{}) error {
	if err, ok := r.(error); ok {
		return err
	}
	return fmt.Errorf("%v", r)
}
