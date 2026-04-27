package application

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const gracefulShoutdownTimeout = time.Second * 15

type Module interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Hook struct {
	OnBeforeStart func(ctx context.Context) error
	OnAfterStart  func(ctx context.Context) error
	OnBeforeStop  func(ctx context.Context) error
	OnAfterStop   func(ctx context.Context) error
}

type Application struct {
	modules []Module
	hooks   []Hook
	started atomic.Bool
}

func New(modules ...Module) *Application {
	return &Application{
		modules: modules,
	}
}

func (a *Application) WithModules(modules ...Module) *Application {
	if a.started.Load() {
		panic("application: cannot add modules after start")
	}

	a.modules = append(a.modules, modules...)

	return a
}

func (a *Application) WithHooks(hooks ...Hook) *Application {
	if a.started.Load() {
		panic("application: cannot add hooks after start")
	}

	a.hooks = append(a.hooks, hooks...)

	return a
}

func (a *Application) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := a.Start(ctx); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), gracefulShoutdownTimeout)
	defer cancel()

	if err := a.Stop(ctx); err != nil {
		return fmt.Errorf("failed to graceful shoutdown: %w", err)
	}

	return nil
}

func (a *Application) Start(ctx context.Context) error {
	pipe := pipeline{
		a.beforeStart,
		a.start,
		a.afterStart,
	}

	return pipe.run(ctx)
}

func (a *Application) Stop(ctx context.Context) error {
	pipe := pipeline{
		a.beforeStop,
		a.stop,
		a.afterStop,
	}

	return pipe.run(ctx)
}

func (a *Application) start(ctx context.Context) error {
	if !a.started.CompareAndSwap(false, true) {
		return errors.New("application: already started") //nolint:err113
	}

	var wg sync.WaitGroup

	done := make(chan error, len(a.modules))

	for _, module := range a.modules {
		wg.Add(1)

		go func(module Module) {
			defer wg.Done()

			if err := module.Start(ctx); err != nil {
				done <- err
			}
		}(module)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	// Возвращаем первую ошибку, если она есть
	if err := <-done; err != nil {
		return err
	}

	return nil
}

func (a *Application) stop(ctx context.Context) error {
	var wg sync.WaitGroup

	errs := make([]error, len(a.modules))

	for i, module := range a.modules {
		wg.Add(1)

		go func(i int, module Module) {
			defer wg.Done()

			if err := module.Stop(ctx); err != nil {
				errs[i] = err
			}
		}(i, module)
	}

	wg.Wait()

	return errors.Join(errs...)
}

func (a *Application) beforeStart(ctx context.Context) error {
	for _, hook := range a.hooks {
		if hook.OnBeforeStart != nil {
			if err := hook.OnBeforeStart(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Application) afterStart(ctx context.Context) error {
	for _, hook := range a.hooks {
		if hook.OnAfterStart != nil {
			if err := hook.OnAfterStart(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Application) beforeStop(ctx context.Context) error {
	for _, hook := range a.hooks {
		if hook.OnBeforeStop != nil {
			if err := hook.OnBeforeStop(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Application) afterStop(ctx context.Context) error {
	for _, hook := range a.hooks {
		if hook.OnAfterStop != nil {
			if err := hook.OnAfterStop(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

type pipeline []func(ctx context.Context) error

func (p pipeline) run(ctx context.Context) error {
	for _, fn := range p {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}
