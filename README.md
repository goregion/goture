# Goture - Go Future Pattern Implementation

[![Go Reference](https://pkg.go.dev/badge/github.com/goregion/goture.svg)](https://pkg.go.dev/github.com/goregion/goture)
[![Go Report Card](https://goreportcard.com/badge/github.com/goregion/goture)](https://goreportcard.com/report/github.com/goregion/goture)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Goture** is a lightweight Go library that provides a Future pattern implementation for asynchronous task execution. It allows you to execute tasks concurrently and wait for their completion with built-in panic recovery and context support.

## Features

- 🚀 **Asynchronous Task Execution**: Execute tasks concurrently without blocking
- 🔄 **Parallel Processing**: Run multiple tasks simultaneously with efficient coordination
- 📦 **Generic Results**: Type-safe result handling with Go generics
- 🛡️ **Panic Recovery**: Automatic panic handling and conversion to errors
- ⏰ **Context Support**: Full context.Context integration for cancellation and timeouts
- 📊 **Result Collection**: Collect results from multiple parallel tasks
- 🎯 **Fail-Fast**: Early termination on first error in parallel execution

## Installation

```bash
go get github.com/goregion/goture
```

## Requirements

- Go 1.23+ (uses generics)

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/goregion/goture"
)

func main() {
    ctx := context.Background()

    // Simple async task
    task := func(ctx context.Context) error {
        // Your async work here
        fmt.Println("Task executed!")
        return nil
    }

    future := goture.NewGoture(ctx, task)
    
    // Do other work while task runs...
    
    // Wait for completion
    if err := future.Wait(); err != nil {
        fmt.Printf("Task failed: %v\n", err)
    }
}
```

### Parallel Task Execution

```go
func main() {
    ctx := context.Background()

    task1 := func(ctx context.Context) error {
        time.Sleep(100 * time.Millisecond)
        fmt.Println("Task 1 completed")
        return nil
    }

    task2 := func(ctx context.Context) error {
        time.Sleep(200 * time.Millisecond)
        fmt.Println("Task 2 completed")
        return nil
    }

    // Execute tasks in parallel
    future := goture.NewParallelGoture(ctx, task1, task2)
    
    if err := future.Wait(); err != nil {
        fmt.Printf("One or more tasks failed: %v\n", err)
    }
    fmt.Println("All tasks completed")
}
```

### Tasks with Results

```go
func main() {
    ctx := context.Background()

    // Task that returns a result
    task := func(ctx context.Context) (string, error) {
        time.Sleep(50 * time.Millisecond)
        return "Hello, Goture!", nil
    }

    future := goture.NewGotureWithResult(ctx, task)
    result, err := future.Wait()
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Result: %s\n", result)
}
```

### Parallel Tasks with Results

```go
func main() {
    ctx := context.Background()

    // Multiple tasks returning results
    tasks := []goture.TaskWithResult[int]{
        func(ctx context.Context) (int, error) { return 10, nil },
        func(ctx context.Context) (int, error) { return 20, nil },
        func(ctx context.Context) (int, error) { return 30, nil },
    }

    future := goture.NewParallelWithResult(ctx, tasks...)
    results, err := future.Wait()
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Results: %v\n", results) // [10 20 30]
    
    // Calculate sum
    sum := 0
    for _, result := range results {
        sum += result
    }
    fmt.Printf("Sum: %d\n", sum) // 60
}
```

## API Reference

### Core Types

#### `Goture`
Represents a future that completes when the associated task finishes.

```go
type Goture struct {
    // private fields
}

// Wait blocks until the task completes and returns any error
func (f Goture) Wait() error
```

#### `GotureWithResult[T]`
Represents a future that completes with a typed result.

```go
type GotureWithResult[ResultType any] struct {
    // private fields
}

// Wait blocks until the task completes and returns the result and any error
func (f GotureWithResult[T]) Wait() (T, error)
```

### Task Types

```go
// Task represents a function executed asynchronously
type Task func(ctx context.Context) error

// TaskWithResult represents a function that returns a typed result
type TaskWithResult[ResultType any] func(ctx context.Context) (ResultType, error)
```

### Constructor Functions

#### Single Task Execution
```go
// Execute a single task asynchronously
func NewGoture(ctx context.Context, fn Task) Goture

// Execute a single task that returns a result
func NewGotureWithResult[T any](ctx context.Context, fn TaskWithResult[T]) GotureWithResult[T]
```

#### Parallel Task Execution
```go
// Execute multiple tasks in parallel
func NewParallelGoture(parentCtx context.Context, tasks ...Task) Goture

// Execute multiple tasks in parallel and collect results
func NewParallelWithResult[T any](parentCtx context.Context, tasks ...TaskWithResult[T]) GotureWithResult[[]T]
```

## Error Handling

Goture provides robust error handling:

- **Panic Recovery**: All panics in tasks are automatically recovered and converted to errors
- **Context Cancellation**: Supports context cancellation and timeout
- **Fail-Fast**: In parallel execution, the first error encountered is returned
- **Error Propagation**: Original errors are preserved and properly propagated

```go
task := func(ctx context.Context) error {
    // This panic will be recovered and converted to an error
    panic("something went wrong")
}

future := goture.NewGoture(ctx, task)
err := future.Wait()
// err will contain the panic message as an error
```

## Context Support

Full integration with Go's context package:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

future := goture.NewGoture(ctx, longRunningTask)
err := future.Wait() // Will respect the timeout
```

## Performance Considerations

- Tasks start executing immediately upon creation
- Minimal overhead for goroutine management
- Efficient parallel coordination using channels
- No polling - uses Go's built-in synchronization primitives

## Use Cases

- **API Calls**: Make multiple HTTP requests concurrently
- **Database Operations**: Execute parallel database queries
- **File Processing**: Process multiple files simultaneously  
- **Pipeline Processing**: Create async processing pipelines
- **Background Tasks**: Execute non-blocking background operations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Examples

More examples can be found in the `examples_test.go` file in the repository.

---
