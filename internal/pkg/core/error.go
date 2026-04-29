package core

import (
	"errors"
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	xtremeres "github.com/globalxtreme/go-core/v2/response"
	"net/http"
	"os"
	"runtime/debug"
)

func ErrorHandler(fn func() error) error {
	errChan := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				bug := false

				var err error
				if panicData, ok := r.(*xtremeres.ResponseError); ok {
					status := panicData.Status
					bug = status.Bug

					errMsg := status.Message
					if status.InternalMsg != "" {
						errMsg += ". " + status.InternalMsg
					}

					err = errors.New(fmt.Sprintf("%s. Code: %d.", errMsg, status.Code))
				} else if panicData, ok := r.(error); ok {
					err = errors.New(fmt.Sprintf("%v. Code: %d.", panicData.Error(), http.StatusInternalServerError))
				} else {
					bug = true
					err = errors.New(fmt.Sprintf("An error Occurred. Code: %d.", http.StatusInternalServerError))
				}

				fmt.Fprintf(os.Stderr, "panic: %v\n", r)
				xtremepkg.LogError(r, bug)

				errChan <- err
			}
		}()

		if err := fn(); err != nil {
			xtremepkg.LogError(err, false)
			errChan <- err
		} else {
			close(errChan)
		}
	}()

	return <-errChan
}

func RabbitMQErrorHandler(fn func() (interface{}, error)) (res interface{}, err error, trace []byte) {
	resChan := make(chan interface{})
	errChan := make(chan error)
	traceChan := make(chan []byte)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				bug := false

				var err error
				if panicData, ok := r.(*xtremeres.ResponseError); ok {
					status := panicData.Status
					bug = status.Bug

					errMsg := status.Message
					if status.InternalMsg != "" {
						errMsg += ". " + status.InternalMsg
					}

					err = errors.New(fmt.Sprintf("%s. Code: %d.", errMsg, status.Code))
				} else if panicData, ok := r.(error); ok {
					err = errors.New(fmt.Sprintf("%v. Code: %d.", panicData.Error(), http.StatusInternalServerError))
				} else {
					bug = true
					err = errors.New(fmt.Sprintf("An error Occurred. Code: %d.", http.StatusInternalServerError))
				}

				fmt.Fprintf(os.Stderr, "panic: %v\n", r)
				xtremepkg.LogError(r, bug)

				errChan <- err
				traceChan <- debug.Stack()
			}
		}()

		res, err := fn()
		if err != nil {
			xtremepkg.LogError(err, false)

			errChan <- err
			traceChan <- nil
		} else {
			resChan <- res
		}
	}()

	select {
	case res := <-resChan:
		return res, nil, nil
	case err := <-errChan:
		trace := <-traceChan
		return nil, err, trace
	}
}
