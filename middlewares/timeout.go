package middlewares

import (
	"context"
	"sync"
	"time"

	g "service/global"
	"service/pkg/errors"

	"github.com/kataras/iris/v12"
)

func Timeout(timeout time.Duration) iris.Handler {
	return func(ctx iris.Context) {
		// Create a channel to wait for the handler to complete.
		ch := make(chan any, 1)

		// Call the next handler in a separate goroutine.
		panicChan := make(chan any, 1)
		writerLock := &sync.Mutex{}
		closedWriter := false
		ctx.Values().Set(g.WriterLock, writerLock)
		ctx.Values().Set(g.ClosedWriter, closedWriter)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			ctx.Next()
			close(ch)
		}()

		newCtx, cancel := context.WithTimeout(ctx.Request().Context(), timeout)
		ctx.ResetRequest(ctx.Request().WithContext(newCtx))
		defer cancel()

		// Wait for either the handler to complete or the timeout to expire.
		select {
		case p := <-panicChan:
			panic(p)
		case <-ch:
			// Handler completed successfully, do nothing.
		case <-newCtx.Done():
			// Handler timed out, return an error response.
			panic(errors.New(errors.ServiceUnavailable, errors.Resend, "TimeoutError", ""))
		}
	}
}
