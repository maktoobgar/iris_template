package middlewares

import (
	"time"

	"service/pkg/errors"

	"github.com/kataras/iris/v12"
)

func Timeout(timeout time.Duration) iris.Handler {
	return func(ctx iris.Context) {
		// Create a channel to wait for the handler to complete.
		ch := make(chan struct{})

		// Call the next handler in a separate goroutine.
		go func() {
			Panic(ctx)
			close(ch)
		}()

		// Wait for either the handler to complete or the timeout to expire.
		select {
		case <-ch:
			// Handler completed successfully, do nothing.
		case <-time.After(timeout):
			// Handler timed out, return an error response.
			panic(errors.New(errors.TimeoutStatus, errors.Resend, "TimeoutError", ""))
		}
	}
}
