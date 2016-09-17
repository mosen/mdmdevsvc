package device

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostDevice(ctx context.Context, d *Device) error {
	mw.logger.Log("method", "PostDevice")
	return mw.next.PostDevice(ctx, d)
}

func (mw loggingMiddleware) GetDevice(ctx context.Context, uuidStr string) (Device, error) {
	mw.logger.Log("method", "GetDevice", "uuid", uuidStr)
	return mw.next.GetDevice(ctx, uuidStr)
}

func (mw loggingMiddleware) PutDevice(ctx context.Context, uuidStr string, d Device) error {
	mw.logger.Log("method", "PutDevice", "uuid", uuidStr)
	return mw.next.PutDevice(ctx, uuidStr, d)
}

func (mw loggingMiddleware) PatchDevice(ctx context.Context, uuidStr string, d Device) error {
	mw.logger.Log("method", "PatchDevice", "uuid", uuidStr)
	return mw.next.PatchDevice(ctx, uuidStr, d)
}

func (mw loggingMiddleware) DeleteDevice(ctx context.Context, uuidStr string) error {
	mw.logger.Log("method", "DeleteDevice", "uuid", uuidStr)
	return mw.next.DeleteDevice(ctx, uuidStr)
}
