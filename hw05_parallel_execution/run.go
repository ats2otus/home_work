package hw05parallelexecution

import (
	"context"
	"errors"
	"math"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// m < 0 - it means all errors will be ignored
func Run(tasks []Task, n, m int) error {
	if m < 0 {
		m = math.MaxInt
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue := make(chan Task, len(tasks))
	errors := make(chan error, len(tasks))

	for i := 0; i < n; i++ {
		go worker(ctx, queue, errors)
	}
	for i := 0; i < len(tasks); i++ {
		queue <- tasks[i]
	}
	close(queue)

	var totalErrs int
	for i := 0; i < len(tasks); i++ {
		if err := <-errors; err != nil {
			totalErrs++
			if totalErrs >= m {
				return ErrErrorsLimitExceeded
			}
		}
	}

	return nil
}

func worker(ctx context.Context, src chan Task, done chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-src:
			if !ok {
				return
			}
			done <- task()
		}
	}
}
