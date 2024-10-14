package tool

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

// Context represents cancelable context which also reacts to os.Interrupt signal.
type Context struct {
	ctx                 context.Context
	ctxSignalNotifyStop context.CancelFunc
	ctxCancel           context.CancelFunc
	canceled            atomic.Bool
}

// NewContext returns new Context based on parent.
func NewContext(parent context.Context) *Context {
	ctx, cancel := context.WithCancel(parent)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	c := &Context{
		ctx:                 ctx,
		ctxSignalNotifyStop: stop,
		ctxCancel:           cancel,
		canceled:            atomic.Bool{},
	}
	return c
}

// Cancel context and stops listening.
// It can be called repeatedly and safely from simultaneous goroutines.
func (c *Context) Cancel() {
	if c.canceled.CompareAndSwap(false, true) {
		c.ctxCancel()
		c.ctxSignalNotifyStop()
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

var _ context.Context = (*Context)(nil)
