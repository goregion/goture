package goture

import (
	"context"
	"fmt"
)

func recoverCancel(cancel context.CancelCauseFunc) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			cancel(err)
			return
		}
		cancel(fmt.Errorf("%v", r))
	}
}

func recoverCancelForParallel(ch chan<- error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			ch <- err
			return
		}
		ch <- fmt.Errorf("%v", r)
	}
}

func makeErrorFromPanic(r interface{}) error {
	if err, ok := r.(error); ok {
		return err
	}
	return fmt.Errorf("%v", r)
}
